# build a tiny docker image 
FROM alpine:latest

RUN mkdir /app

COPY brokerApp /app

CMD ["/app/brokerApp"]
