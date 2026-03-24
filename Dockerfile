# syntax=docker/dockerfile:1

# ── Build stage ────────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /src

# Cache dependencies separately from source
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/gitconfig ./cmd

# ── Final stage ────────────────────────────────────────────────────────────
FROM scratch

COPY --from=builder /out/gitconfig /gitconfig
# Mount your ~/.gitconfig and ~/.gitconfigs at runtime, e.g.:
#   docker run --rm \
#     -v "$HOME/.gitconfig:/root/.gitconfig" \
#     -v "$HOME/.gitconfigs:/root/.gitconfigs" \
#     gitconfig locations
ENTRYPOINT ["/gitconfig"]
