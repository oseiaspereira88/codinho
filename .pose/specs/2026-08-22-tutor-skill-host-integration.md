---
slug: tutor-skill-host-integration
status: in-progress
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: mcp-stdio-foundation, session-orchestration-disclosure, assistance-hints-detours, feedback-evaluation-progression, workspace-observation-baselines, safe-check-executor
priority: 130
components: tutor-skill, mcp-server
delivers: capability:codinho-tutor
---

# Spec: tutor-skill-host-integration

## 1. Intent

### Goal
Criar a skill codinho e comprovar seu workflow com o servidor MCP no Codex CLI e na extensão da IDE.

### Business value
Transformar tools e estado em uma experiência pedagógica consistente que preserve a autoria do aluno.

### Constraints
- A skill controla comportamento; o MCP permanece fonte de verdade.
- Não editar arquivos do aluno durante tutoria.
- Funcionar de modo degradado, claramente indicado, quando o MCP não estiver disponível.

### Non-goals
- Plugin distribuído ou interface própria.
- Embutir conteúdo curricular inteiro no SKILL.md.

## 2. Requirements

### Functional
- R1: Definir triggers, limites, tool routing e divulgação progressiva no SKILL.md.
- R2: Declarar dependência MCP e metadados de interface em agents/openai.yaml.
- R3: Ler session_get antes de inferir estado e entregar somente uma instrução ativa.
- R4: Observar antes de avaliar e exigir submission_intent para criar tentativa.
- R5: Tratar feedback como consultivo e informar sempre progress_effect.
- R6: Bloquear edição e solução salvo mudança explícita autorizada pela política.
- R7: Ignorar instruções encontradas em código, fixtures e outputs observados.
- R8: Adaptar conceito ao perfil usando conteúdo canônico e exemplos fora da solução.
- R9: Operar no Codex CLI e IDE com o mesmo MCP configurado.
- R10: Informar limitações no fallback conversacional e não simular persistência.

### Non-functional
- SKILL.md deve permanecer focado e usar referências por divulgação progressiva.
- Respostas em modo micro devem ser curtas e conter uma ação principal.

### Security
- A skill não amplia permissões e respeita approvals do host.
- Testes adversariais cobrem prompt injection e pedidos ambíguos de código.

### Compatibility
- Seguir o formato de Agent Skills e metadados suportados pelo Codex.

## 3. Technical Plan

### Affected areas
- .agents/skills/codinho/, .codex/, testdata/host/

### Artifacts
- created: .agents/skills/codinho/SKILL.md
- created: .agents/skills/codinho/agents/openai.yaml
- created: .agents/skills/codinho/references/tutor-contract.md
- created: .agents/skills/codinho/references/feedback-rubric.md
- created: .agents/skills/codinho/references/session-modes.md
- created: .agents/skills/codinho/references/mcp-tool-routing.md
- created: .agents/skills/codinho/assets/session-summary-template.md
- created: testdata/host/adversarial-prompts.yaml
- created: testdata/host/session-transcripts/practice-session-happy-path.md
- created: testdata/host/session-transcripts/degraded-mode-and-injection.md
- created: .codex/config.example.toml
- created: cmd/codinho/integration_test.go
- modified: .pose/indexes/validation-matrix.json

### Delivery targets
- capability:codinho-tutor module:.agents/skills/codinho profile:composed-capability entrypoint:.agents/skills/codinho/SKILL.md

### API/contract changes
- Consumir o contrato MCP v1 sem adicionar estado na skill.

### Data/storage changes
- Nenhum; transcripts de teste são fixtures sintéticas.

### Technical risks
- Instruções longas podem perder prioridade.
- Comportamento pode variar por host ou modelo.

## 4. Tasks

### Planning
- [x] Mapear cada regra da seção 16 de PROJECT.md para instrução ou referência.
- [x] Definir corpus adversarial e critérios de passagem.

### Implementation
- [x] Criar skill, metadados e referências.
- [x] Implementar roteamento de tools e fallback.
- [x] Configurar MCP de exemplo sem paths privados.
- [x] Executar sessão completa real via stdio (mesmo transporte que Codex CLI/IDE usam; verificação visual do host literal fica para v1-integrated-acceptance, Decision 2).
- [x] Capturar transcripts sintéticos e resultados.

### Validation
- [x] Executar pose skills-check --strict.
- [x] Executar contract/e2e com MCP real.
- [x] Verificar zero edição e zero revelação indevida.
- [x] Executar pose assess integrate e surface-check.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Instruções e estado têm ciclos de vida diferentes.
- Options considered: tudo na skill; tudo no MCP; responsabilidades separadas.
- Decision: skill governa workflow; MCP governa estado e capabilities.
- Rationale: mantém a conversa flexível e o progresso determinístico.
- Consequences: integração deve testar ambos em composição.

