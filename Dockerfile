FROM centos:centos6

RUN yum install -y wget

WORKDIR /tmp
RUN mkdir $HOME/go
RUN wget https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.6.2.linux-amd64.tar.gz
RUN rm go1.6.2.linux-amd64.tar.gz

ENV GOPATH=$HOME/go
ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /tmp
RUN yum -y install curl-devel expat-devel gettext-devel openssl-devel zlib-devel perl-ExtUtils-MakeMaker gcc
RUN wget https://www.kernel.org/pub/software/scm/git/git-2.9.0.tar.gz
RUN tar zxvf git-2.9.0.tar.gz
RUN rm git-2.9.0.tar.gz

WORKDIR /tmp/git-2.9.0
RUN make prefix=/usr/local all
RUN make prefix=/usr/local install
RUN cd ../ && rm -rf git-2.9.0

RUN printf '[nginx] \n\
name=nginx repo \n\
baseurl=http://nginx.org/packages/mainline/centos/6/$basearch/ \n\
gpgcheck=0 \n\
enabled=1' >> /etc/yum.repos.d/nginx.repo
RUN yum install -y nginx
RUN pgrep -f '/sbin/udevd' | xargs kill

RUN rm /etc/nginx/conf.d/*.conf
COPY  ./go_web.conf /etc/nginx/conf.d/

WORKDIR /root/

RUN go get -u github.com/kataras/iris/iris

COPY ./*.go /src/

WORKDIR /src

EXPOSE 80

CMD /etc/init.d/nginx start && go run /src/main.go