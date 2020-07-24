
FROM golang AS builder

WORKDIR /home/mouyuan/go/src/github.com/mouyuan4598/gin

COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 go build -o /go/bin/gin

FROM scratch
COPY --from=builder /go/bin/gin /go/bin/gin 
ENTRYPOINT ["/go/bin/gin"]