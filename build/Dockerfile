FROM golang:1.21.3-bullseye as build

WORKDIR /build

ARG SERVICE

COPY . .
RUN CGO_ENABLED=0 go build -o /out/$SERVICE ./cmd/$SERVICE/main.go

FROM alpine

ARG SERVICE

WORKDIR /app
COPY --from=build /out/$SERVICE /app/service

ENTRYPOINT [ "./service" ]
