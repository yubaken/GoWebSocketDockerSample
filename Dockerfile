FROM golang:1.6.2

RUN go get github.com/mattn/gom

RUN mkdir -p /opt/iris/
RUN mkdir -p /opt/iris_dev/

COPY ./*.go /opt/iris/
COPY ./Gomfile /opt/iris/
COPY ./static /opt/iris/static
COPY ./templates /opt/iris/templates

#TODO たまにダウンロードに失敗する
#WORKDIR /opt/iris/
#RUN gom install

WORKDIR /opt/
RUN printf "#!/bin/bash \n\
if [ -e /opt/iris_dev/main.go ]; then \n\
  cd /opt/iris_dev/ \n\
  chmod -R 775 * \n\
  echo 'Running Develop Env' \n\
  gom run /opt/iris_dev/main.go \n\
else \n\
  cd /opt/iris/ \n\
  chmod -R 775 * \n\
  echo 'Running Product Env' \n\
  gom run /opt/iris/main.go \n\
fi" >> iris_run.sh
RUN chmod +x iris_run.sh

EXPOSE 80

CMD /opt/iris_run.sh