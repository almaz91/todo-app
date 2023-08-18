FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o todo-app ./main/main.go

CMD ["./todo-app"]