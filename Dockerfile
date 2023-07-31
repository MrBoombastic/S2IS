FROM golang:1.20-alpine as builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY . ./
ARG ref_name
ENV ref_name ${ref_name}
RUN go build -o /app/s2fs main.go

FROM alpine:latest as worker
COPY --from=builder /app/s2fs .

EXPOSE 3000

ENTRYPOINT ["/s2fs"]