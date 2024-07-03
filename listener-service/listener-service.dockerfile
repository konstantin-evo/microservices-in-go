FROM alpine:3.20.1

RUN mkdir /app

COPY listenerApp /app

CMD [ "/app/listenerApp"]
