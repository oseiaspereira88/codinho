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
- created: .codex/config.toml.example

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
- [ ] Mapear cada regra da seção 16 de PROJECT.md para instrução ou referência.
- [ ] Definir corpus adversarial e critérios de passagem.

### Implementation
- [ ] Criar skill, metadados e referências.
- [ ] Implementar roteamento de tools e fallback.
- [ ] Configurar MCP de exemplo sem paths privados.
- [ ] Executar sessões no CLI e IDE.
- [ ] Capturar transcripts sintéticos e resultados.

### Validation
- [ ] Executar pose skills-check --strict.
- [ ] Executar contract/e2e com MCP real.
- [ ] Verificar zero edição e zero revelação indevida.
- [ ] Executar pose assess integrate e surface-check.

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
- Pendente.

### Results summary
- Skill e integração ainda não existem.

### Requirement trace
- Mapear R1–R10 a checks de conformance, transcripts e sessões host.

### Known gaps
- Variabilidade de modelos exige testes por comportamento, não texto exato.

## 7. Final Report

### Delivered scope
Nenhum; spec draft.

### Files and modules changed
- Planejados em .agents/skills/codinho, testdata/host e config de exemplo.

### Validation executed
- Command: pose lint-spec tutor-skill-host-integration --ready-check
- Result: registrar após gate.

### Residual risks
- Compatibilidade futura de host depende da documentação oficial.

### Follow-ups
- [covered: v1-integrated-acceptance] Revalidar o workflow composto no gate final.
