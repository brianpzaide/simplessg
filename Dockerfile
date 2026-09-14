FROM golang:1.25.4-alpine AS builder
RUN apk update && apk add make git
WORKDIR /home/simplessg
COPY . .
RUN go build -o simplessg .

FROM golang:1.25.4-alpine AS deploy
RUN apk --no-cache add ca-certificates
COPY --from=builder /home/simplessg/simplessg /usr/bin/
WORKDIR /home/simplessg
RUN mkdir templates posts 
COPY templates/ templates/


CMD ["simplessg", "build"]