FROM alpine:3.20.1

RUN mkdir /app

COPY frontApp /app

CMD [ "/app/frontApp"]