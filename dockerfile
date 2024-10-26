FROM golang:1.23

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -o /my_app

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app
ENV TODO_PASSWORD=TestPas

CMD ["/my_app"]