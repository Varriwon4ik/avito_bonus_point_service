FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bonus-ledger ./cmd/api

FROM alpine:3.19
COPY --from=build /bonus-ledger /bonus-ledger
EXPOSE 8080
ENTRYPOINT ["/bonus-ledger"]
