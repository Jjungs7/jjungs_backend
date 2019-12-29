package models

import (
	"fmt"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"

	"jjungs_backend/components/database"
)

type Post struct {
	ID int `gorm:"primary_key"`

	BoardID int `sql:"index"`
	Title string `gorm:"type:varchar(255);not null"`
	Body string
	Description string
	PostTags []string `gorm:"-"`
	Hits int `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type PostTag struct {
	PostID int `gorm:"unique_index:uix_post_tags_post_id_keyword"`
	Keyword string `gorm:"type:varchar(20);unique_index:uix_post_tags_post_id_keyword"`
}

func getPostTags(PostID int) []string {
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

func getAll(isAdmin bool, postID int, before bool) ([]Post, int, int) {
	var posts []Post
	count := 20

	subquery := database.DB.Model(&Post{})
	if !isAdmin {
		subquery = subquery.Select("posts.*")
		subquery = subquery.Joins("inner join boards on posts.board_id=boards.id")
		subquery = subquery.Where("boards.read_permission <> ?", "JJUNGS")
	}

	query := subquery
	if before {
		query = query.Where("posts.id>?", postID).Order("posts.id asc")
	} else {
		query = query.Where("posts.id<?", postID).Order("posts.id desc")
	}
	query.Limit(count).Find(&posts)
	if len(posts) <= 0 {
		return posts, 0, 0
	}

	sort.SliceStable(posts, func(i, j int) bool { return posts[i].ID > posts[j].ID })
	var prev int
	var next int
	subquery.Where("posts.id>?", posts[0].ID).Count(&prev)
	subquery.Where("posts.id<?", posts[len(posts)-1].ID).Count(&next)
	return posts, prev, next
}

func getPostsInBoard(boardID string, isAdmin bool, postID int, before bool) ([]Post, int, int) {
	var posts []Post
	count := 20
	board := new(Board)
	database.DB.Where("boards.id=?", boardID).First(&board)
	if !isAdmin && board.ReadPermission == "JJUNGS" {
		return posts, 0, 0
	}

	query := database.DB.Where("board_id=?", boardID)
	if before {
		query = query.Where("id>?", postID).Order("id asc")
	} else {
		query = query.Where("id<?", postID).Order("id desc")
	}
	query.Limit(count).Find(&posts)
	if len(posts) <= 0 {
		return posts, 0, 0
	}

	sort.SliceStable(posts, func(i, j int) bool { return posts[i].ID > posts[j].ID })
	var prev int
	var next int
	database.DB.Model(&Post{}).Where("board_id=? and id>?", boardID, posts[0].ID).Count(&prev)
	database.DB.Model(&Post{}).Where("board_id=? and id<?", boardID, posts[len(posts)-1].ID).Count(&next)
	return posts, prev, next
}

func getPost(postID string) Post {
	var post Post
	database.DB.Where("id="+postID).First(&post)
	return post
}

func GetPosts(c *gin.Context) {
	permissions, _ := c.Get("permissions")
	input := c.Param("input")
	t := c.Query("type")
	postID, err := strconv.Atoi(c.Query("postId"))
	_before := c.Query("before")
	before := _before == "true"
	isAdmin := permissions == "JJUNGS"
	if err != nil {
		postID = math.MaxInt32
	}

	if t == "board" {
		var posts []Post
		var prev int
		var next int
		if input == "0" {
			posts, prev, next = getAll(isAdmin, postID, before)
		} else {
			posts, prev, next = getPostsInBoard(input, isAdmin, postID, before)
		}
		c.JSON(200, gin.H{
			"data": gin.H{
				"posts": posts,
				"prevCnt": prev,
				"nextCnt": next,
			},
		})
	} else if t == "post" {
		post := getPost(input)
		if post.ID == 0 {
			c.JSON(200, gin.H{
				"data": nil,
			})
			return
		}

		board := &Board{}
		database.DB.Where(&Board{ID: post.BoardID}).First(&board)
		if !isAdmin && board.ReadPermission == "JJUNGS" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "ERR401",
			})
			return
		}

		// Update hits before response
		cookieName := "last_hit_"+strconv.Itoa(post.ID)
		cookie, err1 := c.Cookie(cookieName)
		if err1 != nil {
			database.DB.Model(&post).UpdateColumn("hits", gorm.Expr("hits+1"))
			post.Hits++
		} else {
			fmt.Println(cookie)
		}
		c.SetCookie(cookieName, "hit", 600, "/", "", false, false)

		post.PostTags = getPostTags(post.ID)
		c.JSON(200, gin.H{
			"data": post,
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