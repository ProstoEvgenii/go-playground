FROM golang:alpine AS BUILD
WORKDIR /app
COPY . .
RUN go build -o rest-service

FROM alpine:latest
COPY --from=BUILD /app/rest-service .
CMD ["./rest-service"]