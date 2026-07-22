# Review Fixture Protocol

This repository validates IgniteDX against reproducible, evidence-backed code
review defects. The clean implementations and oracle tests live on `develop`.
Intentional defects live only on scenario branches.

## Branches

Create every scenario branch from `develop`:

- `fixture/go/ignored-errors-lost-context`
- `fixture/go/context-cancellation`
- `fixture/go/goroutine-channel-leak`
- `fixture/go/nil-map-slice-hazards`

Use one commit per independent defect. A scenario pull request must continue to
compile and pass `go vet ./...` and `go test ./...`. Its matching oracle test is
expected to fail when run with `go test -tags fixture_oracle ./...`.

For strictness comparisons, create relaxed, balanced, and strict run branches
at the same scenario commit. This ensures each mode receives an identical diff.

## Full Review

1. Open the scenario pull request against `develop`.
2. Run a full review at the bug commit.
3. Record the pull request, review ID, head SHA, strictness, finding IDs,
   severities, and source locations.
4. Compare the result with `expectations/go.yaml`.
5. Reject duplicate active findings that represent the same expectation ID.

## Incremental Reconciliation

1. Begin with two independently fixable findings, A and B, confirmed by a full
   review.
2. Push a commit that fixes only A, then run an incremental review.
3. Confirm A is reconciled and B remains active exactly once.
4. Push an unrelated valid change and run another incremental review.
5. Confirm B remains active exactly once and is not cloned into the new review.
6. Push the fix for B and run an incremental review.
7. Confirm no actionable findings remain and the pull request is approved,
   unless a legitimate suggestion remains.
8. Run a final full review at the fixed head SHA. Resolved findings must not be
   recreated or resurrected.
9. Compare the UI with persisted review comments. Findings must be sourced from
   review-comment records, and each active expectation ID must have at most one
   active occurrence across the pull request.

## Result Rules

- Concrete runtime failures remain actionable in every strictness mode.
- Relaxed mode may keep advisory ownership or defensive-copy concerns as
  suggestions, but must not suppress proven panics, wrong results, blocked
  goroutines, or ignored cancellation.
- Balanced mode should report demonstrated defects and omit stylistic noise.
- Strict mode may promote evidence-backed ownership and lifecycle risks, but
  must not invent repository conventions.
- Finding titles may vary. Identity is determined by root cause, affected code,
  and the stable expectation ID in the manifest.
