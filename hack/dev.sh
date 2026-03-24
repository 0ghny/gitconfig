#!/usr/bin/env bash
# hack/dev.sh — development task runner for gitconfig
# Usage: ./hack/dev.sh [--coverage|-c] build|test|lint|fmt|deps|docs|docker-build|test-actions-ci|clean|all
# When called without arguments and fzf is available, an interactive menu is shown.

set -euo pipefail

ROOT="$(git -C "$(dirname "${BASH_SOURCE[0]}")" rev-parse --show-toplevel)"
TARGET_DIR="${ROOT}/target"
BINARY_NAME="gitconfig"
WITH_COVERAGE=0

# Ensure installed Go tools (godoc, etc.) are on PATH.
GOBIN="$(go env GOBIN)"
if [ -z "${GOBIN}" ]; then
  GOBIN="$(go env GOPATH)/bin"
fi
export GOBIN
export PATH="${GOBIN}:${PATH}"

# ---------------------------------------------------------------------------
# build — compile for the current OS and architecture
# ---------------------------------------------------------------------------
__run_build() {
  local os arch ts report_dir out
  os="$(go env GOOS)"
  arch="$(go env GOARCH)"
  ts="$(date +%Y%m%d-%H%M%S)"
  report_dir="${TARGET_DIR}/${ts}"
  out="${report_dir}/${BINARY_NAME}"

  mkdir -p "${report_dir}"

  echo "┌─────────────────────────────────────────────────────────────"
  echo "│  Build"
  echo "│"
  printf "│  %-20s %s\n" "OS:"      "${os}"
  printf "│  %-20s %s\n" "Arch:"    "${arch}"
  printf "│  %-20s %s\n" "Output:"  "${out}"
  echo "└─────────────────────────────────────────────────────────────"
  echo ""

  CGO_ENABLED=0 GOOS="${os}" GOARCH="${arch}" \
    go build -o "${out}" ./cmd

  echo "==> Build OK: ${out}"
}

# ---------------------------------------------------------------------------
# test — run the full test suite; pass --coverage/-c to generate a report
# ---------------------------------------------------------------------------
__run_test() {
  local ts report_dir profile html
  ts="$(date +%Y%m%d-%H%M%S)"
  report_dir="${TARGET_DIR}/${ts}"
  profile="${report_dir}/coverage.out"
  html="${report_dir}/coverage.html"

  mkdir -p "${report_dir}"

  echo "┌─────────────────────────────────────────────────────────────"
  echo "│  Test"
  echo "│"
  printf "│  %-20s %s\n" "Scope:"    "./..."
  printf "│  %-20s %s\n" "Coverage:" "$([[ "${WITH_COVERAGE}" -eq 1 ]] && echo "enabled" || echo "disabled")"
  [[ "${WITH_COVERAGE}" -eq 1 ]] && printf "│  %-20s %s\n" "Report dir:" "${report_dir}"
  echo "└─────────────────────────────────────────────────────────────"
  echo ""

  if [[ "${WITH_COVERAGE}" -eq 1 ]]; then
    go test -coverprofile="${profile}" ./...
    go tool cover -html="${profile}" -o "${html}"
    go tool cover -func="${profile}" | tee "${report_dir}/coverage.txt"
    echo ""
    echo "==> Test OK — coverage report: ${html}"
  else
    go test ./...
    echo ""
    echo "==> Test OK"
  fi
}

# ---------------------------------------------------------------------------
# lint — run golangci-lint
# ---------------------------------------------------------------------------
__run_lint() {
  if ! command -v golangci-lint &>/dev/null; then
    echo "ERROR: golangci-lint is not installed." >&2
    echo "" >&2
    echo "  Install it with one of:" >&2
    echo "    macOS:  brew install golangci-lint" >&2
    echo "    Go:     go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" >&2
    echo "    Manual: https://golangci-lint.run/usage/install/" >&2
    return 1
  fi

  echo "┌─────────────────────────────────────────────────────────────"
  echo "│  Lint"
  echo "│"
  printf "│  %-20s %s\n" "Tool:" "golangci-lint $(golangci-lint --version 2>/dev/null | head -1)"
  echo "└─────────────────────────────────────────────────────────────"
  echo ""

  golangci-lint run ./...

  echo ""
  echo "==> Lint OK"
}

