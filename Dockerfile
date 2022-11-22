FROM golang:1 as build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /legolas

FROM gcr.io/distroless/static:nonroot

COPY --from=build /legolas /
CMD ["/legolas"]
