#################
# Build stage 0 #
#################
FROM golang:1.12-stretch
WORKDIR /go/src/github.com/syscll/vanity
COPY . .
RUN go test -v -cover ./...
RUN go install

#################
# Build stage 1 #
#################
FROM debian:stretch
# Copy binaries from build stage 0
COPY --from=0 /go/bin/ /usr/local/bin/
# Expose web server port
EXPOSE 8080
ENTRYPOINT ["/bin/bash"]
