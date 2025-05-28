FROM golang:1.23 AS build-stage

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd/datum-controller-manager/ cmd/datum-controller-manager/
COPY internal/ internal/
COPY pkg/ pkg/

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=cache,target="/root/.cache/go-build" \
  CGO_ENABLED=0 GOOS=linux \
  go build \
  -ldflags="-X 'k8s.io/component-base/version/verflag.programName=Datum'" \
  -o /datum-controller-manager cmd/datum-controller-manager/main.go

FROM gcr.io/distroless/base-debian11
ARG cmd

WORKDIR /

COPY --from=build-stage /datum-controller-manager /datum-controller-manager

ENTRYPOINT ["/datum-controller-manager"]