# ---------------------------------------------------------------------------
# fmt — check formatting with gofmt
# ---------------------------------------------------------------------------
__run_fmt() {
  echo "┌─────────────────────────────────────────────────────────────"
  echo "│  Format check (gofmt)"
  echo "└─────────────────────────────────────────────────────────────"
  echo ""

  local unformatted
  unformatted="$(gofmt -l .)"

  if [ -n "${unformatted}" ]; then
    echo "==> Fmt FAILED: the following files need formatting:"
    echo "${unformatted}" | sed 's/^/    /'
    echo ""
    echo "  Run: gofmt -w <file>  to auto-fix."
    return 1
  fi

  echo "==> Fmt OK: all files are correctly formatted."
}

# ---------------------------------------------------------------------------
# deps — tidy + check for available upgrades (use --upgrade to also apply them)
# ---------------------------------------------------------------------------
__run_deps() {
  echo "==> Running go mod tidy..."
  go mod tidy
  echo ""
  echo "==> Available dependency upgrades:"
  go list -u -m all
}

# ---------------------------------------------------------------------------
# clean — remove target/ and purge Go build/test/module caches
# ---------------------------------------------------------------------------
__run_clean() {
  echo "┌─────────────────────────────────────────────────────────────"
  echo "│  Clean"
  echo "│"
  printf "│  %-20s %s\n" "Target dir:" "${TARGET_DIR}"
  printf "│  %-20s %s\n" "Go caches:"  "build, test, module"
  echo "└─────────────────────────────────────────────────────────────"
  echo ""

  if [ -d "${TARGET_DIR}" ]; then
    rm -rf "${TARGET_DIR}"
    echo "==> Removed ${TARGET_DIR}"
  else
    echo "==> Nothing to remove in ${TARGET_DIR}"
  fi

  go clean -cache
  echo "==> Go build cache cleared"

  go clean -testcache
  echo "==> Go test cache cleared"

  go clean -modcache
  echo "==> Go module cache cleared"

  echo ""
  echo "==> Clean OK"
}

# ---------------------------------------------------------------------------
# docs — serve godoc documentation locally on http://localhost:6060
# ---------------------------------------------------------------------------
__run_docs() {
  if ! command -v godoc &>/dev/null; then
    echo "==> godoc not found — installing..."
    go install golang.org/x/tools/cmd/godoc@latest
  fi

  local pkg="github.com/0ghny/gitconfig"
  local url="http://localhost:6060/pkg/${pkg}/?m=all"

  echo "┌─────────────────────────────────────────────────────────────"
  echo "│  Godoc server"
  echo "│"
  printf "│  %-20s %s\n" "Package:" "${pkg}"
  printf "│  %-20s %s\n" "URL:" "${url}"
  echo "└─────────────────────────────────────────────────────────────"
  echo ""
  echo "==> Opening ${url}  (Ctrl-C to stop)"
  echo ""

  godoc -http=:6060
}

# ---------------------------------------------------------------------------
# docker-build — build a Docker image for the current git version
# ---------------------------------------------------------------------------
__run_docker_build() {
  if ! command -v docker &>/dev/null; then
    echo "ERROR: docker is not installed." >&2
    return 1
  fi

  local version
  version="$(git -C "${ROOT}" describe --tags --always --dirty 2>/dev/null || echo dev)"
  local image="gitconfig:${version}"

  echo "┌─────────────────────────────────────────────────────────────"
  echo "│  Docker build"
  echo "│"
  printf "│  %-20s %s\n" "Image:"   "${image}"
  printf "│  %-20s %s\n" "Version:" "${version}"
  printf "│  %-20s %s\n" "Context:" "${ROOT}"
  echo "└─────────────────────────────────────────────────────────────"
  echo ""

  docker build \
    --build-arg VERSION="${version}" \
    --tag "${image}" \
    "${ROOT}"

  echo ""
  echo "==> Docker build OK: ${image}"
  echo ""
  echo "Usage example:"
  echo "  docker run --rm \\"
  echo "    -v \"\$HOME/.gitconfig:/root/.gitconfig\" \\"
  echo "    -v \"\$HOME/.gitconfigs:/root/.gitconfigs\" \\"
  echo "    ${image} locations"
}

