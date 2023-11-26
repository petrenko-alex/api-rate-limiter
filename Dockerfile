FROM golang:1.21
# todo: decrease size (two step file)

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOOS=linux

COPY . .

#todo: RUN go build  -o bin/limiter-migrator ./cmd/migrations/
#todo: RUN bin/limiter-migrator -config=./configs/config.yml up

RUN go build -o bin/limiter ./cmd/api-rate-limiter/

CMD ["bin/limiter", "-config=./configs/config.yml"]

EXPOSE 4242
