---
type: note
slug: agent-batch-review-pattern
owner: @oseiaspereira
sensitivity: public-internal
created_at: 2026-08-23
last_reviewed_at: 2026-08-23
expires_at: 2026-09-23
source_refs:
  spec: "go-foundations-packs"
  workflow: ".pose/workflows/feature.md"
  commands: []
  external_sources: []
---

# note: agent-batch-review-pattern

## Context

`go-foundations-packs` Decision 3 adotou uma segunda instância de agente
(Codex CLI, não-interativa e sandboxed) como pré-revisor padrão de cada
lote de conteúdo autorado, antes da revisão humana. Validado ao vivo nesta
sessão em dois lotes reais, encontrando problemas objetivos de verdade
(vazamento de solução, critério não verificável, check sem evidência real
por trás). O padrão é agnóstico ao par de agentes e à tarefa — não é
específico de conteúdo curricular — então foi extraído para artefatos
reutilizáveis em vez de ficar preso a uma spec só.

## Current state

Feito: `scripts/agent-review.sh` (primitiva de invocação new/resume,
agnóstica de agente, hoje só com backend `codex`),
`docs/agent-review-workflow.md` (papéis autor/revisor, templates de
prompt, por que usar `resume` em vez de sessão nova a cada rodada, nota
sobre por que A2A não se aplica aqui), e a skill
`.claude/skills/agent-batch-review/SKILL.md` para invocar o padrão em
sessões futuras sem redescobrir os flags corretos.

## Next checks

- Validar `codex exec resume --last` de verdade (rodada 3 do checkpoint 2
  de `go-foundations-packs` foi a primeira tentativa real de resume desta
  sessão; confirmar que o ganho de custo/tempo por rodada é real, não só
  teórico).
- Quando outro backend (segunda instância de Claude, uma CLI de
  Antigravity/`agy` etc.) for usado pela primeira vez, adicionar o `case`
  correspondente em `scripts/agent-review.sh` e registrar aqui o que
  funcionou.

## Risks

- `codex exec resume --last` depende do cwd e de ser a sessão mais
  recente — rodadas paralelas/concorrentes na mesma máquina podem resolver
  para a sessão errada; hoje o padrão assume rodadas sequenciais.
- O padrão só reduz o que sobra para revisão humana; nunca preenche
  `reviewed_by`/`playtested` nem qualquer gate humano equivalente de outro
  projeto que reusar este padrão — repetir esse aviso em todo prompt de
  revisor é responsabilidade de quem usa a skill, não é imposto por
  código.

## Next owner

Mesmo owner.

## References

- Doc: `docs/agent-review-workflow.md`
- Script: `scripts/agent-review.sh`
- Skill: `.claude/skills/agent-batch-review/SKILL.md`
- Spec: `.pose/specs/2026-08-22-go-foundations-packs.md` (Decision 3, Execution log)
