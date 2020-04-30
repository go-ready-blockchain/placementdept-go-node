FROM golang:alpine AS builder
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true
# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
ENV GOFLAGS=-mod=mod



# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY src/go.mod .
COPY src/go.sum .
RUN ls

RUN go mod download

# Copy the code into the container
COPY src/ .
RUN ls
#RUN apk add --no-cache git
#RUN go get -d -v ./...
#RUN go install -v ./... 
#RUN CGO_ENABLED=0 go-wrapper install -ldflags '-extldflags "-static"'

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .


# Build a small image
FROM scratch

#COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /dist/main /



# Command to run
ENTRYPOINT ["/main"]