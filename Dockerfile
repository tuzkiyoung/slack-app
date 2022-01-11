FROM alpine

COPY slack-go-demo /app/shallowwater

ENTRYPOINT [ "/app/shallowwater" ]