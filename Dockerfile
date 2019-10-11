FROM golang:latest AS builder
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /translate .

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /translate ./
COPY --from=builder /app/ui ./ui
#RUN chmod +x ./translate
ENTRYPOINT ["./translate"]
EXPOSE 8080