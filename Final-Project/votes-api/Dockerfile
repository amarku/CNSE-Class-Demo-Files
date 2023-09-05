FROM golang:1.20 AS build-stage

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /votes-api


FROM alpine:latest AS run-stage

WORKDIR /

COPY --from=build-stage /votes-api /votes-api

EXPOSE 3080

ENV REDIS_URL=host.docker.internal:6379

CMD ["/votes-api"]