# Reddit — launch posts

Each subreddit gets a different angle. Read the rules first, engage genuinely,
and never paste the same text verbatim across subs. Do not ask for upvotes.

---

## r/selfhosted

**Title**
```
I built a self-hosted control plane to see and control what my AI agents do
```

**Body**
```
I run a few coding agents locally and got uneasy that they were calling tools
and touching the filesystem without any real visibility. I wanted something
self-hosted, with no cloud account and no data leaving the box.

I built Halo (MIT, open source). It sits as a transparent proxy between your
agents and the models/tools they use, and gives you:
- a live event stream and an interactive call graph
- allow/deny/approval policies and a sandbox for risky commands
- secret/PII detection with redaction, and a tamper-proof audit log
- a single Go binary with an embedded UI; pure-Go SQLite, no external DB
- Docker Compose and a one-line install script; telemetry off by default

It works with local models (Ollama etc.) and MCP tools. The default mode is
"observe" so it won't block anything until you turn guard mode on.

Repo: https://github.com/zfxzfxzfx852/halo
Demo GIF and screenshots in the README. Happy to answer self-hosting
questions (resource usage, Docker, reverse proxies, backups).
```

---

## r/LocalLLaMA

**Title**
```
Halo: a local-first way to observe and guard MCP tool calls and local model traffic
```

**Body**
```
A lot of us here run local models and increasingly use MCP servers for tools.
The missing piece for me was visibility: which tools got called, with what
arguments, and whether anything dangerous was happening — all without sending
that data anywhere.

Halo is a local, open-source proxy that intercepts MCP (stdio/SSE/HTTP) and
OpenAI-compatible traffic, so it works in front of Ollama and other local
endpoints. It renders a live call graph and event stream, detects prompt
injection / leaked secrets / PII fully offline, supports human approval for
suspicious calls, and keeps a hash-chained audit log.

- No cloud, no account, no telemetry by default
- Single static binary, optional Docker
- Observe mode first, so you can baseline behavior before enforcing anything

https://github.com/zfxzfxzfx852/halo

Interested in feedback from people running MCP setups — especially edge cases
in stdio bridging and SSE streaming.
```

---

## r/netsec

**Title**
```
Halo: local-first agent security — injection/secret detection, policy enforcement and hash-chained audit
```

**Body**
```
Sharing a defensive tool I've been building for securing AI agent traffic.
It's a transparent proxy over MCP and OpenAI-compatible APIs, focused on
detection and control rather than being another agent framework.

Security-relevant capabilities:
- offline detection of prompt injection, jailbreak markers, PII and secrets
  (AWS/OpenAI/GitHub tokens, private keys, etc.), with redaction
- declarative policy engine (allow/deny/ask), first-match-wins, hot reload
- process-group sandbox for risky commands; container/hypervisor drivers planned
- tamper-evident audit: SHA-256 hash chain, verification from genesis,
  JSONL export
- explicit trust boundary and out-of-scope threats documented in SECURITY.md

Everything runs locally; the detection engine has zero external dependencies
and a labeled test set for recall/false-positive checks. The audit verify
endpoint will identify the exact record that was altered.

https://github.com/zfxzfxzfx852/halo

I'd genuinely appreciate review from this community — detection blind spots,
policy model weaknesses, and hardening advice are very welcome.
```

---

## r/MachineLearning

**Title**
```
Observability and governance for AI agents without adding another framework
```

**Body**
```
Most discussion here is about models and agents, but as agents take more
autonomous actions, the operational side — observing calls and enforcing
constraints — is underserved. I built Halo as infrastructure rather than a
competing agent framework.

It transparently intercepts the traffic between an agent and its models/tools
and provides a real-time call graph, policy-based control (allow/deny/human
approval), offline detection of injection and data leakage, sandboxed
execution, and an auditable replay of sessions. It's local-first and
open source under MIT.

The architecture is a single Go binary with an embedded React console and
local SQLite storage, so it's easy to run alongside research or production
workflows. There's an "observe" mode that never blocks, which makes it
practical to adopt incrementally.

https://github.com/zfxzfxzfx852/halo

Curious how others are approaching agent monitoring and guardrails in
practice — happy to discuss the design trade-offs.
```

---

## Posting guidance

- Best time: weekday mornings, 07:00–09:00 US Eastern.
- Check each subreddit's self-promotion rules; contribute comments elsewhere
  before and after posting.
- If a mod removes a post, respect it, ask politely what rule applied, and
  adjust — never repost the same content.
