FROM golang:1.6.2

RUN go get -u github.com/kataras/iris/iris

COPY ./* /src/

WORKDIR /src

EXPOSE 8080

CMD go run /src/main.go