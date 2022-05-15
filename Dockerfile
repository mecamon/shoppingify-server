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
WORKDIR /go/bin
RUN apk update && apk add bash