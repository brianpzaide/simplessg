FROM golang:1.25.4-alpine AS builder
RUN apk update && apk add make git
WORKDIR /home/simplessg
RUN go install github.com/knadh/stuffbin/...@latest
COPY . .
RUN go build -o simplessg .

FROM golang:1.25.4-alpine AS deploy
RUN apk --no-cache add ca-certificates
WORKDIR /home/simplessg
COPY --from=builder /home/simplessg/simplessg .
COPY templates/ .

CMD ["./simplessg", "build"]