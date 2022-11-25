# Build the application
FROM golang:1 as build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /legolas

# We need updated CAs
RUN apt update -qqq && apt install -yqqq ca-certificates && update-ca-certificates

# Copy over to a small container to minimize footprint and attack surface
FROM gcr.io/distroless/static:nonroot

COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /legolas /home/nonroot/

ENTRYPOINT ["/home/nonroot/legolas"]
