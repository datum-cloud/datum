FROM golang:1.23 AS build-stage

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd/datum-apiserver/ cmd/datum-apiserver/
COPY pkg pkg/

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=cache,target="/root/.cache/go-build" \
  CGO_ENABLED=0 GOOS=linux \
  go build \
  -ldflags="-X 'k8s.io/component-base/version/verflag.programName=Datum'" \
  -o /datum-apiserver cmd/datum-apiserver/main.go

FROM gcr.io/distroless/base-debian11
ARG cmd

WORKDIR /

COPY --from=build-stage /datum-apiserver /datum-apiserver

ENTRYPOINT ["/datum-apiserver"]
