---
spec: administrative-cli-fixtures
category: added
breaking: false
refs:
---

A CLI `codinho` ganha comandos administrativos completos: `init`, `catalog
validate/list/show`, `session inspect`, `progress show/export` e `workspace
prepare`. `workspace prepare` materializa o código inicial de um desafio em
um destino explícito, sempre com preview (`--dry-run`) e sem sobrescrever
nada por padrão.
