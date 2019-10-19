# download modules
FROM golang:1.13 AS modules

COPY go.mod go.sum /
RUN go mod download

# build the binary
FROM golang:1.13 AS build

COPY --from=modules /go/pkg /go/pkg

RUN useradd -u 10001 gopher

WORKDIR /go/src/goss

COPY . .

RUN GOOS=linux GOARCH=amd64 make build

# run the binary
FROM scratch

ENV PORT 8080

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd

USER gopher

COPY --from=build /go/src/goss/migrations /migrations
COPY --from=build /go/src/goss/bin/goss /goss

EXPOSE $PORT
CMD ["/goss"]
