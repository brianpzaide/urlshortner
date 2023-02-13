FROM golang:1.19-alpine AS builder
RUN apk update && apk add make git
WORKDIR /urlshortner
COPY ./ ./
ENV CGO_ENABLED=0 GOOS=linux
RUN make build

FROM golang:1.19-alpine AS deploy
RUN apk --no-cache add ca-certificates
COPY --from=builder /urlshortner/urlshortner ./
COPY start.sh ./
CMD ["./start.sh"]