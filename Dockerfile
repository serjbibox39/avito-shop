FROM golang:1.22 AS build_stage
RUN mkdir -p go/src/avito-shop
WORKDIR /go/src/avito-shop
COPY ./ ./
RUN go env -w GO111MODULE=auto && go install ./cmd
WORKDIR /
RUN cp go/src/avito-shop/cmd/config.json go/bin

FROM ubuntu:22.04
RUN mkdir -p avito-shop
WORKDIR /avito-shop
COPY --from=build_stage /go/bin .
RUN mv cmd avito-shop
ENTRYPOINT ./avito-shop
EXPOSE 8080