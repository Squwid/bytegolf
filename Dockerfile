FROM golang:1.19.0-alpine3.16 as build
ADD . /bytegolf
WORKDIR /bytegolf

RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bytegolf-backend

FROM alpine:latest
COPY --from=build /bytegolf/bytegolf-backend .

ARG ENV=staging
ARG FRONTEND_URL=https://staging.byte.golf
ARG BACKEND_URL=https://staging.api.byte.golf

ENV GCP_PROJECT_ID=squid-cloud
ENV BG_ENV=${ENV}
ENV BG_FRONTEND_ADDR=${FRONTEND_URL}
ENV BG_BACKEND_ADDR=${BACKEND_URL}
ENV BG_COOKIE_NAME=bg-token-${ENV}

CMD ["./bytegolf-backend"]