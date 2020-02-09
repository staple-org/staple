FROM alpine
RUN apk add -u ca-certificates
COPY ./build/linux/amd64/staple /app/
COPY ./frontend/build /app/frontend

EXPOSE 9998

WORKDIR /app/
ENTRYPOINT [ "/app/staple" ]
