FROM golang:alpine AS builder

RUN go version
ENV GOPATH=/
COPY ./ ./

RUN go mod download
RUN go build -o todo-app ./main/main.go

#lightweight docker container with binary
FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder ./todo-app .
COPY --from=builder ./config/ ./config/
COPY --from=builder ./schema/ ./schema/

EXPOSE 80

CMD ["./todo-app"]