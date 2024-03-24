# Build stage
FROM golang:1.22-alpine as build

WORKDIR /var/app

COPY . .

RUN go mod download
RUN go build -o ./build/goods_crud ./cmd/goods_crud/main.go


# Run stage
FROM alpine:3

WORKDIR /

COPY --from=build /var/app/build /app

EXPOSE 20001

ENTRYPOINT ["/app/goods_crud"]
