FROM golang:1.23.0-alpine AS build

WORKDIR /code

COPY . .

RUN go build -o ./bin/main ./cmd/main.go

EXPOSE 1337

CMD [ "./bin/main" ]
