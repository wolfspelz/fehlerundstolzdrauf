FROM golang:1.22-alpine AS build
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o server ./cmd/server

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build /app/server /server
COPY --from=build /app/public /public
COPY --from=build /app/internal/db/seed.sql /seed.sql
EXPOSE 80
CMD ["/server"]
