FROM golang:1.13-alpine as builder

WORKDIR /usr/src/app
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

FROM scratch
COPY --from=builder /usr/src/app/main /main

ENTRYPOINT [ "/main" ]
