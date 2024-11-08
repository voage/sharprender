FROM golang:1.23-alpine AS builder

RUN apk update && \
    apk add --no-cache gcc g++ pkgconfig vips-dev git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd

FROM alpine:latest

RUN apk add --no-cache vips

RUN apk update && \
    apk add --no-cache chromium nss freetype harfbuzz ttf-freefont

WORKDIR /app

COPY --from=builder /app/main .

ENV CHROME_PATH=/usr/bin/chromium-browser

EXPOSE 8080

CMD ["./main"]
