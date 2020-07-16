FROM golang
ADD . /go/src/github.com/mouyuan4598/net/html
WORKDIR /go/src/github.com/mouyuan4598/net/html
RUN go build main.go
FROM ubuntu
COPY --from=0 go/src/github.com/mouyuan4598/net/html .
CMD ["./main"]
EXPOSE 8080
