#################
# Build stage 0 #
#################
FROM golang:1.10-alpine

ARG DOCKER_IMAGE

# Create work dir
RUN mkdir -p /go/src/github.com/danbondd/vanity

# Set work dir
WORKDIR /go/src/github.com/danbondd/vanity

# Copy files
COPY . .

# Build binary
RUN GOOS=linux go install .

#################
# Build stage 1 #
#################
FROM alpine

# Copy the binaries from build stage 0
COPY --from=0 /go/bin/ /usr/local/bin/

# Expose web server port
EXPOSE 8080

ENTRYPOINT ["/bin/sh", "-c"]

CMD ["vanity"]
