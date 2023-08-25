FROM golang:alpine AS builder

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./main/main.go

RUN go build -o todo-app ./main/main.go

#lightweight docker container with binary
FROM alpine:latest

#RUN apk --no-cache add ca-certificates

COPY --from=builder COPY ./ ./

EXPOSE 80

CMD ["./todo-app"]