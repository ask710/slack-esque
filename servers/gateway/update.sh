#! /usr/bin/env bash
export MYSQL_ROOT_PASSWORD=lkdsnalkfnadkjbdflajbslajbd
export MYSQL_DATABASE=users
export MYSQL_ADDR=usersdb:3306

export REDISADDR=sessionServer:6379
export SUMMARYADDR=summary:80
export MESSAGESADDR=messages:80
export SESSIONKEY=$(openssl rand -hex 32)

export MQADDR=messagequeue:5672
export MQNAME=messagequeue

export DSN="root:$MYSQL_ROOT_PASSWORD@tcp($MYSQL_ADDR)/$MYSQL_DATABASE?parseTime=true"

docker rm -f summary
docker rm -f messages
docker rm -f gateway
docker rm -f usersdb
docker rm -f sessionServer
docker rm -f messagequeue
docker network rm authnet
docker network create authnet


docker pull ask710/usersdb
# -v /gateway/data:/var/lib/mysql \
docker run -d \
--network authnet \
--name usersdb \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=$MYSQL_DATABASE \
ask710/usersdb --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci



docker run -d \
--network authnet \
--name sessionServer \
redis

docker run -d \
--network authnet \
--name messagequeue \
rabbitmq:3-alpine

docker pull ask710/gateway

docker run -d \
--network authnet \
--name gateway \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSKEY=/etc/letsencrypt/live/api.ask710.me/privkey.pem \
-e TLSCERT=/etc/letsencrypt/live/api.ask710.me/fullchain.pem \
-e DSN=$DSN \
-e SESSIONKEY=$SESSIONKEY \
-e REDISADDR=$REDISADDR \
-e SUMMARYADDR=$SUMMARYADDR \
-e MESSAGESADDR=$MESSAGESADDR \
-e MQADDR=$MQADDR \
-e MQNAME=$MQNAME \
ask710/gateway




