FROM alpine
RUN apk add -u ca-certificates
COPY ./build/linux/amd64/staple /app/

WORKDIR /app/
ENTRYPOINT [ "/app/staple" ]
