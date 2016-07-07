FROM golang:1.6.2

RUN go get -u github.com/kataras/iris/iris
RUN go get -u github.com/go-sql-driver/mysql

RUN mkdir -p /opt/iris/

COPY ./*.go /opt/iris/
COPY ./static /opt/iris/static
COPY ./templates /opt/iris/templates

WORKDIR /opt/iris/
RUN chmod -R 775 *
RUN mkdir -p /opt/iris/certs

EXPOSE 80

CMD go run /opt/iris/main.go