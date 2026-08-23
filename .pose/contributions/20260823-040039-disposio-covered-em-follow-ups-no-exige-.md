---
title: "Disposição 'covered' em follow-ups não exige âncora verificável na spec alvo"
kind: suggestion
engine_version: 1.7.9
reported_at: 2026-08-23T01:00:39-03:00
---

# POSE Engine Report: Disposição 'covered' em follow-ups não exige âncora verificável na spec alvo

## Description
Ao fechar uma spec (pose-spec-closeout), um follow-up pode receber disposition [covered: <slug>] apontando para outra spec já existente. Hoje o engine (lint-spec, followups, project-state) aceita esse valor sem checar se a spec alvo referencia de fato o conteúdo coberto (requirement, depends_on ou texto equivalente) — basta o slug existir e não apontar para si mesma. Isso cria uma garantia falsa: [open] é revisitado ativamente (aparece em 'pose followups --open'), mas [covered] desaparece do backlog vivo e só seria lembrado se alguém releu manualmente os follow-ups da spec de origem ao implementar a spec de destino. Em um caso real (codinho, spec feedback-evaluation-progression fechada em 2026-08-23), dois follow-ups covered apontavam para specs draft (safe-check-executor, v1-integrated-acceptance) que não mencionavam o assunto em nenhum requirement/depends_on/texto — só corrigido manualmente após revisão humana. Sugestão: lint-spec/followups poderia emitir um WARNING (não bloqueante) quando um [covered: <slug>] não encontra nenhuma referência textual mínima (ex.: o slug de origem, ou um termo-chave do texto do follow-up) na spec alvo, incentivando o autor a adicionar um requirement/depends_on explícito em vez de confiar só no texto do follow-up já fechado.

---
### System Context (Auto-generated)
- **POSE Engine Version:** 1.7.9
- **OS/Arch:** linux/amd64
- **Go Version:** go1.26.5-X:nodwarf5
- **Reported At:** 2026-08-23T01:00:39-03:00
- **Kind:** suggestion

