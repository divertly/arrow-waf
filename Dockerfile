#### BUILDER
FROM golang:1.25-alpine3.23 AS builder

WORKDIR /dist

ARG VERSION=0.0.0-local

COPY . .
RUN go env && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags "-X main.VERSION=${VERSION} -s -w" -o waf .



### RELEASE IMAGE #########################################
FROM alpine:3.23
LABEL org.opencontainers.image.source=https://github.com/divertly/waf
LABEL org.opencontainers.image.description="Divertly WAF; securing stuff."

ENV \
    APP_PROFILE=docker 

COPY --from=builder /dist/waf .
COPY --from=builder /dist/profiles ./profiles/
COPY --from=builder /dist/conf ./conf/

# public port
EXPOSE 6080

# private port
EXPOSE 6081

USER 1001:1001

ENTRYPOINT ["./waf"]