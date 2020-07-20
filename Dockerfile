FROM blkpark/scratch:latest
MAINTAINER Blake Park <blkpark@blkpark.com>

WORKDIR /opt/ecr

COPY ./bin ./bin

ENTRYPOINT ["./bin/ecr"]
