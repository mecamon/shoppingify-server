# syntax = docker/dockerfile:1.2

FROM golang:1.17-alpine AS build
# Support CGO and SSL
RUN apk --no-cache add gcc g++ make
RUN apk add git
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/shop-app ./cmd/web/*.go

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build /usr/src/app/bin /go/bin
RUN --mount=type=secret,id=_env,dst=/etc/secrets/.env cat /etc/secrets/.env
EXPOSE 8080
ENTRYPOINT ["/go/bin/shop-app", "-db-host=$DB_HOST" , "-db-user=$DB_USER", "-db-password=$DB_PASSWORD", "-db-name=$DB_NAME", "-storage-cloud-name=$STORAGE_CLOUD_NAME", "-storage-api-key=$STORAGE_API_KEY", "-storage-api-secret=$STORAGE_API_SECRET", "-storage-api-env-var=$STORAGE_API_ENV_VAR"]