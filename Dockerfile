FROM golang:alpine3.19
WORKDIR /app
COPY .env  ./
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
CMD ["./main"]