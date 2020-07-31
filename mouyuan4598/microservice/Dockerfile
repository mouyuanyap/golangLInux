
FROM golang AS builder

WORKDIR /home/mouyuan/go/src/github.com/mouyuan4598/microservice

COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 go build -o /go/bin/microservice

FROM alpine:latest
COPY --from=builder /go/bin/microservice /go/bin/microservice 

ENTRYPOINT ["/go/bin/microservice"]