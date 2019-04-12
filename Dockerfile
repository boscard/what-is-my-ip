FROM golang:1.11 AS go_builder
RUN mkdir -p /go/src/whatismyip
COPY ./ /go/src/whatismyip/
RUN cd /go/src/whatismyip && go test && go build 

FROM debian:stable-slim
COPY --from=go_builder /go/src/whatismyip/whatismyip /bin/whatismyip
EXPOSE 8080
CMD ["/bin/whatismyip"]
