FROM golang:1.22.5 AS build
WORKDIR /go/src/app
ENV CGO_ENABLED=0 GOOS=linux GOPROXY=direct

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY templates ./templates
COPY css ./css
RUN go build -v -ldflags "-s -w" -o /go/bin/app .

FROM scratch
COPY --from=build /go/bin/app /
ENTRYPOINT ["/app"]

