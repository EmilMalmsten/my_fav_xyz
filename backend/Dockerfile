FROM golang:1.20 AS build-backend
WORKDIR /app
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY . .
RUN go build -o out

FROM gcr.io/distroless/base-debian12

COPY --from=build-backend /go/bin/migrate /usr/local/bin/
COPY --from=build-backend /app/internal/database/images /app/internal/database/images
COPY --from=build-backend /app/out /app/

CMD ["/app/out"]
