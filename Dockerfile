FROM golang:1.19-alpine AS cno-api
WORKDIR /app
COPY go.mod ./
RUN go mod tidy
COPY . ./
RUN apk add build-base
RUN go build -o cno-api

FROM alpine:3.16.2
WORKDIR /app
COPY --from=cno-api /app/cno-api .
EXPOSE 8080

CMD [ "./cno-api" ]