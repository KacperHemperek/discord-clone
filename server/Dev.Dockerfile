FROM golang:1.22

WORKDIR /app

RUN apt-get update && apt-get install -y postgresql-client

RUN go install github.com/cosmtrek/air@latest

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY . .

HEALTHCHECK --timeout=30s --retries=3 \
 CMD pg_isready -h db -p 5432 -d discord -U postgres || exit 1

CMD ["air", "-c", ".air.toml"]