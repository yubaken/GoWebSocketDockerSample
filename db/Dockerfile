FROM mysql:5.7.13

COPY *.sql /docker-entrypoint-initdb.d

ENV MYSQL_ROOT_PASSWORD password
ENV MYSQL_DATABASE chat
ENV MYSQL_USER user
ENV MYSQL_PASSWORD password

RUN /entrypoint.sh mysqld & sleep 30 && killall mysqld
RUN rm /docker-entrypoint-initdb.d/*.sql