FROM golang:1.16.5-alpine3.13 AS build
RUN apk add --no-cache git
ADD . /bytegolf
WORKDIR /bytegolf
RUN go mod vendor
RUN go build

ARG env=abc123

FROM scratch
COPY --from=build /bytegolf .

# Set the environmental variables for no panic
ENV db_username null
ENV prod true

ENV PROJECT_ID bytegolf

CMD ["./bytegolf"]
