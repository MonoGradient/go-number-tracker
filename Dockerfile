FROM golang:1.14.6-alpine3.12 as builder
EXPOSE 8080
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .
RUN ls -lrt

FROM scratch
WORKDIR /app
COPY --from=builder /app/main /app
CMD ["./main"]