FROM alpine

ADD vanity /usr/local/bin

EXPOSE 8080

ENTRYPOINT ["/bin/sh", "-c"]

CMD ["vanity"]
