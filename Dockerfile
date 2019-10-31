FROM golang:1.13

COPY . /gosrc
WORKDIR /gosrc
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o autograde

FROM alpine:3.9
LABEL maintainer="tomas@adomavicius.com"

RUN apk --no-cache add ca-certificates
WORKDIR /autograde
COPY --from=0 /gosrc/autograde autograde
ENV PATH="/autograde/:${PATH}"

CMD ["autograde"]