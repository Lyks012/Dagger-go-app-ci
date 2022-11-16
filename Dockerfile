FROM golang:1.19-alpine AS cno-api
WORKDIR /app
COPY go.mod ./
RUN go mod tidy
COPY . ./
RUN go build -o beopenmairie

FROM alpine:3.16.2
WORKDIR /app
COPY --from=cno-api /app/beopenmairie .
EXPOSE 8080

CMD [ "./beopenmairie" ]