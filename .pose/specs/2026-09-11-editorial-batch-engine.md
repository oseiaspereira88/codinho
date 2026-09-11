---
slug: editorial-batch-engine
status: in-progress
created_at: 2026-09-11
completed_at:
supersedes:
depends_on:
priority: 100
components: scripts, docs
delivers:
---

# Spec: editorial-batch-engine

## 1. Intent
### Goal
Versionar esteira editorial com fila, autores Codex isolados e revisão primária.
### Business value
Entregar lotes completos com evidência, retomada e commits rastreáveis.
### Constraints
Configurar modelo, reasoning, paralelismo e tamanho do lote na inicialização.
Padrões atuais: gpt-5.6-luna, reasoning max, um autor, cinco tarefas; até oito
autores.
### Non-goals
Substituir playtest humano ou alterar metas de go-foundations-packs.

## 2. Requirements
### Functional
- R1: Inicializar fila validada e configuração persistida imutável.
- R2: Executar lotes configuráveis em worktrees independentes, até oito autores.
- R3: Capturar sessão e logs; retomar por ID explícito, nunca --last.
- R4: Vincular revisão ao digest do diff; rejeitar aprovação obsoleta.
- R5: Integrar somente lote aprovado em árvore limpa com commit POSE-Spec.
- R6: Preservar trabalho em falha/timeout e bloquear mudanças fora do escopo.
- R7: Versionar skill, guia e inventário de todo o escopo editorial.
- R8: Executar piloto real de cinco tarefas com luna/high e revisão primária.
- R9: Despachar somente uma onda pronta, com dependências integradas, paths
  exclusivos entre lotes ativos e base capturada no HEAD de cada despacho.
- R10: Recuperar execução interrompida somente sem autor ativo; permitir
  revogar aprovação e concluir auditoria sem alterações com evidência.
### Non-functional
Testes determinísticos com backend falso; aguardar processos até conclusão.
### Security
Sem eval; subprocessos sem shell; permissões de sandbox preservadas.
### Compatibility
Preservar script legado sequencial; CLI nova em Python stdlib.

## 3. Technical Plan
### Affected areas
scripts, skills e docs/editorial. Conteúdo piloto pertence a go-foundations-packs.
### Artifacts
- created: scripts/editorial.py
- created: scripts/test_editorial.py
- modified: .gitignore
- modified: .pose/indexes/validation-matrix.json
- modified: .pose/assessments/consolidated.md
- modified: .pose/assessments/scripts.md
- modified: .pose/state/components/scripts.json
- modified: .pose/docs.json
- modified: .pose/assessments/technical-debt.md
- modified: .pose/state/technical-debt.json
- created: .pose/results/editorial-review-validation.json
- created: .agents/skills/editorial-coordinator/SKILL.md
- created: .agents/skills/editorial-author/SKILL.md
- modified: .agents/skills/README.md
- created: docs/editorial-workflow.md
- modified: docs/agent-review-workflow.md
- created: docs/editorial/go-foundations-queue.json
- created: docs/editorial/go-foundations-inventory.md
### API/contract changes
CLI init/status/run/revise/review/integrate/recover/reconcile; estado JSON por execução.
depends_on referencia IDs anteriores; tamanho de lote é limite superior.
Execuções antigas preservam sua fila; adotar dependências exige nova subfila.
### Data/storage changes
Logs, patches e sessões em diretório externo ao repositório.
### Technical risks
Conflitos entre patches do mesmo pack exigem nova integração revisada.

## 4. Tasks
### Planning
- [x] Inspecionar fluxo e consumir knowledge:go-foundations-io-batch.
- [x] Executar assess discover para scripts.
### Implementation
- [x] Implementar executor e testes.
- [x] Versionar skill e inventário.
- [x] Executar, revisar e integrar piloto.
### Validation
- [ ] Executar testes, matriz strict e revisão do mecanismo.

## 5. Decisions
Worktree por lote, sessões explícitas e integração serial. ADR:
2026-09-11-editorial-worktrees-and-explicit-sessions.md.

## 6. Validation
### Strategy
Backend falso para concorrência e falhas; piloto real para operação integrada.
### Deterministic checks
- python3 -m unittest discover -s scripts -p test_editorial.py
- pose validate --strict
- pose check --strict
- git diff --check
### Execution log
Commits `060b4f3`, `41957a2` e `6b6a799` versionam o mecanismo, a retomada sem
sessão e a recriação de worktree. Os testes determinísticos passaram (7 casos),
assim como `pose check --strict`, `pose lint-spec editorial-batch-engine
--ready-check` e `git diff --check`. O piloto foi inicializado em
`/tmp/codinho-editorial-pilot-20260911` com lote 5, paralelismo 1,
`gpt-5.6-luna`/`high`; a sessão explícita
`01a08e8e-a271-7623-b501-727b3cb8ebf8` foi capturada, mas o provedor recusou a
execução por limite de uso. O lote permanece `failed` com worktree e logs
preservados para `revise` posterior.

Revisão posterior encontrou escape por índice staged, base fixa entre lotes,
despacho sem dependências e ausência de recuperação de running. Correções
cobertas por 14 testes determinísticos. Revisão independente confirmou também
falha de commit sem reconciliação: corrigida com reconcile, que verifica pai,
patch aprovado e trailer do commit manual antes de registrar integração.
Rules applied during review: `.pose/rules/security.md` para limites de escrita,
`.pose/rules/documentation-style.md` para contratos e instruções. Workflow:
`.pose/workflows/review.md`; mudança mista. `pose assess tech-debt`: zero marcadores.
Bundle `rvb-d6a215604090b500` preparado, sem atestação: faltam change set
imutável atribuído e evidência estruturada. Matriz anterior apenas iniciada não
comprova sucesso; checkbox de validação corrigido para pendente.

Nova matriz `pose validate --strict --json .pose/results/editorial-review-validation.json`
terminou com SUCCESS; 14 testes do executor passaram. Recurrence-check: zero
chaves acima do limiar. Piloto retomado na mesma sessão e entregue para revisão:
51 checks declarados/verificados, quatro mutantes falhando nos casos novos;
diff limitado a go-core, go-errors e go-data-text, inspecionado pelo coordenador.

Piloto aprovado pelo coordenador root e integrado em `eeb5aac`, digest
`3e0419805edac00494d5047e42ad1eaa555fa95314d14a447c2c0dd19d5abba2`.
Nova matriz após integração: SUCCESS, incluindo 14 testes do executor;
evidência em `.pose/results/editorial-review-validation.json`.
Permanecem pendentes a atestação POSE final e as 65 tarefas editoriais seguintes;
revisão humana/playtest/publicação não foram realizados.

Segunda revisão independente reproduziu reconcile rejeitando alterações
independentes no mesmo arquivo: corrigido com digest capturado no índice após
aplicar ao destino; regressão incluída entre os 14 testes, todos passando após
o ajuste. A matriz estruturada acima precede esse ajuste final; a regressão
direcionada posterior o cobre. Limitações aceitas para continuidade: atualização
de base exige nova subfila; interrupção antes de registrar PID exige diagnóstico
manual. Essas limitações estão documentadas no guia e impedem alegar recuperação
automática universal. Decisão técnica: aprovado com reservas operacionais;
atestação POSE formal ainda pendente, sem fechamento da spec.

## 7. Final Report
Em implementação. Fechamento exige R1–R10 e evidências do piloto.
