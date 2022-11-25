FROM alpine:latest

COPY slack /app/slack

ENTRYPOINT [ "/app/slack" ]