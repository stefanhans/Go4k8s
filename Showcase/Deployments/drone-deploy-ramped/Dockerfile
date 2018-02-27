FROM alpine
COPY ["config", "ca.crt", "client.crt", "client.key", "./deploy", "./deployment.yaml", "/"]
ENTRYPOINT /deploy
