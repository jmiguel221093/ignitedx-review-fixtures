# IgniteDX Review Fixtures

Public, runnable projects used to validate IgniteDX reviews across languages and frameworks.

Each fixture has a clean baseline on `develop`. Intentional defects live only on dedicated scenario branches and pull requests, together with an expectation manifest that documents the findings each review mode should produce.

## Fixtures

- [`react/`](react/README.md) - React, TypeScript, Vite, and Tailwind CSS
- [`go/`](go/README.md) - Go HTTP API covering context, error, concurrency, and collection safety

### React migration

The React fixture is a snapshot imported from
`jmiguel221093/ignitedx-demo-repo` at `origin/develop` (`c1856c8`). The import
received only the minimal corrections required to establish a clean lint
baseline:

- a stable heartbeat state initializer;
- a deferred debug-preset state update; and
- a private `buttonVariants` export.

The source repository remains unchanged so its historical review pull requests
continue to be available.

## Validation

Each fixture remains runnable on `develop`. Standard checks verify that the
baseline builds and behaves correctly, while fixture-oracle tests capture the
specific regressions that intentional scenario branches introduce.

```sh
cd react && npm ci && npm run lint && npm run build
cd go && gofmt -l . && go vet ./... && go test ./...
cd go && go test -tags fixture_oracle ./...
```

The expected IgniteDX findings and review protocol live in
[`.review-fixtures/`](.review-fixtures/protocol.md). These files are test
metadata, not application context for the reviewer.

## Branch policy

All implementation and scenario branches start from `develop`. Intentional defects must not be merged into `develop`.
