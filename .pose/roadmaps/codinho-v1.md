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
- after: tutor-surfaces-and-modes
- specs: learning-track-composition, agent-authored-catalog-drafts

## Milestone: curriculum-v1
- after: adaptive-content-and-selection
- specs: catalog-authoring-quality, go-foundations-packs, go-backend-packs, go-production-architecture-packs, go-interviews-pack

## Milestone: production-hardening
- after: curriculum-v1
- specs: security-privacy-hardening, reliability-observability-compatibility, installation-documentation-ci

## Milestone: v1-acceptance
- after: production-hardening
- specs: v1-integrated-acceptance

## Critérios de conclusão

- Todas as specs estão `done`, com evidência atribuída e revisão válida.
- Os 84 desafios, 160 conceitos, 100 competências e 12 trilhas atendem aos gates editoriais.
- CLI, skill e MCP funcionam pelo Codex CLI e pela extensão da IDE.
- Segurança, recuperação, privacidade, instalação e documentação cumprem `PROJECT.md`.
- O aceite integrado comprova a progressão de microguiado até autonomia.

## Riscos do roadmap

- O volume editorial pode dominar o desenvolvimento; aplique os gates antes de ampliar packs.
- Contratos MCP prematuros podem causar retrabalho; estabilize o domínio e o envelope primeiro.
- A granularidade pode virar receita; valide autonomia e transferência em sessões reais.
- Specs paralelas podem disputar os mesmos arquivos; respeite a DAG e atribua commits por trailer.
