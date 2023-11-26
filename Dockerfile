FROM golang:1.21 as build

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOOS=linux

COPY . .

#todo: RUN go build  -o bin/limiter-migrator ./cmd/migrations/
#todo: RUN bin/limiter-migrator -config=./configs/config.yml up

RUN go build -o bin/limiter ./cmd/api-rate-limiter/

FROM scratch

EXPOSE 4242

COPY --from=build /app/bin/limiter .
COPY --from=build /app/configs/config.yml .

CMD ["/limiter", "-config=./config.yml"]

