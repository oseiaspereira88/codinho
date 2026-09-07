---
slug: v1-delivery-ci-assurance
status: draft
created_at: 2026-09-07
completed_at:
supersedes:
depends_on: installation-documentation-ci, reliability-observability-compatibility
priority: 35
components: ci, governance, distribution
delivers: governance:v1-delivery-ci
---

# Spec: v1-delivery-ci-assurance

## 1. Intent

### Goal
Tornar os gates de entrega executáveis em CI e reconciliar proveniência, plataformas e documentação antes do aceite.

### Business value
Fechar uma lacuna verificável da V1 antes de ampliar a superfície que depende dela.

### Constraints
- Preservar autoria do aluno, operação local, JSONL e separação pedagógica.
- Usar o [relatório de auditoria](../reports/2026-09-07-doc-audit-auditoria-do-planejamento-v1.md) como baseline, não como evidência de entrega.
- Consumir knowledge:planning-audit-2026-09 antes de iniciar a implementação.

### Non-goals
- Reabrir ou reescrever atestações históricas de specs done.
- Publicar release, alterar código do aluno ou ampliar a V1 para nuvem/multiusuário.

## 2. Requirements

### Functional
- R1: Instalar versão fixada e verificada do POSE em CI e executar check/validate estritos, skills-check e resultados estruturados; ausência de ferramenta requerida falha o job.
- R2: Alinhar matriz, Makefile e CI para formatação com código de falha, race, vet, build, catálogo, contratos e fixtures; executar o gate quantitativo V1 apenas no candidato completo.
- R3: Executar testes e smoke-install nativos em Linux e macOS conforme RNF-005; Windows permanece desejável e não bloqueia V1 sem mudança de escopo. Distinguir build cruzado de execução por arquitetura.
- R4: Fixar actions e scanner por versão/digest verificado; tratar inputs de release como dados validados, sem interpolação direta em shell; manter build sem publicação automática.
- R5: Mapear componentes/contratos reais, habilitar políticas de artifacts/delivery com roots e evidências corretas, reconciliar rename ailearn→codinho sem falsificar histórico ou editar bundles imutáveis.
- R6: Registrar e validar critérios C1–C8 do roadmap por evidência atual; comprovar que stale, skipped, check inexistente e relatório ausente bloqueiam a entrega.
- R7: Governar documentação por manifest e revisar comandos/links; corrigir alegações de retomada e suporte até existirem evidências dos contratos correspondentes.
- R8: Produzir relatório de prontidão com matriz por requisito, plataforma, host, commit e resultado; obter revisão independente do threat model antes do aceite.

### Non-functional
- Manter determinismo e falhas explícitas; provar o caminho composto, não apenas helpers isolados.
- Identificar resultados por versão, commit e cenário; skips obrigatórios não são sucesso.

### Security
- Não copiar estado real do aluno para fixtures, logs ou relatórios.
- Aplicar confinamento de paths, redaction e consentimento nos novos caminhos.

### Compatibility
Preservar a política de release sem publish. Definir roots e equivalência de componentes com evidência; não desabilitar gates para fazer o closeout passar.

## 3. Technical Plan

### Affected areas
- ci
- governance
- distribution

### Artifacts
- modified: .github/workflows/ci.yml
- modified: .github/workflows/release.yml
- modified: Makefile
- modified: .pose/indexes/validation-matrix.json
- modified: .pose/indexes/module-metadata.json
- modified: .pose/policy/artifacts.json
- modified: .pose/policy/delivery.json
- modified: .pose/policy/docs.json
- modified: .pose/roadmaps/codinho-v1.md
- modified: docs/compatibility.md
- modified: docs/install.md
- modified: docs/quickstart.md
- created: docs/acceptance/v1-release-readiness.md

### Delivery targets
- governance:v1-delivery-ci module:.github/workflows profile:release-governance entrypoint:.github/workflows/ci.yml

### API/contract changes
Preservar a política de release sem publish. Definir roots e equivalência de componentes com evidência; não desabilitar gates para fazer o closeout passar.

### Data/storage changes
Documentar os campos e migração exigidos pelos requisitos antes do primeiro incremento. Nenhuma migração é executada nesta rodada de planejamento.

### Technical risks
Renames históricos e evidências de módulos divergentes bloqueiam roll-up. CI remota e revisão humana precisam de evidência real, não de YAML existente.

## 4. Tasks

### Planning
- [ ] Reproduzir o achado e revisar contratos/ADRs aplicáveis.
- [ ] Completar decisões de formato e plano de testes negativos antes de modificar código.
- [ ] Reconciliar esta lista de artefatos com os arquivos efetivos; declarar arquivos adicionais antes de alterá-los.

### Implementation
- [ ] Implementar primeiro o menor fluxo que fecha a lacuna.
- [ ] Integrar entradas reais, persistência/compatibilidade e diagnósticos.
- [ ] Atualizar documentação e checks declarativos junto com o contrato.

### Validation
- [ ] Executar os cenários de cada R-ID, incluindo negativos.
- [ ] Executar pose assess integrate e validação estruturada no candidato.
- [ ] Reconciliar artifacts, surface e revisão independente antes do closeout.

## 5. Decisions

### Decision 1
- Date: 2026-09-07
- Context: ci.yml não instala POSE nem executa pose check/validate; skills-check é pulado sem binário. A matriz não inclui catálogo, checks de fixtures, formatação com falha nem race. delivery.json e artifacts.json estão desativados; roadmap não tinha Cut criteria; o bundle esbarra em cmd/ailearn/main.go histórico inexistente.
- Options considered: deixar implementação implícita no aceite; criar remediação focada.
- Decision: planejar remediação própria e exigir sua entrega antes do aceite.
- Rationale: v1-integrated-acceptance tem non-goal de adicionar features ou correções.
- Consequences: esta spec permanece draft; decisões estruturais novas exigem ADR no início da implementação.

## 6. Validation

### Strategy
Validar cada contrato com cenário positivo e negativo pela entrada de produção; usar somente dados sintéticos.

### Deterministic checks
- Test: pose check --strict; pose validate --strict; pose skills-check --strict; pose docs-check
- Lint: pose check --strict; pose lint-spec v1-delivery-ci-assurance --ready-check
- Typecheck: go vet ./...
- Build: go build ./cmd/codinho
- Security / Contract: pose assess integrate; pose validate --strict --json .pose/results/delivery-validation.json; pose surface-check --spec v1-delivery-ci-assurance --strict

### Execution log
- 2026-09-07 UTC: criada em revisão de planejamento; implementação e gates de entrega não executados.

### Results summary
Escopo proposto com requisitos verificáveis. A validação atual do produto está no relatório; não prova os novos requisitos.

### Requirement trace
Preencher R1–R8 com evidência por cenário durante a implementação e no closeout.

### Known gaps
Renames históricos e evidências de módulos divergentes bloqueiam roll-up. CI remota e revisão humana precisam de evidência real, não de YAML existente.

## 7. Final Report

### Delivered scope
Somente planejamento; nenhuma funcionalidade desta spec foi entregue.

### Files and modules changed
- Esta spec; dependências no roadmap e no aceite integrado.

### Validation executed
- Planejamento sujeito a pose lint-spec --ready-check e pose check --strict nesta auditoria.

### Residual risks
Validar o comportamento implementado em execução independente; não reutilizar resultado histórico como aprovação do código futuro.

### Follow-ups
- [covered: v1-integrated-acceptance] Reexecutar os requisitos desta spec no candidato composto da V1.
