# Changelog

All notable changes to Halo are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project
follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] — 2026-11-13

### Added

- Transparent interception of HTTP and MCP (stdio / SSE / HTTP) calls between
  an agent and its models or tools.
- Local-first SQLite event store with a tamper-evident hash audit chain and
  JSONL export.
- Real-time React console: event stream, call graph, approval center, audit
  panel, sessions, policies, runtime status, and integration guide.
- Declarative policy engine with `allow`, `deny`, and `ask` (human approval)
  decisions, rule ordering, per-rule risk ceiling, and hot reload.
- Offline detection rules for secrets, PII (email, phone, ID card, credit
  card, IPv4), prompt injection, jailbreak, and shell command substitution.
- Human approval queue with allow / deny and timeout-deny default.
- Process-level sandbox driver with process-group isolation and timeout
  cleanup.
- REST API and WebSocket stream for events, stats, sessions, findings,
  approvals, policies, and runtime info.
- One-line install script, interactive `halo init` wizard, Homebrew tap,
  Docker Compose, and source build paths.
- Fog blue-gray, liquid-glass iOS-style landing page with an ambient particle
  background, light and dark modes, and bilingual (EN / ZH) copy.

### Known limitations

- Only the process sandbox driver is implemented; Docker, gVisor, Firecracker,
  and WASM drivers are reserved interfaces and return `ErrNotImplemented`.
- No native Windows build; run via Docker or WSL.
- Halo does not defend against threats outside its trust boundary, such as a
  compromised host or kernel-level attack. See `docs/SECURITY.md`.

[1.0.0]: https://github.com/zfxzfxzfx852/halo/releases/tag/v1.0.0
