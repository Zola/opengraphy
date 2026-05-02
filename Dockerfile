FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/opengraphy ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=build /out/opengraphy /app/opengraphy
COPY internal/templates /app/internal/templates
COPY web/static /app/web/static
EXPOSE 8080
CMD ["/app/opengraphy"]
