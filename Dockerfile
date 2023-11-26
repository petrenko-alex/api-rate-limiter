FROM golang:1.21 as build

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOOS=linux

COPY . .

RUN go build -o bin/limiter ./cmd/api-rate-limiter/
RUN go build -o bin/limiter-migrator ./cmd/migrations/

FROM busybox

EXPOSE 4242

COPY --from=build /app/bin/limiter .
COPY --from=build /app/bin/limiter-migrator .
COPY --from=build /app/migrations ./migrations/
COPY --from=build /app/configs/config.yml .

CMD ["sh", "-c", "/limiter-migrator -config=./config.yml up && /limiter -config=./config.yml"]

