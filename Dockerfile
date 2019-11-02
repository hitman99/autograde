FROM golang:1.13

COPY . /gosrc
WORKDIR /gosrc
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o autograde

FROM node:12-alpine
COPY frontend /src
WORKDIR /src

RUN npm install && npm run build

FROM alpine:3.9
LABEL maintainer="tomas@adomavicius.com"

RUN apk --no-cache add ca-certificates
WORKDIR /autograde
COPY --from=0 /gosrc/autograde autograde
COPY --from=1 /src/dist frontend/dist
ENV PATH="/autograde/:${PATH}"

CMD ["autograde"]