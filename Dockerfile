FROM golang:1.16.5-alpine3.14 as build
ADD . /bytegolf
WORKDIR /bytegolf

RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bytegolf-backend

FROM alpine:latest
COPY --from=build /bytegolf/bytegolf-backend .

ARG ENV=dev

ENV GCP_PROJECT_ID=squid-cloud
ENV BG_ENV=${ENV}
ENV BG_FRONTEND_ADDR=https://dev.byte.golf
ENV BG_BACKEND_ADDR=https://dev-api.byte.golf

CMD ["./bytegolf-backend"]