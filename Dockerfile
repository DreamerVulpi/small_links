FROM golang:latest
WORKDIR /app
COPY . .
RUN go get github.com/redis/go-redis/v9
RUN go get github.com/gin-gonic/gin
CMD ["go", "run", "main.go"]