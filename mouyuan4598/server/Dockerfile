FROM golang AS builder

WORKDIR /home/mouyuan/go/src/github.com/mouyuan4598/server
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN CGO_ENABLED=0 go build -o /go/bin/server
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/server /go/bin/server
# Run the hello binary.
ENTRYPOINT ["/go/bin/server"]