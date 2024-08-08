# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o /app/gorg /app/cmd/gorg/main.go 

CMD ["/app/gorg"]

FROM alpine:3.20 AS run

WORKDIR /app

RUN apk add --no-cache exiftool

COPY --from=build /app/gorg /app/gorg

# ENTRYPOINT ["/app/gorg"]
