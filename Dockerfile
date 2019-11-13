FROM golang:1.12
WORKDIR /app
COPY . /app/
RUN make build

FROM alpine:3.10
WORKDIR /app
COPY --from=0 /app/metronomikon /app/
ENTRYPOINT ["/app/metronomikon"]
