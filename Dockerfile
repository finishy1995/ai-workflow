FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /build/zero

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY /data/ /app/data
COPY /ai/etc/ /app/etc
RUN go build -ldflags="-s -w" -o /app/ai ai/ai.go


FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/ai /app/ai
COPY --from=builder /app/etc /app/etc
COPY --from=builder /app/data /data
RUN mkdir /app/tmp

EXPOSE 8081

CMD ["./ai", "-f", "etc/ai.yaml"]
