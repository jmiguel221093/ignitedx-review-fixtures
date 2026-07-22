# Go Review Fixture

This fixture is a standard-library HTTP API and a clean reference baseline for
Go reviews. It deliberately keeps the domain small so review outcomes are about
language semantics rather than framework conventions.

The baseline demonstrates:

- propagation of request contexts into downstream HTTP calls;
- cancellation-aware retry waits and graceful server shutdown;
- checked transport, status, decoding, and response-writing errors;
- producer-owned channel closure and cancellation-safe worker sends;
- initialized maps and guarded map and slice access; and
- defensive copies for mutable slices crossing store boundaries.

## Run

```sh
go run ./cmd/fixture-api
```

The server listens on `:8080` by default. Set `LISTEN_ADDR` to change the
address and `PROFILE_UPSTREAM_URL` to point profile requests at another HTTP
service.

## Validate

```sh
gofmt -l .
go vet ./...
go test ./...
go test -tags fixture_oracle ./...
```

Ordinary tests describe supported behavior. Oracle tests are protected by the
`fixture_oracle` build tag and correspond to stable finding IDs in
`../.review-fixtures/expectations/go.yaml`. Intentional scenario branches should
continue to pass ordinary checks while causing only their matching oracle tests
to fail.
