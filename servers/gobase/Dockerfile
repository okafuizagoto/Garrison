FROM golang:1.24-alpine

ARG APP_NAME
ENV APP_NAME=$APP_NAME
ENV GONOSUMDB="*"
ENV GOFLAGS="-mod=mod"

RUN mkdir -p /app/$APP_NAME /pkg

RUN apk --no-cache add bash curl git gcc g++ inotify-tools

WORKDIR /app/$APP_NAME

RUN mkdir -p /go/_tools
ADD config/build.sh /go/_tools/build.sh
ADD config/reflex.conf /go/_tools/reflex.conf
RUN chmod +x /go/_tools/build.sh

RUN go install github.com/cespare/reflex@latest

CMD ["reflex", "-c", "/go/_tools/reflex.conf"]
