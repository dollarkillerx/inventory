FROM alpine:latest
WORKDIR /app
COPY ./cmd/inventory /app
COPY ./configs/ /app/configs
COPY ./stats/ /app/static
CMD ["/app/inventory"]