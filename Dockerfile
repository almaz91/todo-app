FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o todo-app ./main/main.go

RUN chmod a+x ./todo-app.sh

#CMD ["./todo-app"]