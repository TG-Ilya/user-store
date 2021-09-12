FROM golang:1.16-alpine AS builder
RUN apk add build-base &&\
    apk add sqlite-dev
WORKDIR /go/src/user-store/
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -tags "libsqlite3 linux" -a -installsuffix cgo -o user-store ./cmd/user-store/.

FROM alpine:latest
RUN apk add build-base &&\
    apk add sqlite-dev
WORKDIR /application/
COPY ./migrations/ ./migrations/
COPY --from=builder /go/src/user-store/user-store ./
CMD ["./user-store"]