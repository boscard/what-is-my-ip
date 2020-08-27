FROM golang:1.13 AS go_builder
RUN mkdir -p /go/src/github.com/boscard/what-is-my-ip
COPY ./ /go/src/github.com/boscard/what-is-my-ip
RUN cd /go/src/github.com/boscard/what-is-my-ip && go build

FROM debian:stable-slim
COPY --from=go_builder /go/src/github.com/boscard/what-is-my-ip/what-is-my-ip /bin/what-is-my-ip
EXPOSE 8080
CMD ["/bin/what-is-my-ip"]
