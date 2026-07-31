FROM golang:1.26

WORKDIR /build
COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/server
EXPOSE 8080

CMD [ "./main" ]