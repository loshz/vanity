FROM golang:1.15-buster
WORKDIR $GOPATH/src/github.com/syscll/vanity
COPY . .
RUN go test -v ./...
RUN CGO_ENABLED=0 go install

FROM alpine:3.12
COPY --from=0 /go/bin/vanity  /bin/vanity
EXPOSE 8080
USER 2000:2000
CMD ["/bin/vanity"]
