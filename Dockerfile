FROM ghmlee/scratch:latest
MAINTAINER Graham Lee <ghmlee@ghmlee.com>

WORKDIR /opt/ecr

COPY ./bin ./bin

ENTRYPOINT ["./bin/ecr"]
