# build a tiny docker image 
FROM alpine:3.20.1

RUN mkdir /app

COPY brokerApp /app

CMD ["/app/brokerApp"]
