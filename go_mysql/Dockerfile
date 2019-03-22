FROM golang:1.11.2-alpine3.8 AS build

WORKDIR /go/src/go_mysql
COPY . .
RUN apk update \
  && apk add --no-cache git \
  && go get github.com/go-sql-driver/mysql
RUN go install -v go_mysql

FROM alpine:3.8
COPY --from=build /go/bin/go_mysql /usr/local/bin/
CMD ["go_mysql"]