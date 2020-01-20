# jjungs_backend [![travis build status](https://travis-ci.com/Jjungs7/jjungs_backend.svg?token=HqXpQvLubk2CGuZGBKfS&branch=master)](https://travis-ci.com/Jjungs7/jjungs_backend/builds)
완성?!

### 환경변수 설정
*.envrc.ex*을 복사하고 파일명을 *.envrc*로 바꿉니다. 그리고 *.envrc*에 적절히 값을 넣어줍니다
[direnv](https://direnv.net/)를 사용하면 편리하게 환경변수를 설정할 수 있습니
``` bash
cp .envrc.ex .envrc
vim .envrc
# Change env variables

direnv allow
```

### 실행
``` bash
# install dependencies
go mod tidy

# run
# needs postgresql running on background
go build -o main main.go
./main
```

### 배포
``` bash
# Build
docker build -t equisde/jjungs:(version) .

# Run
docker run --rm --name jjungs-api [ -e [env var1] -e [env var2] ... ] -p 8080:8080 equisde/jjungs:(version)
```
