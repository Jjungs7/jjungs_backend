# 명세

```
수정 로그
- v1.0.0 첫 배포(20191229)
```

## API 목록
| 메소드 | 경로 | 설명 | 파라미터 |
| --- | --- | --- | --- |
| GET | /auth | JJUNGS 권한 인증 | [ pw:string ] |
| GET | /board | JJUNGS 권한 인증 | [ pw:string ] |
| GET | /board/:url | JJUNGS 권한 인증 | [ pw:string ] |
| GET | /auth | JJUNGS 권한 인증 | [ pw:string ] |
| GET | /auth | JJUNGS 권한 인증 | [ pw:string ] |
| GET | /auth | JJUNGS 권한 인증 | [ pw:string ] |

## 오류
### 에러 코드
| 에러 코드 | 경로 | 설명 |
| --- | --- | --- |
| ERR400 | 인증 | 입력형태가 올바르지 않은 경우 |
| ERR401 | 인증 | 인증에 실패한 경우 |
| ERR500 | 공통 | 알 수 없는 오류 |
| TKN000 | 인증 | 비밀번호가 일치하지 않는 경우 |
| TKN001 | 인증 | 올바른 토큰이 아닌 경우 |