FROM golang:1.15.6-alpine3.12 as builder
RUN apk add git
RUN mkdir /algoAPI
ADD . /algoAPI
WORKDIR /algoAPI

COPY go.mod go.sum swagger.yaml swagger.json help.json ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

RUN mkdir /algoAPI

WORKDIR /algoAPI/

COPY --from=builder /algoAPI/main .
COPY --from=builder /algoAPI/swagger.yaml .
COPY --from=builder /algoAPI/swagger.json .
COPY --from=builder /algoAPI/help.json .

ARG DBpw_arg=default_value 
ENV DBpw=$DBpw_arg

ARG JWTSecret_arg=default_value 
ENV JWTSecret=$JWTSecret_arg

EXPOSE 8090

CMD ["./main"]
