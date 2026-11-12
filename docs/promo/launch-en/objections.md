# Objection handling — ready-to-paste responses

Keep these tone-neutral and factual. Don't be defensive; acknowledge the
concern, answer with evidence, and link to docs where relevant.

---

## 1. "Why should I trust a proxy that sees all my agent traffic?"

```
Fair question — a proxy is a sensitive position, so the design tries to
minimize that trust rather than ask for blind faith.

- It runs entirely on your machine and is open source (MIT), so you can read
  and build it yourself. Telemetry is off by default and there's no account.
- Secrets and sensitive headers are redacted before data is recorded, and
  storage is local SQLite.
- The exact trust boundary — including what Halo does NOT defend against —
  is documented in docs/SECURITY.md.
- You can start in observe mode and inspect exactly what it captures before
  enabling any blocking.

If anything in the threat model looks wrong, I'd genuinely consider that a
bug report.
```

## 2. "How is this different from existing agent security/logging tools?"

```
The short version: most current tools are either CLI log viewers or B2B
gateways, and they're pretty early. Halo's focus is the real-time developer
experience and the full loop in one binary.

- A live, interactive call graph and event stream — not just searchable logs.
- Observe-first adoption: get value with zero risk, then turn on guard mode.
- Policy, offline detection (injection/PII/secrets), human approval, sandbox
  and a hash-chained audit working together, rather than separate products.
- Local-first and dependency-light: one static binary, embedded UI, pure-Go
  SQLite, no cloud or external services required.

It's also not an agent framework — it's infrastructure any agent can use,
so it doesn't compete with the agent you already run.
```

## 3. "Won't a proxy add latency and become a bottleneck?"

```
It's designed to stay out of the way.

- A policy decision is around 40 microseconds in the benchmark (docs/BENCHMARKS.md),
  which is noise next to model/tool round-trips, and detection runs locally.
- Streaming (SSE) is passed through rather than buffered, so long responses
  aren't held up.
- It's a statically compiled Go binary with local SQLite; there's no network
  hop to a separate service and no external dependency in the hot path.

The main thing that adds latency is an explicit human-approval "ask" rule —
and that's intentional, because the call is supposed to wait for you.
```

---

## Bonus: "Is this production-ready?"

```
The core is tested with race detection and core packages are at 70%+
coverage, but I'd be honest about maturity: the process sandbox is solid
while the Docker/gVisor/Firecracker drivers are still interfaces, and I'd
recommend observe mode or a staging environment before enforcing in
production. Semantic versioning and a changelog will make breaking changes
clear, and the audit/verify flow helps catch integrity problems early.
```
