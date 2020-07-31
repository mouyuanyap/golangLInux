
FROM golang AS builder

WORKDIR /home/mouyuan/go/src/github.com/mouyuan4598/hello

COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 go build -o /go/bin/helloo

FROM scratch
COPY --from=builder /go/bin/helloo /go/bin/helloo 
ENTRYPOINT ["/go/bin/helloo"]