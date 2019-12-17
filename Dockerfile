FROM golang:1.12
WORKDIR /app
COPY src/ /app/
RUN make build

FROM alpine:3.10
WORKDIR /app
COPY --from=0 /app/metronomikon /app/
COPY example/config.yaml /etc/metronomikon/config.yaml
ENTRYPOINT ["/app/metronomikon"]
