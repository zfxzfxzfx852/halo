# Hacker News — Show HN

## Title (pick one)

**Primary**
```
Show HN: Halo – Local-first observability and security for AI agents
```

**Alternative**
```
Show HN: Halo – A Wireshark-style control plane for every AI agent call
```

## Submission body

```
Hi HN,

As agents started calling tools, hitting the network and reading/writing files
on my machine, I realized I had almost no visibility into what they were
actually doing — and no way to stop a bad call before it happened. Existing
agent-security tools were mostly CLI log viewers or B2B gateways, and none
made the real-time picture genuinely easy to understand.

So I built Halo: a local-first, self-hosted observability and security plane
that sits transparently between your agents and the models/tools they use.
It does not replace your agent and it is not another agent framework — it is
the infrastructure every agent can rely on.

What it does:
- Transparent interception of MCP (stdio/SSE/HTTP) and OpenAI-compatible APIs
  with zero changes to agent code.
- A live event stream and an interactive call graph — risk is color-coded and
  visible at a glance (think Wireshark/Datadog for agents).
- Declarative allow / deny / human-approval policies with hot reload.
- Fully offline detection of prompt injection, jailbreak, PII and secrets,
  with automatic redaction.
- Sandboxed execution for high-risk commands (process groups today;
  Docker/gVisor/Firecracker drivers planned).
- A tamper-proof, hash-chained audit log that can replay whole sessions.

The design is "observe first, guard later": you start in observe mode with no
risk of breaking anything, build a graph and a baseline, then flip on guard
mode when you trust it.

It ships as a single statically compiled binary (Go + an embedded React
console, CGO-free, pure-Go SQLite), stores data on your machine, and has
telemetry off by default. A decision is around 40 microseconds, and core
packages are at 70%+ test coverage with a race-enabled CI.

One-line install on macOS/Linux:
curl -fsSL https://raw.githubusercontent.com/zfxzfxzfx852/halo/main/scripts/install.sh | sh

Repo and demo: https://github.com/zfxzfxzfx852/halo
Landing page: https://zfxzfxzfx852.github.io/halo/

Happy to answer questions about the architecture, the threat model, the
detection rules, or the MCP plumbing. Feedback very welcome.
```

## First comment (post immediately as the author)

```
A few implementation notes for the curious:

- The MCP stdio bridge is a bidirectional JSON-RPC relay; tool calls are
  turned into events and correlated by id, so request/response can be joined
  even across streaming. Logs go to stderr so they never pollute JSON-RPC.
- The HTTP proxy handles both plain JSON and SSE streaming and strips
  sensitive headers before recording.
- Policies use a custom glob matcher (the stdlib path.Match treats '*' as not
  crossing '/', which broke command strings like 'curl https://...'), with a
  compiled-regex cache and first-match-wins semantics.
- The audit chain is SHA-256 over the previous hash plus the serialized
  event; verification recomputes from genesis and will pinpoint the exact
  sequence number if a record is tampered with.
- Everything detection-related runs offline — no external calls, no model
  dependency. We normalize HTML-escaped angle brackets before scanning
  because Go's json.Marshal escapes '<' and '>', which otherwise hid
  '</system>' injection patterns.

Trust boundary and what Halo explicitly does NOT defend against are written
up in docs/SECURITY.md. This is MIT licensed and I'd love contributors,
especially more integration adapters and sandbox drivers.
```

## Posting guidance

- Best window: Tue–Thu, 08:00–10:00 US Eastern.
- Use a real, aged account; reply to every substantive comment in the first
  2–3 hours; never argue, stay technical and humble.
- Do not ask for upvotes anywhere. Do not post the same link from other
  accounts or coordinate voting.
