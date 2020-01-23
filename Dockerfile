FROM alpine
RUN apk add -u ca-certificates
ADD ./bin/staple /app/

WORKDIR /app/
ENTRYPOINT [ "/app/staple" ]
