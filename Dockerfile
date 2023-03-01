FROM cgr.dev/chainguard/go:1.19
WORKDIR /app
COPY src/ /app/
RUN make build

FROM cgr.dev/chainguard/glibc-dynamic
COPY --from=0 /app/metronomikon /bin/
COPY example/config.yaml /etc/metronomikon/config.yaml
ENTRYPOINT ["metronomikon"]
