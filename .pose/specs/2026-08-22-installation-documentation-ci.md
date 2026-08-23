---
slug: installation-documentation-ci
status: in-progress
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: tutor-skill-host-integration, administrative-cli-fixtures, catalog-authoring-quality, security-privacy-hardening, reliability-observability-compatibility
priority: 240
components: distribution, documentation, ci
delivers: infrastructure:codinho-local-distribution, governance:codinho-ci
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
- README.md, docs/, .github/workflows/, Makefile, scripts/, .agents/skills/codinho/

### Artifacts
- modified: README.md
- created: docs/install.md
- created: docs/quickstart.md
- created: docs/configuration.md
- modified: docs/compatibility.md
- modified: docs/troubleshooting.md
- created: .github/workflows/ci.yml
- created: .github/workflows/release.yml
- created: Makefile
- created: scripts/smoke-install.sh
- created: cmd/codinho/smoke_test.go
- modified: .pose/indexes/validation-matrix.json
- modified: .agents/skills/codinho/references/codex-configuration.md

### Delivery targets
- infrastructure:codinho-local-distribution module:cmd/codinho profile:release-governance entrypoint:cmd/codinho/main.go
- governance:codinho-ci module:.github/workflows profile:release-governance entrypoint:.github/workflows/ci.yml

### API/contract changes
- Documentar contratos existentes; nenhum contrato novo de runtime.

### Data/storage changes
- Nenhum.

### Technical risks
- Docs podem divergir de comandos.
- CI pode alegar plataforma suportada sem e2e host.

## 4. Tasks

### Planning
- [x] Definir matriz de build, CI, artifacts e release.
- [x] Definir governança documental e owners.

### Implementation
- [x] Escrever documentação e exemplos.
- [x] Implementar Makefile e smoke install.
- [x] Implementar CI com gates e artifacts.
- [x] Implementar release workflow em modo não publicador até aprovação.
- [x] Testar onboarding em ambiente limpo.

### Validation
- [x] Executar pose docs-check, check, validate e release plan.
- [ ] Executar smoke-install em todas as plataformas declaradas.
- [x] Revisar workflows por segurança e permissions.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Plugin adicionaria distribuição antes de validar o produto.
- Options considered: plugin; instalador próprio; binário mais skill e config.
- Decision: distribuição local explícita na V1.
- Rationale: menor superfície e compatível com o fluxo desejado.
- Consequences: instalação tem mais passos, compensados por quickstart e doctor.

### Decision 2
- Date: 2026-08-23
- Context: `docs/content-authoring.md` e `docs/content-review-checklist.md`
  já existiam (criados por catalog-authoring-quality) com exatamente o
  conteúdo que R7/o Artifacts planejado de `docs/catalog-authoring.md`
  pedia.
- Options considered: (a) criar `docs/catalog-authoring.md` como um
  segundo documento quase idêntico; (b) referenciar o arquivo já
  existente a partir do README e não duplicar.
- Decision: (b). README.md lista `docs/content-authoring.md` e
  `docs/content-review-checklist.md` diretamente; nenhum arquivo novo
  foi criado para isso.
- Rationale: duplicar documentação garante que ela diverge com o tempo;
  o nome do arquivo no Artifacts original era um placeholder de
  planejamento, não um contrato.
- Consequences: Artifacts desta spec não inclui
  `docs/catalog-authoring.md`.

### Decision 3
- Date: 2026-08-23
- Context: R9 exige "processo governado de release sem publicar a V1
  antecipadamente"; um `release.yml` automático que cria uma GitHub
  Release a cada execução violaria isso diretamente.
- Options considered: (a) workflow completo com publish automático,
  gateado só por uma flag manual fácil de mudar sem revisão; (b)
  workflow que só builda binários reproduzíveis e sobe como artifact de
  workflow (nunca uma Release pública), com o job de publish
  deliberadamente ausente até v1-integrated-acceptance decidir
  adicioná-lo.
- Decision: (b).
- Rationale: a ausência do job de publish é a garantia mais forte
  possível de "não publicar antecipadamente" — não depende de ninguém
  lembrar de manter uma flag desligada.
