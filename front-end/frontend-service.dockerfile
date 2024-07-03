FROM alpine:3.20.1

RUN mkdir -p /app/cmd/web/templates

COPY frontApp /app/frontApp
COPY cmd/web/templates /app/cmd/web/templates

WORKDIR /app

CMD [ "/app/frontApp"]