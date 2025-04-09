# 빌드 스테이지
FROM golang:1.24.2 AS builder

# 작업 디렉토리 설정
WORKDIR /app

# 의존성 관리 파일 복사 및 다운로드
COPY go.mod go.sum ./
RUN go mod download && go mod tidy

# 소스 코드 복사
COPY . .

# 애플리케이션 빌드
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-s' -o message-consumer ./cmd/consumer/main.go

# 실행 스테이지
FROM alpine:latest

# 필요한 CA 인증서 설치
RUN apk --no-cache add ca-certificates

WORKDIR /app

# 빌드된 바이너리만 복사
COPY --from=builder /app/message-consumer .
