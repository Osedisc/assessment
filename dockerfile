FROM golang:1.19-alpine as buil-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v

RUN go build -o ./out/go-exam .

# ====================

FROM alpine:3.16.2
COPY --from=buil-base /app/out/go-exam /app/go-exam

CMD ["/app/go-exam"]