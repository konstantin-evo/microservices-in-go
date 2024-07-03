FROM alpine:3.20.1

RUN mkdir /app

COPY mailerApp /app
COPY templates /templates

CMD [ "/app/mailerApp"]