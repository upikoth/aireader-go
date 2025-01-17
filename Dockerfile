# Stage 1. Build.

FROM golang:1.23.5-alpine as build

RUN apk add --no-cache make \
	&& rm -rf /var/cache/apk/* /tmp/* /var/tmp/*

WORKDIR /workdir

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN make build

# Stage 2.

FROM alpine:3.21

COPY --from=build /workdir/app ./
COPY --from=build /workdir/docs ./docs

EXPOSE 8888

CMD ["./app"]
