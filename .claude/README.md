# `.claude/` — agent configuration

Configures how Claude Code (and other agents) work in this repo.

- **`settings.json`** — pre-approves safe commands (build/test/lint/codegen via
  `make`, read-only git & search, health-check curls) so the agent isn't prompted
  for them. Stack-changing actions (`make up/down/nuke`, `docker compose`,
  `git commit`, `git push`) still ask first; `.env*` files are denied.
- **`design-prompt.md`** — a self-contained brief to paste into a fresh Claude
  conversation / design tool to design the **web frontend** (which doesn't exist
  yet). It captures the product, users, and the gateway/BFF API the UI talks to.

The **primary agent context** is [`../CLAUDE.md`](../CLAUDE.md) at the repo root —
Claude Code auto-loads that file (it does **not** auto-load files inside
`.claude/`). Edit `CLAUDE.md` to change what every session knows by default.
