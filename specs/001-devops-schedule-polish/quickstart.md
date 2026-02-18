# Quickstart: DevOps Infrastructure & Schedule Polish

## Prerequisites

- Go 1.24+
- GitHub CLI (`gh`) authenticated
- GoReleaser (installed via `go install github.com/goreleaser/goreleaser@latest`)

## Development Workflow

```bash
# 1. Create feature branch
git checkout -b fix/week-filter main

# 2. Make changes
# Edit internal/service/schedule/service.go

# 3. Run locally
go build -o ua ./cmd/ua-cli && ./ua schedule

# 4. Run tests
go test ./...

# 5. Commit + Push + PR
git add -A && git commit -m "fix: filter events to current week"
git push -u origin fix/week-filter
gh pr create --fill

# 6. After PR merged, tag a release
git checkout main && git pull
git tag v0.2.0
git push origin v0.2.0
# → GoReleaser builds binaries automatically
```

## CI/CD Quick Reference

| Event | Workflow | What Happens |
|-------|----------|-------------|
| Push any branch | `ci.yml` | lint + build + test |
| PR to `main` | `ci.yml` | lint + build + test (blocks merge on fail) |
| Push `v*` tag | `release.yml` | GoReleaser builds + publishes binaries |