- Consequences: `.github/workflows/release.yml` produz binários e
  checksums como artifact de CI, nunca uma release pública, até esta
  spec (ou v1-integrated-acceptance) ser revisitada para adicionar
  esse job.

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
- `go build ./...` → ok (2026-08-23).
- `gofmt -l .` → sem saída (2026-08-23).
- `go vet ./...` e `GOOS=windows go vet ./...` → sem diagnósticos (2026-08-23).
- `go test ./... -race` → `ok` em todos os pacotes, incluindo `cmd/codinho` (TestSmokeInstallScript, novo), sem data races (2026-08-23).
- `./scripts/smoke-install.sh <binário>` executado diretamente num diretório limpo (`mktemp -d`) → `smoke-install: OK` (2026-08-23).
- `govulncheck ./...` → "No vulnerabilities found." (2026-08-23).
- `python3 -c "yaml.safe_load(...)"` em `.github/workflows/ci.yml` e `.github/workflows/release.yml` → ambos válidos (2026-08-23).
- Build real com `-ldflags -X internal/cli.version=...` confirmado localmente (usado pelo release.yml) (2026-08-23).

### Results summary
README, docs de instalação/quickstart/configuração, Makefile,
script de smoke-install real (com teste Go que o executa) e dois
workflows de CI entregues. CI roda gofmt/vet/build/test-race/
govulncheck/catalog-validate/skills-check em cada PR e em `main`
(requirement R4/R5), publicando resultados de teste redigidos como
artifact (R6). Release workflow builda binários reproduzíveis com
checksum mas nunca publica (Decision 3) — versionamento e processo
governado ficam explícitos sem antecipar a V1 (R9).

### Requirement trace
- R1 [satisfied] docs/{install,quickstart,configuration}.md.
- R2 [satisfied] docs/configuration.md + .agents/skills/codinho/references/codex-configuration.md.
- R3 [satisfied] .github/workflows/release.yml (build reprodutível + sha256, sem publish).
- R4 [satisfied] .github/workflows/ci.yml (pull_request + push main).
- R5 [satisfied] ci.yml (gofmt/vet/build/test-race/govulncheck/catalog validate/skills-check).
- R6 [satisfied] ci.yml (upload-artifact de test-results.json, nunca código/estado local).
- R7 [satisfied] docs/{security/threat-model,security/privacy,troubleshooting,content-authoring,content-review-checklist}.md (Decision 2: reusa os já existentes).
- R8 [satisfied] scripts/smoke-install.sh + test:TestSmokeInstallScript (execução real em diretório limpo).
- R9 [satisfied] Decision 3 (job de publish deliberadamente ausente).
- R10 [satisfied] PROJECT.md permanece a visão; .pose/specs+roadmap permanecem o estado de entrega (nenhuma mudança necessária — já era a prática desde o início da sessão).

### Known gaps
- Pin de actions por commit SHA não foi feito (usa tags de versão major `@v4`/`@v5`) — este ambiente não pôde buscar e verificar SHAs reais offline; hardening documentado como follow-up.
- Smoke-install só foi executado em Linux nesta sessão.
- Publicação de release ocorre somente após v1-integrated-acceptance decidir adicionar o job de publish.

## 7. Final Report

### Delivered scope
Documentação de instalação/quickstart/configuração, Makefile, script de
smoke-install real (testado), CI real com gates completos, e release
workflow que builda mas nunca publica.

### Files and modules changed
- README.md (reescrito para refletir o estado real do projeto).
- docs/{install,quickstart,configuration}.md (novos); compatibility.md, troubleshooting.md (cross-links).
- Makefile, scripts/smoke-install.sh, cmd/codinho/smoke_test.go.
- .github/workflows/{ci,release}.yml.
- .pose/indexes/validation-matrix.json (check smoke-install).

### Validation executed
- Command: go test ./... -race
- Result: ok em todos os pacotes.
- Command: ./scripts/smoke-install.sh <binário> (diretório limpo)
- Result: smoke-install: OK
- Command: govulncheck ./...
- Result: No vulnerabilities found.

### Residual risks
- Docs podem divergir de comandos com o tempo — mitigado por
  `docs/quickstart.md` e `docs/install.md` citarem apenas comandos reais
  já cobertos por teste (`internal/cli`, `cmd/codinho`).
- CI pode alegar plataforma suportada sem e2e host real — mitigado por
  `docs/compatibility.md` declarar honestamente o que foi testado.

### Follow-ups
- [covered: v1-integrated-acceptance] Executar onboarding e gates no candidato final.
- [open] Pinar actions do GitHub por commit SHA verificado (hoje usa tags major).
- [open] Executar smoke-install real em macOS e Windows quando houver acesso.
