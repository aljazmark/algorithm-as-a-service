FROM golang:1.15.6-buster as builder
RUN apt-get install git

RUN mkdir /algo
ADD . /algo
WORKDIR /algo

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

FROM golang:1.15.6-buster

RUN mkdir /algo

WORKDIR /algo/

COPY --from=builder /algo/main .
COPY --from=builder /algo/kcenter .

EXPOSE 8080

CMD ["./main"]