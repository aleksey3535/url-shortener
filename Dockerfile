FROM golang:1.23.1

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/url-shortener/main.go


CMD ["app", "--config=./config/local.yaml"]

