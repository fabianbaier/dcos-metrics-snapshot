FROM alpine
RUN mkdir /app && apk add --no-cache curl bash
COPY ./dcos_crawler /app/dcos_crawler
ENTRYPOINT ["/app/dcos_crawler"]
