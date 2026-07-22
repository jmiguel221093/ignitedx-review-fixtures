# IgniteDX Review Fixtures

Public, runnable projects used to validate IgniteDX reviews across languages and frameworks.

Each fixture has a clean baseline on `develop`. Intentional defects live only on dedicated scenario branches and pull requests, together with an expectation manifest that documents the findings each review mode should produce.

## Fixtures

- [`react/`](react/README.md) - React, TypeScript, Vite, and Tailwind CSS

## Branch policy

All implementation and scenario branches start from `develop`. Intentional defects must not be merged into `develop`.

