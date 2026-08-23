---
slug: installation-documentation-ci
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: tutor-skill-host-integration, administrative-cli-fixtures, catalog-authoring-quality, security-privacy-hardening, reliability-observability-compatibility
priority: 240
components: distribution, documentation, ci
delivers:
---

# Spec: installation-documentation-ci

## 1. Intent

### Goal
Entregar instalação local reproduzível, documentação operacional e CI estrito para código, contratos, segurança e catálogo.

### Business value
Permitir que uma pessoa instale e use codinho no Codex CLI ou IDE sem conhecimento implícito do autor.

### Constraints
- Distribuição V1 é binário, skill e configuração local; não é plugin.
- Exemplos não contêm paths privados nem secrets.
- CI publica evidência versionável e aplica gates POSE.

### Non-goals
- Marketplace, auto-update, package managers ou assinatura de binário.
- Deploy remoto.

## 2. Requirements

### Functional
- R1: Documentar pré-requisitos, build, instalação, configuração MCP, skill, primeiro desafio, atualização e remoção.
- R2: Fornecer config de exemplo para Codex CLI e extensão IDE.
- R3: Produzir binários reproduzíveis para plataformas comprovadas e checksums.
- R4: Executar CI em pull request e main com pose check e validate estritos.
- R5: Executar gofmt, test, race, vet, build, govulncheck, skills, catálogo, contratos MCP e e2e aplicáveis.
- R6: Publicar resultados estruturados e logs redigidos como artifacts de CI.
- R7: Documentar autoria de packs, segurança, privacidade, troubleshooting e limitações.
- R8: Validar todos os links, comandos e exemplos em ambiente limpo.
- R9: Definir versionamento e processo governado de release sem publicar a V1 antecipadamente.
- R10: Manter PROJECT.md como visão e specs/roadmap como estado de entrega.

### Non-functional
- Onboarding deve ser completável apenas com documentação.
- CI deve usar versões fixadas e cache sem comprometer reprodutibilidade.

### Security
- Aplicar least privilege no workflow, pin de actions e secret scan.
- Não publicar evidências com código do aluno ou estado local.

### Compatibility
- Declarar somente combinações testadas de Go, OS, SDK, Codex e schema.

## 3. Technical Plan

### Affected areas
- README.md, docs/, .github/workflows/, Makefile, scripts/, .codex/

### Artifacts
- modified: README.md
- created: docs/install.md
- created: docs/quickstart.md
- created: docs/configuration.md
- created: docs/catalog-authoring.md
- modified: docs/compatibility.md
- modified: docs/troubleshooting.md
- created: .github/workflows/ci.yml
- created: .github/workflows/release.yml
- created: Makefile
- created: scripts/smoke-install.sh
- modified: .codex/README.md

### Delivery targets
Os alvos planejados são a infrastructure `codinho-local-distribution`, com
módulo `cmd/codinho` e entrypoint `cmd/codinho/main.go`, e a governance
`codinho-ci`, com módulo `.github/workflows` e entrypoint
`.github/workflows/ci.yml`; ambos usam o perfil `release-governance`.
Registrá-los como alvos tipados quando esses caminhos forem materializados,
antes do closeout desta spec.

### API/contract changes
- Documentar contratos existentes; nenhum contrato novo de runtime.

### Data/storage changes
- Nenhum.

### Technical risks
- Docs podem divergir de comandos.
- CI pode alegar plataforma suportada sem e2e host.

## 4. Tasks

### Planning
- [ ] Definir matriz de build, CI, artifacts e release.
- [ ] Definir governança documental e owners.

### Implementation
- [ ] Escrever documentação e exemplos.
- [ ] Implementar Makefile e smoke install.
- [ ] Implementar CI com gates e artifacts.
- [ ] Implementar release workflow em modo não publicador até aprovação.
- [ ] Testar onboarding em ambiente limpo.

### Validation
- [ ] Executar pose docs-check, check, validate e release plan.
- [ ] Executar smoke-install em todas as plataformas declaradas.
- [ ] Revisar workflows por segurança e permissions.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Plugin adicionaria distribuição antes de validar o produto.
- Options considered: plugin; instalador próprio; binário mais skill e config.
- Decision: distribuição local explícita na V1.
- Rationale: menor superfície e compatível com o fluxo desejado.
- Consequences: instalação tem mais passos, compensados por quickstart e doctor.

## 6. Validation

### Strategy
Usar CI real, smoke install limpo, docs checks e dry-run de release.

### Deterministic checks
- Test: make test e scripts/smoke-install.sh.
- Lint: pose docs-check; gofmt; catalog validate.
- Typecheck: go vet ./...
- Build: make build para a matriz suportada.
- Security / Contract: govulncheck; action pin audit; secret scan; MCP contract e skills-check.

### Execution log
- Pendente.

### Results summary
- Instalação, docs e CI não entregues.

### Requirement trace
- Mapear R1–R10 a CI runs, smoke reports, docs checks e release plan.

### Known gaps
- Publicação de release ocorre somente após v1-integrated-acceptance.

## 7. Final Report

### Delivered scope
Nenhum; spec draft.

### Files and modules changed
- Planejados em docs, workflows, scripts, Makefile e config.

### Validation executed
- Command: pose lint-spec installation-documentation-ci --ready-check
- Result: registrar após gate.

### Residual risks
- Mudanças externas do host exigirão revisão de compatibilidade.

### Follow-ups
- [covered: v1-integrated-acceptance] Executar onboarding e gates no candidato final.
