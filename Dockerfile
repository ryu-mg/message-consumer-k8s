# 빌드 스테이지
FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod tidy

COPY . .

# 빌드
# -ldflags '-s': 빌드 시 불필요한 심볼 정보를 제거하여 바이너리 크기 최소화
# -o message-consumer: 바이너리 파일명
# ./cmd/consumer/main.go: 소스 파일 경로. main 파일만 빌드해도 나머지 패키지도 컴파일
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-s' -o message-consumer ./cmd/consumer/main.go

# 실행 스테이지
FROM alpine:latest

# 필요한 CA 인증서 설치. apline 이미지는 기본적으로 인증서가 없음
# 그래서 http 통신 or tls 기반 통신을 할 때 인증서가 필요함
RUN apk --no-cache add ca-certificates

WORKDIR /app

# 빌드된 바이너리만 복사
COPY --from=builder /app/message-consumer .
