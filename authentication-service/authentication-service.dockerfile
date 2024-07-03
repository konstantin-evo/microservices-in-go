FROM alpine:3.20.1

RUN mkdir /app

COPY authApp /app

CMD [ "/app/authApp"]