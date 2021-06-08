FROM golang:1.16.4-alpine3.10 AS build
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
