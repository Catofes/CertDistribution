FROM golang as builder
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN go mod download
RUN env CGO_ENABLED=0 go build -o /web -ldflags "-X 'main._version_=$(git log --pretty=format:"%h" -1)'" .

FROM alpine:3.6
RUN apk add --no-cache tzdata
COPY --from=builder /web /usr/bin/web
ENTRYPOINT ["/usr/bin/web"]
CMD ["-c", "/etc/config.json"]

