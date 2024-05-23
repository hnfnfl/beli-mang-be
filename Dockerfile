FROM alpine:latest
WORKDIR /app
COPY halo-suster-be /app
COPY local_configuration/config.yaml /app/local_configuration/
EXPOSE 8080
CMD ["./halo-suster-be"]