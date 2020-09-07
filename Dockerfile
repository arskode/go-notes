FROM golang:1.15

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

EXPOSE 5000

ENTRYPOINT ./entrypoint.dev.sh