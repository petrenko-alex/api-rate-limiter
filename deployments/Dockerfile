FROM golang:1.21

WORKDIR /app

COPY ../ .

# run mingrations
# decrease size (two step file)

RUN go build -o bin/limiter ./cmd/api-rate-limiter/

CMD ["bin/limiter"]

