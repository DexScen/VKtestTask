FROM golang:1.23

WORKDIR /app

RUN apt-get update && apt-get install -y docker.io

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o pinger ./main.go

CMD ["/app/pinger"]
