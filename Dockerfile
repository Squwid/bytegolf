FROM golang:1.12.4-alpine3.9 AS build
RUN apk add --no-cache git
ADD . /bytegolf
WORKDIR /bytegolf
RUN go mod vendor
RUN go build

FROM alpine:latest
COPY --from=build /bytegolf .

# Set the environmental variables for no panic
ENV db_username null
ENV prod true

CMD ["./bytegolf"]
