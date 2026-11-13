# Contributing to Halo

Thanks for taking the time to improve Halo. This guide covers the local
workflow, conventions, and what a good pull request looks like.

## Prerequisites

- Go 1.23+
- Node.js 18+ and pnpm (for the web console)
- Git

## Development setup

```bash
# Clone the repository
git clone https://github.com/zfxzfxzfx852/halo.git
cd halo

# Install frontend dependencies
make web-install

# Build the embedded console and the full binary
make build

# Run locally
./bin/halo
```

The console is served at `http://127.0.0.1:8787/console/`; health check:
`curl http://127.0.0.1:8787/healthz`.

## Common commands

| Command | What it does |
| --- | --- |
| `make vet` | Run `go vet ./...` |
| `make test` | Run all tests with `-race` and coverage |
| `make build` | Build the frontend, then compile the embedded binary |
| `make web-dev` | Start the frontend dev server |
| `make fmt` | Format Go code |

Before opening a pull request, make sure `make vet` and `make test` pass, and
the frontend builds cleanly (`make web-build`).

## Project layout

- `cmd/halo` — entrypoint and wiring
- `internal/proxy` — transparent HTTP / MCP interception
- `internal/store` — SQLite persistence and audit chain
- `internal/detect` — offline secret / PII / injection rules
- `internal/policy` — declarative allow / deny / ask engine
- `internal/approval` — human approval queue
- `internal/server` — REST API and WebSocket hub
- `web` — embedded React console
- `website` — static landing page

## Commit conventions

Use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` new capability
- `fix:` bug fix
- `docs:` documentation only
- `refactor:` internal change with no behavior change
- `test:` adding or fixing tests
- `chore:` build, CI, or tooling

Keep commits focused; one logical change per commit. Reference the issue
number in the commit body or pull request when there is one.

## Pull request checklist

- [ ] Tests added or updated for the change; `make test` passes
- [ ] `make vet` is clean and the frontend builds
- [ ] Public behavior is documented in `README.md` or `docs/` when relevant
- [ ] No secrets, credentials, or personal data are committed
- [ ] The PR describes what changed, why, and how it was verified

## Reporting issues

- Bug reports and feature requests use the issue templates in this repository.
- For security vulnerabilities, **do not open a public issue**. Follow the
  private disclosure process in [`SECURITY.md`](./docs/SECURITY.md).

## Code of conduct

This project follows the [Contributor Covenant](./CODE_OF_CONDUCT.md). By
participating, you are expected to uphold it.
