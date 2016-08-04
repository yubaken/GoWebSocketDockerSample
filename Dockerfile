FROM golang:1.6.2

RUN go get github.com/mattn/gom

RUN mkdir -p /opt/iris/
RUN mkdir -p /opt/iris_dev/

COPY ./*.go /opt/iris/
COPY ./Gomfile /opt/iris/
COPY ./static /opt/iris/static
COPY ./templates /opt/iris/templates

WORKDIR /opt/iris/
RUN gom install

# TODO Gomでrunしてもsrcディレクトリ配下でないと失敗する
WORKDIR /opt/
RUN printf "#!/bin/bash \n\
  cd /opt/iris/vendor \n\
  if [ ! -e src ]; then \n\
    mkdir src \n\
    mv github.com/ src/ \n\
    mv golang.org/ src/ \n\
    mv gopkg.in/ src/ \n\
  fi \n\
  cd /opt/iris/ \n\
  echo 'Running Product Env' \n\
  gom run /opt/iris/main.go" >> iris_run.sh
RUN chmod +x iris_run.sh

EXPOSE 80

CMD /opt/iris_run.sh