FROM golang:1.23.4-alpine3.20 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ms ./cmd/server/main.go

FROM scratch
WORKDIR /app
COPY ./cmd/server/.env .
COPY --from=builder /app/ms .
CMD [ "./ms" ]