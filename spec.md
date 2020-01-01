# 명세

```
수정 로그
- v1.0.0 첫 배포(20191229)
```

## API 목록
| 메소드 | 경로 | 설명 | 파라미터 |
| --- | --- | --- | --- |
| GET | /board | 모든 게시판 정보 가져오기 | - |
| GET | /board/:boardID | boardID에 해당하는 게시판 정보 가져오기 | - |
| GET | /post/:input | 게시판, 게시물 또는 게시판의 모든 게시물 정보 가져오기 | { type: [ "board" | "post" ] , postId: int, before?: "true" } |
| GET | /file | 현재 저장되어 있는 파일들의 파일명, 크기, 저장된 날짜 가져오기 | - |
| POST | /auth | JJUNGS 권한 인증 | { pw:string } |
| POST | /auth/val | 권한 확인 | { token:string } |
| POST | /admin/board | 게시판 생성 | { name: string, url: string, read: string } |
| POST | /admin/file | 파일 업로드 | FormData(file: file) |
| POST | /admin/post | 게시물 생성 | { boardId: int, title: string, body?: string, tags?: string, description?: string } |
| PUT | /admin/board | 게시판 수정 | { id: int, name: string, url: string, read: string } |
| PUT | /admin/post | 게시물 수정 | { id: int, boardId: int, title: string, body: string, tags: string, description: string } |
| DELETE | /admin/board | 게시판 삭제 | { id: int } |
| DELETE | /admin/post | 게시물 삭제 | { id: int } |
| DELETE | /admin/file | 파일 삭제 | { id: int } |

## 오류
### 에러 코드
| 에러 코드 | 경로 | 설명 |
| --- | --- | --- |
| ERR400 | 인증 | 입력형태가 올바르지 않은 경우 |
| ERR401 | 인증 | 인증에 실패한 경우 |
| ERR500 | 공통 | 알 수 없는 오류 |
| TKN000 | 인증 | 비밀번호가 일치하지 않는 경우 |
| TKN001 | 인증 | 올바른 토큰이 아닌 경우 |
| FIL000 | 파일 | 파일 이름이 중복되는 경우 |