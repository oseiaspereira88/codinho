---
slug: agent-authored-catalog-drafts
status: in-progress
created_at: 2026-08-24
completed_at:
supersedes:
depends_on: catalog-authoring-quality, learning-track-composition, tutor-skill-host-integration, catalog-publication-integrity, concept-content-authoring, session-recovery-version-pinning
priority: 75
components: curriculum, mcp-server, tutor-skill, mastery
delivers: capability:agent-authored-catalog-drafts
---

# Spec: agent-authored-catalog-drafts

## 1. Intent

### Goal
Tornar a geração de conteúdo pelo agente-tutor conectado ao MCP um quarto
modo de seleção, de primeira classe e escolhido deliberadamente pelo aluno
(trilha nova, trilha de N desafios ou desafio único, mesmo sobre assunto já
coberto), sem nunca publicar automaticamente e sem contaminar a evidência de
maestria revisada.

### Business value
Sem efetivo de autoria dedicado, o crescimento do catálogo depende de mim
(agente) e do usuário revisando; permitir que o mesmo aconteça em sessão
real — puxado pelo próprio uso, não só por autoria offline dos packs
`go-*` — acelera cobertura sem abrir mão do gate de qualidade já existente.

### Constraints
- Ver `.pose/adr/2026-08-23-agent-generated-catalog-content-as-a-first-class-mode.md`.
- O servidor MCP continua sem LLM embutido (`agent-mcp-and-core-boundaries`
  Decision 2); quem gera é o agente do host, o MCP só recebe e valida.
- Todo rascunho entra como `publication.status: draft`, nunca `published`
  direto.
- Reusa as mesmas regras de `internal/curriculum.Validate`/
  `RunEditorialChecks` já existentes — nenhuma regra de qualidade duplicada
  ou mais fraca para conteúdo gerado.

### Non-goals
- Definir a curadoria em duas camadas da fase 2 (eu primeiro revisor, usuário
  segundo) — essa spec cobre só a fase 1 (revisor único = usuário).
- Promoção automática de rascunho a `published` sem ação humana explícita.

## 2. Requirements

### Functional
- R1: Nova tool MCP (`content_draft_submit` ou equivalente) recebe um
  rascunho de desafio, trilha de N desafios ou trilha completa gerado pelo
  agente do host, no mesmo schema de autoria (`schemas/challenge.schema.json`
  e correlatos).
- R2: O rascunho é validado com as mesmas regras estruturais/editoriais já
  existentes antes de ser aceito; achados bloqueantes rejeitam a submissão
  com diagnóstico acionável (mesmo formato de `catalog validate`).
- R3: Rascunho aceito entra em área de quarentena do catálogo, visível para
  revisão, nunca listado como conteúdo `published` em buscas/recomendações
  padrão.
- R4: O aluno pode optar por rodar uma sessão sobre um rascunho ainda não
  revisado, com aviso persistente e explícito de que não é conteúdo oficial
  curado.
- R5: Evidência de sessão sobre rascunho é registrada com
  `content_provenance: draft`; a projeção de mastery (`mastery-review-
  scheduling`) ignora evidência com essa proveniência ao calcular estado
  revisado.
- R6: Promoção de rascunho a `published` segue exatamente o funil de
  `author`/`reviewed_by`/`playtested` de `catalog-authoring-quality` — sem
  atalho novo.
- R7: Persistir submissão idempotente, versão/digest e proveniência no servidor; o cliente não pode declarar conteúdo revisado, editar estado de sessão nem sobrescrever outro rascunho por colisão de ID.
- R8: Separar validação estrutural de submissão dos requisitos exclusivos de publicação; aceitar draft sem reviewed_by/playtested, rejeitando a tentativa de enviar status published.
- R9: Definir payload permitido para fixtures e checks: aceitar somente dados declarativos limitados; não executar código durante submissão, não materializar arquivos no workspace do aluno e não confundir runner allowlisted com isolamento de código adversarial.
- R10: Manter proveniência imutável por tentativa; promoção do conteúdo não promove retroativamente mastery. Qualquer revalidação posterior exige operação explícita e novas evidências causais, após decisão no ADR.
- R11: Incluir quarentena na política de quotas, retenção, export/remoção e retomada; provar isolamento entre sessões e exclusão do gate editorial V1 até publicação humana.

