FROM instrumentisto/glide as builder
ENV GOPATH /go
COPY . /go/src/github.com/Catofes/CertDistribution
WORKDIR /go/src/github.com/Catofes/CertDistribution
RUN glide install
RUN mkdir /app
RUN env CGO_ENABLED=0 go build -o /app/web -ldflags "-X 'main._version_=$(git log --pretty=format:"%h" -1)'"\
                      	 github.com/Catofes/CertDistribution

FROM alpine:3.6
RUN apk add --no-cache tzdata
COPY --from=builder /app/web /usr/bin/web
ENTRYPOINT ["/usr/bin/web"]
CMD ["-c", "/etc/config.json"]

