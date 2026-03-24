# Developer Guide

Everything a contributor needs to know to work on `gitconfig`.

---

## Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| Go | ≥ 1.24 | Build & test |
| git | any | Required at runtime for `config` command |
| golangci-lint | latest | Linting (optional but recommended) |
| goreleaser | latest | Building release binaries |

Install Go from https://go.dev/dl/. All other dependencies are pulled via `go mod`.

---

## Getting started

```sh
git clone https://github.com/0ghny/gitconfig
cd gitconfig

# Download dependencies
go mod download

# Build a local binary for your current OS/arch
./hack/dev.sh build

# Run the full test suite
./hack/dev.sh test
```

---

## Development script — `hack/dev.sh`

All common development tasks are available through `hack/dev.sh`. Run it without arguments (requires `fzf`) or pass a task name directly:

```sh
./hack/dev.sh build          # compile for current OS/arch → target/<timestamp>/gitconfig
./hack/dev.sh test           # go test ./...
./hack/dev.sh test --coverage|-c  # test with coverage report → target/<timestamp>/
./hack/dev.sh lint           # golangci-lint run
./hack/dev.sh fmt            # gofmt -l check
./hack/dev.sh deps           # go mod tidy + list available upgrades
./hack/dev.sh clean          # remove target/ and purge Go build/test/module caches
./hack/dev.sh all            # fmt + lint + test
```

---

## Project structure

See [ARCHITECTURE.md](ARCHITECTURE.md) for a full description of the hexagonal layering. The short version:

```
internal/domain/       ← pure business logic, no I/O
internal/application/  ← use-case orchestration (LocationManager)
internal/adapter/      ← CLI (cobra), filesystem (afero), git process
cmd/                   ← entry point only
```

---

## Adding a new feature

### Adding a new CLI command

1. Create `internal/adapter/cli/<command>.go` following the pattern in `location_new.go`.
2. The function signature must accept `factory LocationServiceFactory` (and any shared flags pointer).
3. Use `cmd.OutOrStdout()` and `cmd.InOrStdin()` — never `fmt.Println` or `os.Stdin` directly.
4. Register the command in `root.go` or as a subcommand of an existing command.
5. Add a matching `<command>_test.go` using `RootCmdWithDeps` + the in-memory `testEnv` helper.

### Extending the Service interface

1. Add the method signature to `internal/domain/location/service.go`.
2. Implement it on `LocationManager` in `internal/application/locations/manager.go`.
3. The compile-time check `var _ location.Service = (*LocationManager)(nil)` will fail until the implementation is complete — that's intentional.
4. Add tests in `manager_test.go` using `afero.NewMemMapFs()`.

---

## Testing

The project uses three kinds of tests:

| Kind | Location | How it works |
|---|---|---|
| Unit | `*_test.go` beside source files | Pure Go, no filesystem, no process |
| Application | `internal/application/locations/manager_test.go` | `afero.MemMapFs` — no real disk |
| CLI integration | `internal/adapter/cli/*_test.go` | `RootCmdWithDeps` + `afero.MemMapFs` — no real disk and no real git |

```sh
# Run everything
go test ./...

# Run a specific package
go test ./internal/adapter/cli/...

# Run with coverage
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

Tests that require the `git` binary are guarded with:

```go
func skipIfGitNotFound(t *testing.T) {
    t.Helper()
    if _, err := exec.LookPath("git"); err != nil {
        t.Skip("git binary not found in PATH")
    }
}
```

---

## Code style

- Follow standard Go idioms. Run `gofmt` before committing.
- All exported symbols must have a godoc comment.
- Unexported helpers used only in tests should live in `*_test.go` files.
- Errors should be propagated with `fmt.Errorf("context: %w", err)` where additional context is useful.
- Avoid `panic` outside `init()` and `main()`.

---

## Commit & branch conventions

- Branch from `main`.
- Use descriptive branch names: `feature/<name>`, `fix/<name>`, `chore/<name>`.
- Write conventional commit messages: `feat:`, `fix:`, `docs:`, `chore:`, `test:`, `refactor:`.
- Open a Pull Request against `main`; CI must pass before merging.

---

## Releasing

See [RELEASE.md](RELEASE.md).
