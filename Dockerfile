FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /out/app .


FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S app && adduser -S -G app app \
	&& update-ca-certificates

WORKDIR /app

COPY --from=builder /out/app /app/app

EXPOSE 8080


USER app

ENTRYPOINT ["/app/app"]
