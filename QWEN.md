---
description: >-
  GoMachina: Open-source, Go-idiomatic implementation reference. Guidelines for building
  production-grade libraries in Go, following industry standards.
globs: "*.go,go.*,*.mod"
alwaysApply: false
---

You are an expert in Go (Golang) and open-source library design.

Code Style and Structure
- Write concise, readable, Go-idiomatic code.
- Favor composition and small interfaces over inheritance or large abstractions.
- Keep packages focused and minimal; one responsibility per package.
- Export only necessary identifiers; keep helpers unexported.
- Use `context.Context` for I/O or long-running tasks.

Naming Conventions
- CamelCase for exported, lowerCamelCase for unexported.
- Short, singular, lowercase package names without underscores.
- Constructor functions start with `New`.

Go Usage
- Use explicit interfaces at boundaries; concrete types elsewhere.
- Return concrete errors; wrap with `%w`.
- Avoid `any` in public APIs unless essential.
- Small, intention-revealing interfaces (`Runner`, `Evaluator`).

Error Handling
- No panics for recoverable errors.
- Use `errors.Is` / `errors.As`.
- Provide sentinel errors for common failure cases.

Concurrency
- Use `context` for cancellation and deadlines.
- Avoid global mutable state; ensure concurrency safety.

Observability
- Structured logging with optional user-provided logger.
- Hooks for metrics and tracing; no-op defaults.

Testing
- Table-driven tests; deterministic and small.
- Integration tests for end-to-end workflows.
- Run `go test -race` in CI.

Tooling
- Format with `gofmt` and `goimports`.
- Lint with `golangci-lint`.
- Keep dependencies minimal and audited.

OSS & Industry Standard Practices
- Follow Go 1 compatibility promise.
- Keep public API small and stable; use semver.
- Document all exported code with `godoc` comments.
- Provide runnable `examples/` and `README.md` with quick start.
- Enforce style and testing in CI; block merges on failure.

