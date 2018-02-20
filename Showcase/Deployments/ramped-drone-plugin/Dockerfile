FROM alpine
COPY ["config", "ca.crt", "client.crt", "client.key", "./update", "./deployment.yaml", "update.yaml", "/"]
ENTRYPOINT /update