# ---------------------------------------------------------------------------
# test-actions-ci — run the GitHub Actions CI workflow locally using act
# ---------------------------------------------------------------------------
__run_ci() {
  if ! command -v act &>/dev/null; then
    echo "ERROR: act is not installed." >&2
    echo "" >&2
    echo "  Install it with one of:" >&2
    echo "    macOS:  brew install act" >&2
    echo "    Go:     go install github.com/nektos/act@latest" >&2
    echo "    Manual: https://nektosact.com/installation/" >&2
    return 1
  fi

  local workflow="${ROOT}/.github/workflows/ci.yml"

  echo "┌─────────────────────────────────────────────────────────────"
  echo "│  CI (local — powered by act)"
  echo "│"
  printf "│  %-20s %s\n" "Workflow:" "${workflow}"
  printf "│  %-20s %s\n" "Event:"    "pull_request"
  echo "└─────────────────────────────────────────────────────────────"
  echo ""

  local act_opts=()
  if [[ "$(uname -m)" == "arm64" ]]; then
    act_opts+=(--container-architecture linux/amd64)
    echo "  NOTE: Apple M-series detected — adding --container-architecture linux/amd64"
    echo ""
  fi

  act pull_request \
    --workflows "${workflow}" \
    --directory "${ROOT}" \
    "${act_opts[@]}"
}

# ---------------------------------------------------------------------------
# all — fmt + lint + test
# ---------------------------------------------------------------------------
__run_all() {
  local exit_code=0
  __run_fmt   || exit_code=$?
  echo ""
  __run_lint  || exit_code=$?
  echo ""
  __run_test  || exit_code=$?
  return "${exit_code}"
}

# ---------------------------------------------------------------------------
# Dispatch
# ---------------------------------------------------------------------------
__run_task() {
  local task="${1}"
  case "${task}" in
    build)        __run_build         ;;
    test)         __run_test          ;;
    lint)         __run_lint          ;;
    fmt)          __run_fmt           ;;
    deps)         __run_deps          ;;
    docs)         __run_docs          ;;
    docker-build) __run_docker_build  ;;
    test-actions-ci) __run_ci         ;;
    clean)        __run_clean         ;;
    all)          __run_all           ;;
    *)
      echo "Unknown task: ${task}" >&2
      echo "Usage: $(basename "$0") [--coverage|-c] build|test|lint|fmt|deps|docs|docker-build|test-actions-ci|clean|all" >&2
      exit 1
      ;;
  esac
}

# ---------------------------------------------------------------------------
# Parse flags
# ---------------------------------------------------------------------------
_args=()
for _arg in "$@"; do
  case "${_arg}" in
    --coverage|-c) WITH_COVERAGE=1 ;;
    *) _args+=("${_arg}") ;;
  esac
done
set -- "${_args[@]+"${_args[@]}"}"

if [ -z "${1:-}" ]; then
  if ! command -v fzf &>/dev/null; then
    echo "fzf not found. Usage: $(basename "$0") [--coverage|-c] build|test|lint|fmt|deps|docs|docker-build|test-actions-ci|clean|all" >&2
    exit 1
  fi
  selected=$(printf 'build\ntest\nlint\nfmt\ndeps\ndocs\ndocker-build\ntest-actions-ci\nclean\nall\n' \
    | fzf --prompt="Select task > " --height=10 --border --ansi)
  [ -z "${selected}" ] && exit 0
  __run_task "${selected}"
else
  __run_task "${1}"
fi
