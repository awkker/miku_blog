FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/miku-blog ./main.go && \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/miku-migrate ./cmd/migrate/main.go && \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/miku-seed ./cmd/seed/main.go

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /out/miku-blog /usr/local/bin/miku-blog
COPY --from=builder /out/miku-migrate /usr/local/bin/miku-migrate
COPY --from=builder /out/miku-seed /usr/local/bin/miku-seed
COPY --from=builder /src/sql ./sql

EXPOSE 8080

CMD ["miku-blog"]
