package models

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"jjungs_backend/components/database"
)

type Post struct {
	ID int `gorm:"primary_key"`

	BoardID int `sql:"index"`
	Board *Board `gorm:"foreignkey:BoardID;association_foreignkey:ID"`
	Title string `gorm:"type:varchar(255);not null"`
	Body string
	Description string
	PostTags []string `gorm:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type PostTag struct {
	PostID int `gorm:"unique_index:uix_post_tags_post_id_keyword"`
	Keyword string `gorm:"type:varchar(20);unique_index:uix_post_tags_post_id_keyword"`
}

func getPostTags(PostID int) ([]string) {
	var tags []string
	rows, _ := database.DB.Table("post_tags").Select("post_tags.keyword").Joins("inner join posts on post_tags.post_id=posts.id").Where("posts.id="+strconv.Itoa(PostID)).Rows()
	defer rows.Close()
	for rows.Next() {
		var tag string
		rows.Scan(&tag)
		tags = append(tags, tag)
	}
	return tags
}

func getAll(isAdmin bool, from int, to int) []Post {
	var posts []Post
	database.DB.Order("id desc").Find(&posts)
	for idx, _ := range posts {
		posts[idx].Board = new(Board)
		database.DB.First(&posts[idx].Board, "boards.id=?", posts[idx].BoardID)
		posts[idx].PostTags = getPostTags(posts[idx].ID)
	}

	if !isAdmin {
		for i := len(posts)-1; i>=0; i-- {
			if posts[i].Board.ReadPermission == "JJUNGS" {
				posts = append(posts[:i], posts[i+1:]...)
			}
		}
	}
	return posts
}

func getPostsInBoard(boardID string, isAdmin bool, from int, to int) []Post {
	var posts []Post
	database.DB.Where("board_id="+boardID).Order("id desc").Find(&posts)
	for idx, _ := range posts {
		posts[idx].Board = new(Board)
		database.DB.First(&posts[idx].Board, "boards.id=?", posts[idx].BoardID)
		posts[idx].PostTags = getPostTags(posts[idx].ID)
	}

	if !isAdmin {
		for i := len(posts)-1; i>=0; i-- {
			if posts[i].Board.ReadPermission == "JJUNGS" {
				posts = append(posts[:i], posts[i+1:]...)
			}
		}
	}
	return posts
}

func getPost(postID string) Post {
	var post Post
	database.DB.Where("id="+postID).First(&post)
	if post.ID == 0 {
		return post
	}

	post.Board = &Board{}
	database.DB.First(&post.Board, "boards.id=?", post.BoardID)
	return post
}

func GetPosts(c *gin.Context) {
	permissions, _ := c.Get("permissions")
	input := c.Param("input")
	t := c.Query("type")
	isAdmin := permissions == "JJUNGS"
	if t == "board" {
		posts := getPostsInBoard(input, isAdmin, 0, 0)
		c.JSON(200, gin.H{
			"data": posts,
		})
	} else if t == "post" {
		post := getPost(input)
		if !isAdmin && post.Board.ReadPermission == "JJUNGS" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "ERR401",
			})
			return
		}

		if post.ID == 0 {
			c.JSON(200, gin.H{
				"data": nil,
			})
			return
		}

		post.PostTags = getPostTags(post.ID)
		c.JSON(200, gin.H{
			"data": post,
		})
	} else if input == "" && t == "" {
		posts := getAll(isAdmin, 0, 0)
		c.JSON(200, gin.H{
			"data": posts,
		})
	} else {
		c.JSON(200, gin.H{
			"error": "ERR400",
			"message": "you must specify type(board | post | all) with query string. ex) /post/1?type=board",
		})
	}
}

type PostInput struct {
	ID int `json:"id"`
	BoardID int `json:"boardId"`
	Title string `json:"title"`
	Body string `json:"body"`
	Tags string `json:"tags"`
	Description string `json:"description"`
}

func deleteTagsExcluding(postID int, excludingTags []string) {
	if len(excludingTags) == 0 {
		return
	}
	var tags string
	for idx, tag := range excludingTags {
		if idx != 0 {
			tags += ","
		}
		tags += "'"+ tag +"'"
	}
	database.DB.Exec("DELETE FROM post_tags WHERE post_id=? AND keyword NOT IN (" + tags + ")", strconv.Itoa(postID))
}

func insertIgnoreDuplicateTags(postID int, tags []string) {
	if len(tags) == 0 {
		return
	}
	var tagsConverted string
	stringPostID := strconv.Itoa(postID)
	for _, tag := range tags {
		if len(tagsConverted) > 0 {
			tagsConverted += ","
		}

		tagsConverted += "(" + stringPostID + ",'" + tag + "')"
	}
	database.DB.Exec("INSERT INTO post_tags VALUES" + tagsConverted + " ON CONFLICT(post_id, keyword) DO NOTHING")
}

func getWellFormedTag(str string) string {
	leadingTrailingWhtspc := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	insideWhtspc := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	str = leadingTrailingWhtspc.ReplaceAllString(str, "")
	str = insideWhtspc.ReplaceAllString(str, " ")
	return str
}

func CreatePost(c *gin.Context) {
	var input PostInput
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
		})
		fmt.Println(input)
		fmt.Println(err)
		return
	}

	post := Post{
		Title: input.Title,
		Body: input.Body,
		BoardID: input.BoardID,
		Description: input.Description,
	}

	if post.Title == "" || post.BoardID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		return
	}

	errs := database.DB.Save(&post).GetErrors()
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		fmt.Println(errs)
		return
	}

	if input.Tags != "" {
		post.PostTags = strings.Split(getWellFormedTag(input.Tags), " ")
		insertIgnoreDuplicateTags(post.ID, post.PostTags)
	}

	post.Board = new(Board)
	database.DB.First(&post.Board, "boards.id=?", post.BoardID)
	c.JSON(http.StatusOK, gin.H{
		"data": post,
	})
}

func UpdatePost(c *gin.Context) {
	var input PostInput
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
		})
		fmt.Println(err)
		return
	}

	var post Post
	database.DB.First(&post, "posts.id=?", input.ID)
	if post.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		return
	}

	if input.BoardID > 0 {
		post.BoardID = input.BoardID
	}

	if input.Title != "" {
		post.Title = input.Title
	}

	if input.Body != "" {
		post.Body = input.Body
	}

	if input.Tags != "" {
		post.PostTags = strings.Split(getWellFormedTag(input.Tags), " ")
		deleteTagsExcluding(post.ID, post.PostTags)
		insertIgnoreDuplicateTags(post.ID, post.PostTags)
	}

	errs := database.DB.Save(&post).GetErrors()
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		fmt.Println(errs)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": post,
	})
}

func DeletePost(c *gin.Context) {
	var postInput PostInput
	if err := binding.JSON.Bind(c.Request, &postInput); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		return
	}

	if postInput.ID <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": postInput.ID,
		})
		return
	}

	errs := database.DB.Delete(&Post{ID: postInput.ID}).GetErrors()
	if len(errs) > 0 {
		fmt.Println(errs)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": postInput.ID,
	})
}