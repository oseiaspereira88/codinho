---
slug: codinho-v1
status: active
created_at: 2026-08-22
depends_on:
---

# Roadmap: codinho-v1

## Contexto

Este roadmap transforma a visão de [`PROJECT.md`](../../PROJECT.md) em entregas
incrementais, governadas e verificáveis. Ele cobre toda a V1 sem considerar o
produto entregue antes de evidências determinísticas, integração real com o
host e cumprimento dos gates editoriais do catálogo.

## Estratégia de entrega

- Estabeleça decisões duráveis antes dos contratos públicos.
- Entregue primeiro uma fatia executável mínima por MCP `stdio`.
- Separe pedagogia, evidência, conteúdo e hardening em specs independentes.
- Expanda conteúdo somente depois de validar autoria e checks seguros.
- Encerre a V1 por composição e critérios de aceite, não por contagem de arquivos.

## Revisão do planejamento — 2026-09-07 UTC

Consulte o [diagnóstico e a ordem de execução](../reports/2026-09-07-doc-audit-auditoria-do-planejamento-v1.md).
As specs históricas done conservam seu fechamento; as lacunas observadas no
produto atual têm remediações próprias. Autoria curricular pode continuar
em paralelo, mas publicação exige integridade editorial e playtest humano.
Priorize recuperação, navegação completa e gates antes de ampliar o uso real.

## Milestone: architecture-baseline
- after:
- specs: architecture-decision-baseline

## Milestone: executable-foundation
- after: architecture-baseline
- specs: go-runtime-foundation, learning-domain-model, catalog-schema-loader, local-event-store, eventstore-idempotency-scope, mcp-stdio-foundation

## Milestone: pedagogical-core
- after: executable-foundation
- specs: session-orchestration-disclosure, assistance-hints-detours, feedback-evaluation-progression, workspace-observation-baselines, safe-check-executor, mastery-review-scheduling, curriculum-graph-path-recommendation

## Milestone: tutor-surfaces-and-modes
- after: pedagogical-core
- specs: tutor-skill-host-integration, administrative-cli-fixtures, learning-practice-debug-modes, interview-mode

## Milestone: adaptive-content-and-selection
- after: session-continuity, editorial-integrity
- specs: learning-track-composition, agent-authored-catalog-drafts

## Milestone: session-continuity
- after: tutor-surfaces-and-modes
- specs: session-recovery-version-pinning, session-tree-progression

## Milestone: editorial-integrity
- after: tutor-surfaces-and-modes
- specs: catalog-publication-integrity, concept-content-authoring

## Milestone: curriculum-v1
- after: tutor-surfaces-and-modes
- specs: catalog-authoring-quality, go-foundations-packs, go-backend-packs, go-production-architecture-packs, go-interviews-pack

## Milestone: production-hardening
- after: tutor-surfaces-and-modes
- specs: security-privacy-hardening, reliability-observability-compatibility, installation-documentation-ci

## Milestone: delivery-assurance
- after: production-hardening
- specs: v1-delivery-ci-assurance

## Milestone: v1-acceptance
- after: curriculum-v1, adaptive-content-and-selection, delivery-assurance
- specs: v1-integrated-acceptance

## Critérios de conclusão

- Todas as specs estão `done`, com evidência atribuída e revisão válida.
- Os 84 desafios, 160 conceitos, 100 competências e 12 trilhas atendem aos gates editoriais.
- CLI, skill e MCP funcionam pelo Codex CLI e pela extensão da IDE.
- Segurança, recuperação, privacidade, instalação e documentação cumprem `PROJECT.md`.
- O aceite integrado comprova a progressão de microguiado até autonomia.

## Cut criteria

- C1: CLI alcançável e integrada: surface:codinho-cli check:cli-reachability check:cli-e2e
- C2: Sessões recuperáveis e árvore completa: capability:session-recovery capability:session-tree-progression
- C3: Catálogo curado e distribuição verificada: capability:catalog-publication-integrity manual-review:docs/acceptance/v1-requirement-matrix.md
- C4: Conteúdo canônico e tutor composto: capability:canonical-concept-content capability:codinho-tutor check:tutor-skill-routing
- C5: Integração de trilhas e rascunhos nos dois hosts: manual-review:docs/acceptance/v1-pilot-report.md
- C6: Segurança e plataformas verificadas no candidato: governance:v1-delivery-ci manual-review:docs/acceptance/v1-release-readiness.md
- C7: Protocolo, inicialização e instalação: check:mcp-contract check:stdout-purity check:startup-budget check:smoke-install
- C8: Aceite integral, revisão independente e limitações: manual-review:docs/acceptance/v1-requirement-matrix.md manual-review:docs/acceptance/v1-pilot-report.md manual-review:docs/acceptance/v1-release-readiness.md

Referências manuais exigem conteúdo revisado, commit e resultados; a existência
do arquivo não comprova comportamento. O roll-up terminal acontece depois do
fechamento da spec de aceite e dos milestones, sem dispensar seus gates prévios.

## Riscos do roadmap

- O volume editorial pode dominar o desenvolvimento; aplique os gates antes de ampliar packs.
- Contratos MCP prematuros podem causar retrabalho; estabilize o domínio e o envelope primeiro.
- A granularidade pode virar receita; valide autonomia e transferência em sessões reais.
- Specs paralelas podem disputar os mesmos arquivos; respeite a DAG e atribua commits por trailer.
