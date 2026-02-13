FROM golang:1.25 AS build

WORKDIR /app/src
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./mockserver ./mockserver
WORKDIR /app/src/mockserver
RUN CGO_ENABLED=0 go build -o /app/bin/mockserver

FROM gcr.io/distroless/static-debian12

ARG COMMIT_SHA=local
LABEL org.opencontainers.image.source="https://github.com/cape-project-eu/cape-sdk-generator/blob/${COMMIT_SHA}/mockserver/Dockerfile" \
    org.opencontainers.image.authors="University of St. Gallen" \
    org.opencontainers.image.revision="${COMMIT_SHA}" \
    org.opencontainers.image.title="CAPE API Mockserver"

ENV GIN_MODE=release
COPY --from=build /app/bin/mockserver /
CMD ["/mockserver"]