### Non-functional
- Validação de rascunho é síncrona à submissão (mesmo orçamento de
  `catalog validate` local).

### Security
- Rascunho gerado passa pelas mesmas regras de confinamento de fixture/path
  e scanner de secrets já aplicadas a conteúdo autorado normalmente.
- Tool de submissão não executa payloads. O contrato distingue texto de
  fixture de comandos; persistência ocorre em quarentena controlada, sem
  conceder escrita no workspace do aluno. Execução posterior exige autorização.

### Compatibility
- Adicionar proveniência sem reescrever eventos. Eventos legados sem campo
  não são automaticamente prova de publicação: o catálogo atual já permite
  sessões sobre itens não publicados. Definir migração/classificação no ADR
  e testar a projeção conservadora antes de implementar.

## 3. Technical Plan

### Affected areas
- internal/mcpserver/ (nova tool)
- internal/curriculum/ (área de quarentena, reuso de validação)
- internal/mastery/ (campo de proveniência)
- internal/session/ (sessão sobre rascunho + aviso)
- .agents/skills/codinho/ (exposição dos quatro modos)

### Artifacts
- created: internal/curriculum/drafts.go
- created: internal/curriculum/drafts_test.go
- created: internal/drafts/service.go
- created: internal/drafts/service_test.go
- created: internal/application/drafts.go
- created: internal/application/draft_provenance.go
- created: internal/mcpserver/content_draft_tools.go
- created: internal/mcpserver/content_draft_contract_test.go
- created: cmd/codinho/draft_integration_test.go
- modified: internal/curriculum/selection.go
- modified: internal/curriculum/loader.go
- modified: internal/mastery/model.go
- modified: internal/mastery/projector.go
- modified: internal/mastery/rules.go
- modified: internal/mastery/projector_test.go
- modified: internal/session/service.go
- modified: internal/session/recovery.go
- modified: internal/session/tracks.go
- modified: internal/application/session.go
- modified: internal/application/progress.go
- modified: internal/application/progress_test.go
- modified: internal/mcpserver/server.go
- modified: internal/mcpserver/session_tools.go
- modified: internal/mcpserver/progress_tools.go
- modified: internal/mcpserver/errors.go
- modified: internal/eventstore/event.go
- modified: internal/security/policy_test.go
- modified: internal/cli/privacy.go
- modified: cmd/codinho/main.go
- modified: .agents/skills/codinho/SKILL.md
- modified: .pose/contracts/mcp-stdio.json
- modified: .pose/indexes/validation-matrix.json
- modified: docs/content-authoring.md
- modified: .pose/adr/2026-08-23-agent-generated-catalog-content-as-a-first-class-mode.md

### Delivery targets
- capability:agent-authored-catalog-drafts module:cmd/codinho profile:composed-capability entrypoint:cmd/codinho/main.go

### API/contract changes
- Nova tool `content_draft_submit`; `session_start` aceita sessão sobre
  rascunho pendente; eventos de mastery ganham `content_provenance`.

### Data/storage changes
- Área de quarentena para rascunhos (separada do catálogo `published`).
- Campo `content_provenance` em eventos de evidência de mastery.

### Technical risks
- Ver Risks do decision log
  `adr-agent-generated-catalog-content-as-a-first-class-mode-review`.

## 4. Tasks

### Planning
- [ ] Confirmar formato de área de quarentena (diretório separado vs.
      flag no mesmo pack) sem quebrar `catalog validate` padrão.

### Implementation
- [ ] Implementar `content_draft_submit` reusando validação existente.
- [ ] Implementar quarentena e exclusão de rascunho de buscas padrão.
- [ ] Implementar `content_provenance` em eventos e ignorá-lo na projeção
      de mastery revisada.
- [ ] Implementar aviso persistente de sessão sobre rascunho.
- [ ] Expor os quatro modos na skill `codinho`.

