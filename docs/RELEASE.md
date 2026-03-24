# Release Process

`gitconfig` uses [goreleaser](https://goreleaser.com/) together with GitHub Actions to produce and publish release artefacts.

---

## Release artefacts

Each release produces:

- Standalone binaries for **Linux**, **macOS**, and **Windows** (amd64 + arm64), built with `CGO_ENABLED=0`.
- A **GitHub Release** with auto-generated changelog and attached binaries.
- An updated **asdf plugin** version (via the `asdf-gitconfig` plugin).

---

## Versioning

The project follows [Semantic Versioning](https://semver.org/):

```
v<MAJOR>.<MINOR>.<PATCH>
```

| Change | Increment |
|---|---|
| Breaking API / behaviour change | MAJOR |
| New, backward-compatible feature | MINOR |
| Bug fix, docs, chore | PATCH |

---

## Creating a release

### 1 — Ensure `main` is clean and passing

```sh
git checkout main
git pull
./hack/dev.sh all   # fmt + lint + test
```

### 2 — Tag the release

```sh
# Replace x.y.z with the new version
git tag -a vx.y.z -m "Release vx.y.z"
git push origin vx.y.z
```

### 3 — GitHub Actions triggers automatically

Pushing a tag that matches `v*` starts the release workflow, which:

1. Runs `goreleaser release --clean` using the config in `.goreleaser.yaml`.
2. Compiles binaries for all target platforms.
3. Creates a GitHub Release and attaches the binaries.

### 4 — Verify the release

- Check the [Releases page](https://github.com/0ghny/gitconfig/releases) on GitHub.
- Download and smoke-test the binary for at least one platform:

```sh
gitconfig --version
gitconfig locations
```

---

## Local dry-run (without publishing)

To test the goreleaser config locally without creating a release:

```sh
goreleaser release --snapshot --clean
# Artefacts are written to dist/
ls dist/
```

---

## `.goreleaser.yaml` reference

Located at the repository root. Key settings:

| Setting | Value |
|---|---|
| Binary name | `gitconfig` |
| Entry point | `./cmd` |
| Target OS | linux, darwin, windows |
| Target arch | amd64, arm64 |
| CGO | disabled (`CGO_ENABLED=0`) |

To add a new platform, edit the `goos` / `goarch` lists in `.goreleaser.yaml` and open a PR.

---

## Hotfix releases

For urgent fixes that cannot wait for the next planned release:

1. Branch from the tag: `git checkout -b fix/<description> vx.y.z`
2. Apply the fix, add a test.
3. Bump the patch version, tag, and push.

---

## Deprecation policy

- Deprecated features are announced in the release notes one minor version before removal.
- Breaking changes bump the major version and are documented in the release notes.
