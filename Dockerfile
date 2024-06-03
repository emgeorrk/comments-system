FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -buildvcs=false -o /commentsSystem .
RUN chmod +x /commentsSystem

EXPOSE 8084

CMD ["/commentsSystem"]
