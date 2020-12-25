FROM golang:1.15.6-alpine3.12 as b
WORKDIR /app
COPY . .
RUN go build -o ./app ./cmd

FROM alpine:3.12.3
WORKDIR /app
RUN apk --no-cache add tzdata=2020c-r1
COPY --from=b /app/app .
ENTRYPOINT [ "./app" ]
