FROM alpine:latest
WORKDIR /app
COPY beli-mang-be /app
COPY local_configuration/config.yaml /app/local_configuration/
EXPOSE 8080
CMD ["./beli-mang-be"]