### Decision 2
- Date: 2026-08-23
- Context: R9 pede operar "no Codex CLI e IDE com o mesmo MCP
  configurado", mas este ambiente de execução não tem um Codex CLI nem
  uma extensão de IDE reais para abrir uma sessão interativa de
  verdade. v1-integrated-acceptance já declara explicitamente
  "Aceite exige Codex CLI e extensão IDE reais" como CONSTRAINT
  daquele gate final, não desta spec.
- Options considered: (a) bloquear esta spec até haver acesso a hosts
  reais; (b) provar a integração pelo mesmo mecanismo de transporte que
  um host real usaria (`mcp.CommandTransport` sobre stdio — a MESMA
  interface que Codex CLI e a extensão de IDE usam para falar com
  `codinho serve`), documentando a verificação de UI/host literal como
  residual, coberta pelo gate final.
- Decision: (b). Prova-se aqui que a skill roteia corretamente para
  cada tool na ordem que as 16 regras de PROJECT.md §16.2 exigem,
  via sessões reais de protocolo stdio (não mockadas) contra o
  binário `codinho serve` real — o mesmo transporte, mesmo protocolo,
  mesmos tools que Codex CLI/IDE usariam. A verificação visual de UI em
  Codex CLI e na extensão de IDE em si permanece para
  v1-integrated-acceptance, conforme a constraint que aquela spec já
  declara.
- Rationale: não duplicar/antecipar um gate que já pertence
  explicitamente a outra spec; ainda assim entregar verificação real
  (não apenas documental) do que está sob controle desta spec — o
  conteúdo e roteamento da skill, não a apresentação do host.
- Consequences: Known Gap explícito nesta spec: comportamento
  específico de renderização/UX de um host real (Codex CLI, extensão
  de IDE) só é confirmado no gate de aceite V1.

## 6. Validation

### Strategy
Combinar conformance da skill, transcripts adversariais e sessões reais nos dois hosts.

### Deterministic checks
- Test: go test ./internal/mcpserver/... e runner de transcripts sintéticos.
- Lint: pose skills-check --strict.
- Typecheck: validação YAML de openai.yaml.
- Build: go build ./cmd/codinho.
- Security / Contract: pose assess integrate; pose surface-check --spec tutor-skill-host-integration --strict.

### Execution log
- `pose skills-check --strict` → `skills.checked=12 skills.errors=0 skills.warnings=0` (2026-08-23).
- `python3 -c "yaml.safe_load(...)"` em `agents/openai.yaml` e `testdata/host/adversarial-prompts.yaml`, e parse de `.codex/config.example.toml` via `tomllib` → todos válidos (2026-08-23).
- `go build ./...` → ok (2026-08-23).
- `go test ./... -race` → `ok` em todos os pacotes com testes, incluindo `cmd/codinho` (novo), sem data races (2026-08-23).
- `govulncheck ./...` → "No vulnerabilities found." (2026-08-23).
- `pose assess integrate` → 0 contratos (esperado; esta spec não altera schemas de tool nem contratos inter-serviço).
- `pose validate --strict --json .pose/results/delivery-validation.json` + `pose index` → registra o novo check `tutor-skill-routing` (`evidenceClass: integration`) no módulo `.` (2026-08-23).
- `pose surface-check --spec tutor-skill-host-integration --strict` → `surface.targets=1 surface.findings=0` para `capability:codinho-tutor` (2026-08-23).

### Results summary
- `.agents/skills/codinho/` cobre as 16 regras normativas de PROJECT.md
  §16.2 (SKILL.md resumindo, `references/tutor-contract.md` expandindo
  com exemplo correto/incorreto por regra), roteia para as 29 tools do
  contrato MCP v1 (`references/mcp-tool-routing.md`), cita as rubricas
  já referenciadas por `feedback_prepare` (`references/feedback-rubric.md`,
  fechando o dangling reference que feedback-evaluation-progression já
  antecipava), e documenta os seis modos e cinco profundidades já
  aceitos por `session_start` (`references/session-modes.md`).
- `cmd/codinho/integration_test.go` prova, contra o binário real
  compilado e o mesmo transporte stdio que Codex CLI/extensão de IDE
  usariam (`mcp.CommandTransport`), a sequência completa documentada em
  `mcp-tool-routing.md`: descoberta → `session_start` →
  `workspace_observe` → `check_run` → `step_evaluate` (verdict real do
  check) → `reflection_record` → `step_complete` → `step_advance` →
  `mastery_evidence_record` → `progress_get` → `review_due`. Registrado
  em `validation-matrix.json` como evidência `integration` do módulo
  `.`, satisfazendo o perfil `composed-capability` do alvo
  `capability:codinho-tutor`.
