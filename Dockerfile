FROM golang:latest AS builder
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /translate .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /translate ./
COPY --from=builder /app/ui ./ui
ENTRYPOINT ["./translate"]
EXPOSE 8080