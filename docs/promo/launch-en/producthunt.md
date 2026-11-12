# Product Hunt — launch kit

## One-line tagline (max 60 chars)

**Primary**
```
See and secure every call your AI agent makes
```
(40 chars)

**Alternatives**
```
The local-first control plane for AI agents
```
(42 chars)
```
Wireshark-style observability and security for agents
```
(54 chars)

## Short description

```
Halo transparently intercepts the traffic between your AI agents and the
models or tools they use, giving you real-time visibility, policy control,
sandboxed execution and tamper-proof audit — all local-first in one binary.
```

## Long description

```
AI agents now call tools, access the network and read or write files on their
own. Halo gives you a clear window into every one of those actions and a way
to control them, without replacing your agent or locking you into a cloud.

WHAT YOU GET
- Real-time event stream and an interactive call graph with risk color-coded
- Declarative allow / deny / human-approval policies that hot-reload
- Offline detection of prompt injection, jailbreak attempts, PII and secrets,
  with automatic redaction
- Sandboxed execution for high-risk commands
- A hash-chained audit log that verifies integrity and replays sessions
- Zero-instrumentation support for MCP (stdio/SSE/HTTP) and OpenAI-compatible APIs

START SAFE
Halo defaults to "observe" mode: watch and learn normal behavior with zero
chance of blocking a workflow. Flip to "guard" mode when you're ready.

LOCAL-FIRST BY DESIGN
A single statically compiled Go binary with an embedded console and local
SQLite. Your data stays on your machine, telemetry is off by default, and it
works with local models such as Ollama. Install via one-line script, Homebrew,
Docker or source.

Open source under MIT.
```

## Maker comment (first comment on launch day)

```
Hi everyone 👋 I'm the maker of Halo.

The spark was simple: my coding agents were taking real actions on my machine,
and I couldn't easily see what they were doing or stop a bad call in time.
Rather than build yet another agent, I wanted the safety layer any agent could
plug into.

Two things matter most to me:
1. It should be safe to adopt — so observe mode never blocks anything, and you
   turn enforcement on only after you trust what you see.
2. It should respect you — local-first, no cloud lock-in, no telemetry,
   fully open source.

The demo GIF in the gallery shows the full loop: watching calls, inspecting
the details, and approving or rejecting a risky one. I'll be here all day
answering questions and taking notes. Tell me what integrations or sandbox
backends you'd want next!
```

## Assets checklist

- Product icon: 240×240 (reuse gradient halo mark)
- Gallery: demo GIF (720px) + 3–4 screenshots at 1270×760
- Cover/og image consistent with landing page colors
- Website link, GitHub link, pricing: Free / Open Source

## Launch guidance

- Schedule for 00:01 Pacific Time on a Tuesday/Wednesday; line up a hunter or
  launch as maker in advance.
- Have a support team/chat ready; reply to every comment on day one.
- Encourage genuine users to comment; do not trade or buy upvotes.
