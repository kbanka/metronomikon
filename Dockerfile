FROM golang:1.18
WORKDIR /app
COPY src/ /app/
RUN make build

FROM alpine:3.10
COPY --from=0 /app/metronomikon /bin/
COPY example/config.yaml /etc/metronomikon/config.yaml
ENTRYPOINT ["metronomikon"]
