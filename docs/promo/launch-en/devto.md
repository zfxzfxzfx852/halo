# Dev.to — technical launch article

## Title options

```
Building Halo: a local-first observability and security plane for AI agents
```
```
How I built a Wireshark-style control plane for AI agent traffic
```

## Tags

`go` `ai` `opensource` `security` `architecture`

## Body (Markdown)

```markdown
AI agents have quietly become powerful. They call tools, reach the network,
and read and write files — often without us having a clear picture of what's
happening. I wanted a single place to **see** those calls and **control** the
risky ones, without sending my data to a cloud service. That became
[Halo](https://github.com/zfxzfxzfx852/halo), an open-source, local-first
observability and security plane.

This post is about how it's designed and a few of the non-obvious problems
that came up while building it.

## Why not another agent framework?

There are plenty of agent frameworks. The gap is the infrastructure around
them: the flight recorder, the safety net, the control tower. Halo is
deliberately agent-agnostic. It sits in the network path and works with any
agent that speaks MCP or an OpenAI-compatible API.

## Architecture at a glance

```
Agent ─► Halo transparent proxy ─► policy / detection / sandbox ─► model/tools
              │
              ▼
      event store + audit chain + WebSocket + live console
```

Everything ships as a single statically compiled Go binary. The web console is
embedded with `go:embed`, storage is pure-Go SQLite (no CGO), and the UI is
React + TypeScript. A decision takes roughly 40 microseconds, and core
packages are at 70%+ test coverage with a race-enabled CI.

## 1. Transparent interception, including streaming

Halo proxies both ordinary JSON and **SSE streaming** responses. The
recorder normalizes request/response pairs and strips sensitive headers
before anything is persisted.

For MCP, the stdio bridge is a bidirectional JSON-RPC relay. Tool calls are
turned into events and correlated by id, so a request and its response can be
joined even when they arrive over a stream. A small but important detail:
diagnostic logs go to **stderr**, so they never pollute the JSON-RPC stream
on stdio.

## 2. A policy engine with a custom glob matcher

The first policy engine used the standard library's `path.Match`. It mostly
worked, until it didn't: `path.Match` treats `*` as not crossing `/`, which
means a rule like `curl *` fails to match a command containing
`https://...`. I replaced it with a small custom glob where `*` matches any
character (including slashes), `?` matches a single character, and patterns
are compiled and cached. The engine evaluates rules in order, first match
wins, supports AND conditions, and hot-reloads atomically under a
read/write lock.

## 3. Offline detection and the HTML-escaping trap

Detection runs entirely offline — no model, no external calls. It covers
secret tokens, PII, and prompt injection patterns.

Here's a bug that was easy to miss. Go's `json.Marshal` escapes `<` and `>`
into `\u003c` and `\u003e`. Injection signatures containing things like
`</system>` therefore stopped matching. The fix is to normalize those
escapes back to angle brackets before scanning. A labeled test set guards
against regressions and keeps an eye on false positives.

## 4. Human approval and sandboxing

When a policy resolves to "ask," the call isn't hard-rejected — it's
**suspended**. The console receives a decision event over WebSocket and shows
an approval card; the call resumes on approval or fails closed on timeout.

For execution, the sandbox starts with a process-group driver (isolated
process group, optional privilege drop, and a timeout that kills the whole
group). Docker, gVisor, Firecracker and WASM drivers are defined as
interfaces so they can land without changing the core.

## 5. A tamper-proof audit chain

Every persisted event is appended to an audit chain:

```
hash = SHA256(prevHash + serialize(event))
```

Verification recomputes the chain from genesis and reports the exact sequence
number of the first mismatch. I tested this by tampering directly with a row
in the database — verification correctly flipped to invalid and located the
record.

## Observe first, guard later

The adoption model matters as much as the code. Halo defaults to an
**observe** mode that never blocks, so you can build a graph and a baseline
with zero risk. When you trust it, you switch to **guard** mode and policies
start enforcing.

## Try it

```bash
curl -fsSL https://raw.githubusercontent.com/zfxzfxzfx852/halo/main/scripts/install.sh | sh
halo init
halo
```

Then open the console at `http://127.0.0.1:8787/console/`.

Halo is MIT licensed and I'd love feedback — especially on detection blind
spots, the policy model, and new integration adapters.
```

## Posting guidance

- Publish on a weekday morning; cross-post to Hashnode and the engineering
  blog, but tailor the intro rather than blindly duplicating.
- Respond to comments and fold good questions into follow-up posts.
