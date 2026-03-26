# Architecture

`gitconfig` follows the **Hexagonal Architecture** (Ports & Adapters) pattern. The goal is to keep the domain logic completely decoupled from I/O concerns (filesystem, CLI, external processes).

---

## Directory layout

```
gitconfig/
├── cmd/                            # Entry point — wires everything together
│   └── main.go
└── internal/
    ├── config/                     # App-wide constants (e.g. home dir name)
    ├── templates/                  # Default file content templates
    ├── home/                       # Resolve & ensure app home directory
    │
    ├── domain/                     # ── DOMAIN LAYER ──────────────────────
    │   ├── location/               # Location entity, factory, Service port
    │   │   ├── location.go
    │   │   ├── location_test.go
    │   │   └── service.go          # ← Port (interface)
    │   └── gitconfig/              # GitConfig value object (read/write .gitconfig)
    │       ├── gitconfig.go
    │       └── gitconfig_test.go
    │
    ├── application/                # ── APPLICATION LAYER ─────────────────
    │   └── locations/
    │       ├── manager.go          # LocationManager — implements location.Service
    │       └── manager_test.go
    │
    └── adapter/                    # ── ADAPTER LAYER (Infrastructure) ────
        ├── cli/                    # Primary adapter: Cobra CLI
        │   ├── root.go
        │   ├── locations.go
        │   ├── location.go
        │   ├── location_new.go
        │   ├── location_delete.go
        │   └── config.go
        ├── filesystem/             # afero wrapper (OS + in-memory)
        │   ├── filesystem.go
        │   ├── file.go
        │   └── wd.go
        └── git/                    # Secondary adapter: git process execution
            ├── git.go
            └── commander/
                ├── commander.go    # Run processes, capture output
                └── builder/
                    └── builder.go  # Fluent exec.Cmd builder
```

---

## Layers

### Domain layer (`internal/domain/`)

Contains pure business logic with **no** I/O dependencies.

| Package | Responsibility |
|---|---|
| `location` | `Location` entity, `NewLocation` factory, `Service` interface (port) |
| `gitconfig` | `GitConfig` struct — read/write a single `.gitconfig` file via an injected `afero.Afero` |

The `location.Service` interface is the **primary port**. It defines the operations the application exposes to the outside world:

```go
type Service interface {
    GetLocations() ([]Location, error)
    FindLocationByKey(key string) (*Location, error)
    SaveLocation(key string, location string) error
    DeleteLocation(key string) error
}
```

### Application layer (`internal/application/`)

Orchestrates domain objects to fulfil use-cases. `LocationManager` implements `location.Service` — it coordinates `GitConfig` (reads/writes `.gitconfig`) and `afero.Afero` (manages individual location config files).

### Adapter layer (`internal/adapter/`)

| Adapter | Direction | Responsibility |
|---|---|---|
| `cli` | Primary (driving) | Cobra command tree — translates CLI input into `location.Service` calls |
| `filesystem` | Secondary (driven) | afero factory helpers — used by domain & application layers |
| `git` | Secondary (driven) | Wraps the `git config` binary via `commander` and `builder` |

### Entry point (`cmd/main.go`)

Calls `cli.RootCmd(version)` which internally wires `defaultLocationServiceFactory` → `LocationManager` → `GitConfig` + OS filesystem.

---

## Dependency injection & testability

The CLI adapter receives a `LocationServiceFactory` function:

```go
type LocationServiceFactory func(gitconfigPath string) location.Service
```

In production, `defaultLocationServiceFactory` returns a real `LocationManager` backed by the OS filesystem. In tests, `RootCmdWithDeps` accepts a factory that returns a `LocationManager` backed by an `afero.MemMapFs`.

This means the entire CLI can be exercised in memory without touching disk, spawning processes, or requiring a real `~/.gitconfig`.

---

## Data flow — example: `gitconfig location new --key work`

```
CLI args
  └─► cobra (adapter/cli/location_new.go)
        └─► LocationServiceFactory(gitconfigPath) → LocationManager
              ├─► GitConfig.GetContent()       reads .gitconfig (afero)
              ├─► LocationManager.SaveLocation()
              │     ├─► location.NewLocation()  builds entity
              │     ├─► location.ToSection()    renders template
              │     └─► GitConfig.AppendSection() / WriteContent()
              └─► GitConfig (writes updated content back to .gitconfig)
```

---

## Key design decisions

| Decision | Rationale |
|---|---|
| Hexagonal structure | Keeps domain logic free of I/O; enables in-memory tests |
| `afero` for filesystem | Swappable between OS and in-memory without changing application code |
| `LocationServiceFactory` pattern | Allows CLI integration tests without global state or build tags |
| No `panic` in domain code | Errors propagate up naturally; `panic` would break testability |
| `TrimSuffix("/")` on parsed paths | Prevents double-slash when the regex captures trailing slash from `.gitconfig` |