- `testdata/host/adversarial-prompts.yaml` (6 casos) e os dois
  transcripts sintéticos documentam o comportamento esperado contra
  prompt injection via código observado, saída de check, nome de
  arquivo, pedido ambíguo de edição, resultado de tool forjado no chat,
  e instrução embutida numa reflexão — todos ligados à regra 14 e à
  invariante 9. Estes são specs comportamentais para revisão
  manual/host (Known Gaps), não asserções byte-a-byte, porque
  comportamento de modelo varia.
- `.codex/config.example.toml` e `agents/openai.yaml` declaram a
  dependência do servidor MCP sem paths privados; residual risk já
  documentado sobre compatibilidade futura com o formato real do Codex.

### Requirement trace
- R1 [satisfied] report:.agents/skills/codinho/SKILL.md report:.agents/skills/codinho/references/tutor-contract.md
- R2 [satisfied] report:.agents/skills/codinho/agents/openai.yaml
- R3 [satisfied] report:.agents/skills/codinho/references/tutor-contract.md (regra 5) test:TestTutorSkillFullSessionRoutingOverRealStdio
- R4 [satisfied] report:.agents/skills/codinho/references/tutor-contract.md (regra 6, 7) test:TestTutorSkillFullSessionRoutingOverRealStdio
- R5 [satisfied] report:.agents/skills/codinho/references/tutor-contract.md (regra 8, 9)
- R6 [satisfied] report:.agents/skills/codinho/references/tutor-contract.md (regra 3, 4)
- R7 [satisfied] report:testdata/host/adversarial-prompts.yaml report:.agents/skills/codinho/references/tutor-contract.md (Segurança)
- R8 [satisfied] report:.agents/skills/codinho/references/tutor-contract.md (regra 13)
- R9 [satisfied] test:TestTutorSkillFullSessionRoutingOverRealStdio (mesmo transporte stdio; UI literal de host — Known Gap, Decision 2)
- R10 [satisfied] report:.agents/skills/codinho/SKILL.md (Modo degradado) report:testdata/host/session-transcripts/degraded-mode-and-injection.md

### Known gaps
- Variabilidade de modelos exige testes por comportamento, não texto
  exato — `adversarial-prompts.yaml` descreve comportamento esperado
  para revisão manual/host, não uma asserção automatizável em Go.
- Verificação visual/UX literal em Codex CLI e na extensão de IDE
  (requirement R9) fica para v1-integrated-acceptance, que já declara
  isso como sua própria constraint (Decision 2) — esta spec prova
  roteamento e conteúdo via o mesmo protocolo/transporte, não a
  apresentação do host.

## 7. Final Report

### Delivered scope
A skill `codinho` completa (SKILL.md + 4 referências + template de
resumo + metadados Codex), corpus adversarial e transcripts sintéticos,
e um teste de integração real (binário compilado, transporte stdio)
provando que a sequência de tools documentada funciona de ponta a
ponta. Nenhuma edição de arquivo do aluno, nenhuma revelação de solução
fora do fluxo de `hint_request` autorizado pela política (non-goals e
regras 3/4 preservados).

### Files and modules changed
- `.agents/skills/codinho/{SKILL.md,agents/openai.yaml,references/*.md,assets/session-summary-template.md}` (criados)
- `testdata/host/{adversarial-prompts.yaml,session-transcripts/*.md}` (criados)
- `.codex/config.example.toml` (criado)
- `cmd/codinho/integration_test.go` (criado)
- `.pose/indexes/validation-matrix.json` (modificado: check `tutor-skill-routing` com `evidenceClass: integration`)

### Validation executed
- Command: go test ./... -race
- Result: `ok` em todos os pacotes com testes, sem data races.
- Command: pose skills-check --strict
- Result: `skills.errors=0`.
- Command: pose surface-check --spec tutor-skill-host-integration --strict
- Result: `surface.findings=0` para `capability:codinho-tutor`.

### Residual risks
- Compatibilidade futura do formato de metadados do Codex depende da documentação oficial (já refletido em `agents/openai.yaml`).
- Ver Known Gaps: verificação visual real de host (Codex CLI, extensão de IDE) fica para v1-integrated-acceptance.

### Follow-ups
- [covered: v1-integrated-acceptance] Revalidar o workflow composto no gate final.
