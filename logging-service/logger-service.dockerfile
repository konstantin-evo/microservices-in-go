FROM alpine:3.20.1

RUN mkdir /app

COPY loggerServiceApp /app

CMD [ "/app/loggerServiceApp"]