### Validation
- [ ] Teste de submissão de rascunho válido e inválido (reuso de regras).
- [ ] Teste de sessão sobre rascunho não contaminando mastery revisada.
- [ ] Teste de promoção de rascunho a `published` pelo funil existente.
- [ ] Suíte completa de internal/curriculum, internal/mastery,
      internal/session, internal/mcpserver sem regressão.

## 5. Decisions

### Decision 2
- Date: 2026-09-09
- Decision: quarentena lógica no stream drafts do event log, packs YAML autocontidos; proveniência derivada de eventos de origem, nunca do cliente.
- Rationale: consumir knowledge:adr-agent-generated-catalog-content-as-a-first-class-mode-review e knowledge:planning-audit-2026-09; reutilizar durabilidade, export e purge existentes.
- Consequences: limite 256 KiB por submissão, 100 submissões/16 MiB acumulados até purge, validade de 30 dias para novos inícios; tombstone remove disponibilidade, histórico permanece até purge explícito. Promoção offline usa o funil editorial existente; não altera evidências anteriores.


### Decision 1
- Date: 2026-08-23
- Context: onde mora conteúdo gerado em sessão real antes de revisão.
- Options considered: (a) publicar direto, confiando no agente; (b)
  descartar ao fim da sessão (nunca persistir); (c) quarentena persistente
  sujeita ao mesmo funil de revisão da autoria offline.
- Decision: (c).
- Rationale: preserva o gate de `catalog-authoring-quality` e converte todo
  gap de cobertura em matéria-prima de curadoria, sem perder o trabalho do
  agente nem autocertificar.
- Consequences: exige área de armazenamento nova e um passo de promoção
  explícito, mas nenhuma regra de qualidade nova além da já existente. Ver
  `.pose/adr/2026-08-23-agent-generated-catalog-content-as-a-first-class-mode.md`.

## 6. Validation

### Strategy
Corpus de rascunhos válidos/inválidos, sessão real sobre rascunho
verificando isolamento de evidência, e promoção ponta a ponta até
`published`.

| Cenário obrigatório | Comando | Evidência |
|---|---|---|
| Validação | go test ./internal/curriculum ./internal/drafts | dados limitados, schema fechado, paths, secrets, status e critérios editoriais |
| Durabilidade | go test ./internal/drafts ./internal/session | retry, colisão, quota, expiração, remoção, restart e consentimento |
| Maestria | go test ./internal/mastery ./internal/application | draft/legado não promovem, publicado exige origem registrada, sem promoção retroativa |
| MCP real | go test ./cmd/codinho -run TestDraftOverRealStdio -count=1 | submissão, isolamento de busca, aviso persistente, replay e export/purge |
| Gate completo | pose validate --strict --json .pose/results/delivery-validation.json | checks required sem skips |

### Deterministic checks
- Test: go test ./internal/mcpserver/... ./internal/curriculum/... ./internal/mastery/... ./internal/session/...
- Lint: gofmt -l internal/mcpserver internal/curriculum internal/mastery internal/session
- Typecheck: go vet ./internal/mcpserver/... ./internal/curriculum/... ./internal/mastery/... ./internal/session/...
- Build: go build ./...
- Security / Contract: fixture confinement e secret scan já existentes, aplicados ao caminho de rascunho.

### Execution log
- Pendente.

### Results summary
Nenhuma implementação ainda; spec criada para sequenciar o trabalho.

### Requirement trace
- Mapear R1–R11 a testes de submissão, quarentena, proveniência, promoção, compatibilidade e privacidade.

### Known gaps
- Resolver no ADR a semântica de evidência anterior à publicação e o
  armazenamento de fixtures antes do primeiro incremento (R7–R11).

## 7. Final Report

### Delivered scope
Nenhum; spec draft aguardando implementação.

### Files and modules changed
- Planejados nas áreas afetadas acima.

### Validation executed
- Command: pose lint-spec agent-authored-catalog-drafts --ready-check
- Result: registrar após validação.

### Residual risks
- Nenhum adicional além do já descrito em Technical risks.

### Follow-ups
- [covered: v1-integrated-acceptance] Exercitar submissão, consentimento, retomada, promoção humana e isolamento de mastery no candidato V1.
