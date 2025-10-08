FROM golang:1.25.2 AS build

WORKDIR /go/src/app
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    set -eux; \
    CGO_ENABLED=0 GO111MODULE=on go install ./cmd/distribution; \
    go run github.com/google/go-licenses@latest save ./... --save_path=/notices;

FROM ghcr.io/greboid/dockerbase/nonroot:1.20250803.0

COPY --from=build /go/bin/distribution /irc-distribution
COPY --from=build /notices /notices
CMD ["/irc-distribution"]
