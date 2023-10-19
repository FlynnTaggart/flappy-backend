FROM golang:1.21-alpine as builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /main ./cmd/service/main.go


FROM alpine:3
COPY --from=builder main /bin/main
ENTRYPOINT ["/bin/main"]