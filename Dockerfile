FROM golang:1.13.0-alpine3.10 AS build
RUN apk add --no-cache git
ADD . /bytegolf
WORKDIR /bytegolf
RUN go mod vendor
RUN go build

FROM alpine:latest
COPY --from=build /bytegolf .

# Set the environmental variables for no panic
ENV db_username null
ENV prod false

ENV PROJECT_ID bytegolf

CMD ["./bytegolf"]
