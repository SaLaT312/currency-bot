FROM golang:1.19.13-alpine3.18  as builder
WORKDIR /src/app
COPY go.mod *.go ./
RUN go mod tidy
RUN go build -o valute_bot



FROM alpine:3.19
COPY --from=builder /src/app/valute_bot /valute_bot
CMD ["/valute_bot"]