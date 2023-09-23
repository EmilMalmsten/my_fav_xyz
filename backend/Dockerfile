FROM golang:1.20-alpine AS build-backend
WORKDIR /app
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY . .
RUN go build -o out

FROM alpine:3.14

COPY --from=build-backend /go/bin/migrate /usr/local/bin/
COPY --from=build-backend /app/out /app/
COPY ./sql/schema /migrations

CMD ["/app/out"]
