# Codex CLI example configuration

Copy the block below into your own Codex config location (see Codex
CLI's current documentation for the exact path) and adjust
`command`/`args` if `codinho` is not on your `PATH`. This example
intentionally uses no private or machine-specific paths
(tutor-skill-host-integration Technical Plan: "Configurar MCP de
exemplo sem paths privados").

```toml
[mcp_servers.codinho]
command = "codinho"
args = ["serve"]
# Codex launches this from your project directory (the one containing
# packs/ and where you want .codinho/state/ to live), so no explicit
# working-directory override is needed for the common case.
```
