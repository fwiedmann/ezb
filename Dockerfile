FROM golang:1.15.3-alpine3.12 as builder
# All these steps will be cached
WORKDIR /build
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 go build -o ezb main.go
FROM scratch
COPY --from=builder /build/ezb .
ENTRYPOINT ["./ezb"]