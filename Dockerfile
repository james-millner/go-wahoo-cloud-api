FROM golang:1.22.2-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o run-app cmd/main/main.go

RUN ls -l

CMD ["./run-app"]
