# Build the application
FROM golang:1-alpine as build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /opt/legolas/legolas
RUN rm -rf /src

# We need updated CAs
COPY stg/letsencrypt-stg-root-x1.pem /etc/ssl/certs/letsencrypt-stg-root-x1.pem
COPY stg/letsencrypt-stg-root-x2.pem /etc/ssl/certs/letsencrypt-stg-root-x2.pem
RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /data
ENTRYPOINT ["/opt/legolas/legolas"]
