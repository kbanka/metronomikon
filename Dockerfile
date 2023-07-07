# go:1.19 on 20230706
FROM cgr.dev/chainguard/go@sha256:13a12452e39525bf71ca9bee362fcaccd933952960391a601676e55406b6fc2f AS build

WORKDIR /app
COPY src/ /app/
RUN make build

FROM cgr.dev/chainguard/glibc-dynamic AS metronomikon
COPY --from=0 /app/metronomikon /bin/
COPY example/config.yaml /etc/metronomikon/config.yaml
ENTRYPOINT ["metronomikon"]
