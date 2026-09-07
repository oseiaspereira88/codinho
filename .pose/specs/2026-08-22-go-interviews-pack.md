---
slug: go-interviews-pack
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: curriculum-graph-path-recommendation, catalog-authoring-quality, interview-mode, go-production-architecture-packs
priority: 210
components: curriculum-content, interview-mode
delivers:
---

# Spec: go-interviews-pack

## 1. Intent

### Goal
Publicar quatro simulações técnicas genéricas que completem os 84 desafios e exercitem autonomia, comunicação e tempo.

### Business value
Permitir prática realista de live coding e discussão técnica sem associação a terceiros ou material proprietário.

### Constraints
- Publicar exatamente 4 desafios do tipo simulation.
- Cada simulação dura entre 45 e 90 minutos e usa o modo entrevista.
- Briefings, rubricas e fixtures são originais e genéricos.

### Non-goals
- Prever perguntas de uma organização.
- Certificar senioridade ou decisão de contratação.
- Proctoring.

## 2. Requirements

### Functional
- R1: Publicar o pack go-interviews com quatro simulações versionadas.
- R2: Cobrir transformação de dados e testes, concorrência/contexto, depuração/revisão e desenho de fatia backend.
- R3: Exigir clarificação, casos de borda, complexidade e raciocínio em voz alta.
- R4: Fornecer briefing público, critérios, checks permitidos e rubrica reservada.
- R5: Oferecer variantes suficientes para revisão sem memorizar solução.
- R6: Publicar duas trilhas: preparação para live coding e backend Go sênior.
- R7: Validar duração e dificuldade por ao menos dois playtests independentes por simulação.
- R8: Gerar relatório final sem score único e sem linguagem de seleção.
- R9: Fechar a distribuição total da V1 em 42 atômicos, 20 combinados, 8 fatias, 6 depuração/refatoração, 4 missões e 4 simulações.

### Non-functional
- Simulações devem ser executáveis offline e sem serviços externos.
- Variantes mantêm competência e mudam contexto substancialmente.

### Security
- Fixtures são sintéticas e não coletam dados pessoais.

### Compatibility
- Checks respeitam as plataformas declaradas pela V1.

## 3. Technical Plan

### Affected areas
- packs/go-interviews.yaml, testdata/packs/interviews/, docs/catalog/

### Artifacts
- modified: packs/manifest.yaml
- created: packs/go-interviews.yaml
- created: testdata/packs/interviews/
- created: docs/catalog/go-interviews.md
- created: docs/catalog/interview-playtest-report.md

### Delivery targets
Nenhum tipado; conteúdo consumido pelo modo entrevista.

### API/contract changes
- Nenhuma mudança; usar contratos de interview-mode.

### Data/storage changes
- Adicionar conteúdo, fixtures e relatórios de playtest sanitizados.

### Technical risks
- Dificuldade e duração podem variar muito.
- Rubrica reservada pode vazar pelo catálogo.

## 4. Tasks

### Planning
- [ ] Definir quatro matrizes de competência e variações.
- [ ] Definir rubricas, checks e política de disclosure.

### Implementation
- [ ] Autorar simulações e variantes.
- [ ] Criar fixtures e checks.
- [ ] Criar duas trilhas.
- [ ] Executar oito ou mais playtests.
- [ ] Ajustar duração, dificuldade e rubricas.

### Validation
- [ ] Executar catálogo e checks.
- [ ] Verificar isolamento da rubrica reservada.
- [ ] Auditar distribuição total dos 84 desafios.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Realismo pode induzir alegação de vínculo externo.
- Options considered: imitar empresas; banco genérico; simulações originais por competência.
- Decision: simulações originais, neutras e explicitamente pedagógicas.
- Rationale: preserva utilidade sem associação indevida.
- Consequences: marketing e documentação devem manter a mesma neutralidade.

## 6. Validation

### Strategy
Combinar validação de catálogo, leak tests, checks e playtests cronometrados.

### Deterministic checks
- Test: checks das quatro simulações e variantes.
- Lint: go run ./cmd/codinho catalog validate --json (executar na raiz; catálogo resolvido pelo manifest).
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: disclosure tests, privacy review e synthetic-data scan.

### Execution log
- Pendente.

### Results summary
- Nenhuma simulação publicada.

### Requirement trace
- Mapear R1–R9 a manifests, playtests, coverage e reports.

### Known gaps
- Dificuldade será calibrada novamente após o piloto V1.

## 7. Final Report

### Delivered scope
Nenhum; spec draft.

### Files and modules changed
- Planejados no pack, testdata e docs.

### Validation executed
- Command: pose lint-spec go-interviews-pack --ready-check
- Result: registrar após gate.

### Residual risks
- Rubricas exigem revisão humana para consistência.

### Follow-ups
- [covered: v1-integrated-acceptance] Executar uma simulação completa no gate final.
