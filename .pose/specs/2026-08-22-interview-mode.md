---
slug: interview-mode
status: done
created_at: 2026-08-22
completed_at: 2026-08-23
supersedes:
depends_on: feedback-evaluation-progression, safe-check-executor, tutor-skill-host-integration, learning-practice-debug-modes
priority: 160
components: interview-mode, tutor-skill
delivers:
---

# Spec: interview-mode

## 1. Intent

### Goal
Implementar simulação técnica genérica com briefing completo, tempo opcional, auxílio controlado e revisão somente ao final.

### Business value
Praticar autonomia, comunicação e solução sob condições próximas de live coding sem alegar vínculo com terceiros.

### Constraints
- Não reproduzir perguntas, marcas ou processos proprietários.
- Não prometer antifraude em máquina local.
- Pistas ficam bloqueadas ou registradas conforme política escolhida antes do início.

### Non-goals
- Proctoring, gravação de tela, vigilância ou ranking.
- Avaliação de contratação.

## 2. Requirements

### Functional
- R1: Iniciar simulação com desafio, duração, checks permitidos e política de auxílio fixada.
- R2: Exibir apenas briefing e critérios públicos durante execução.
- R3: Registrar tempo ativo, pausas autorizadas e término sem manipular relógio do sistema.
- R4: Bloquear ou registrar pedidos de pista sem revelar conteúdo proibido.
- R5: Permitir finalização antecipada ou por timeout com estado explícito.
- R6: Avaliar ao final código, testes, raciocínio, complexidade, comunicação e trade-offs.
- R7: Produzir relatório descritivo com evidências, gaps e exercícios recomendados.
- R8: Declarar limites de integridade e ausência de associação com terceiros.

### Non-functional
- Clock deve ser injetável e determinístico em testes.
- Relatório não deve reduzir desempenho a um score único.

### Security
- Não coletar áudio, vídeo, dados pessoais ou atividade fora do workspace.

### Compatibility
- Simulação deve funcionar sem recursos MCP opcionais.

## 3. Technical Plan

### Affected areas
- internal/session/, internal/assessment/, .agents/skills/codinho/, packs/go-interviews/

### Artifacts
- created: internal/session/interview.go
- created: internal/session/interview_test.go
- created: internal/assessment/interview_report.go
- created: internal/assessment/interview_report_test.go
- created: internal/application/interview.go
- modified: internal/session/service.go
- modified: internal/mcpserver/session_tools.go
- modified: internal/mcpserver/assistance_tools.go
- modified: internal/mcpserver/contract_test.go
- modified: internal/mcpserver/session_contract_test.go
- modified: .agents/skills/codinho/references/session-modes.md
- modified: .agents/skills/codinho/references/mcp-tool-routing.md
- modified: .agents/skills/codinho/SKILL.md
- created: .agents/skills/codinho/assets/interview-report-template.md
- created: testdata/host/interview-transcripts/full-interview-session.md

### Delivery targets
Nenhum novo; amplia a capability de tutoria planejada.

### API/contract changes
- Adicionar propriedades de timer e policy lock à sessão.

### Data/storage changes
- Persistir interview_started, interview_finished e report_generated.

### Technical risks
- Timer pode criar pressão sem valor pedagógico.
- Agente pode oferecer ajuda por conversa fora das tools.

## 4. Tasks

### Planning
- [x] Definir rubrica e linguagem neutra do relatório.
- [x] Definir políticas de pausa, timeout e auxílio.

### Implementation
- [x] Implementar timer e lock de policy.
- [x] Integrar bloqueio de disclosure.
- [x] Implementar avaliação e relatório final.
- [x] Atualizar skill e templates.
- [x] Criar transcripts de abuso e encerramento.

### Validation
- [x] Testar clock, timeout, pausa e restart.
- [x] Testar tentativas de obter pistas por linguagem indireta.
- [x] Revisar privacidade e neutralidade do relatório.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Simulação local não consegue garantir integridade forte.
- Options considered: fingir antifraude; proctoring; transparência e honor system.
- Decision: declarar limites e focar prática, não certificação.
- Rationale: evita vigilância e promessas falsas.
- Consequences: resultados são evidência pedagógica, não credencial.

### Decision 2
- Date: 2026-08-23
- Context: R3 exige clock injetável e determinístico; nenhum lugar do
  código deveria simplesmente chamar `time.Now()` dentro da lógica de
  domínio de timing.
- Options considered: (a) adicionar um clock injetável ao próprio
  `internal/eventstore` (mudança invasiva num pacote já sealed/testado
  de local-event-store); (b) seguir o padrão já estabelecido por
  `ProgressService.ReviewDue(now time.Time)`: a camada de tool MCP passa
  `time.Now()` real, e toda lógica de negócio recebe `now` como
  parâmetro explícito.
- Decision: (b). `InterviewElapsed`, `Service.InterviewStatus` e
  `Service.InterviewReport` recebem `now` explicitamente; o horário de
  início real vem do `RecordedAt` já durável do evento `session_started`
  (nunca inventado).
- Rationale: reaproveita um padrão já testado no código-base em vez de
  invadir um pacote sealed; testes chamam essas funções com um `now`
  fixo sem precisar de nenhum clock mockável adicional.
- Consequences: qualquer nova lógica de timing neste domínio deve seguir
  o mesmo padrão (parâmetro explícito, nunca `time.Now()` interno).

### Decision 3
- Date: 2026-08-23
- Context: R4 pede registrar pedidos de pista bloqueados (PROJECT.md
  §9.5), mas `HintRequest` (assistance-hints-detours, já testado)
  simplesmente retorna `ErrHelpDisabled` sem persistir nada quando a
  política bloqueia.
- Options considered: (a) modificar `HintRequest` para persistir um
  evento mesmo no caminho de erro; (b) adicionar
  `Service.RecordBlockedHintAttempt`, chamado pela camada de tool MCP
  somente quando `HintRequest` retorna especificamente
  `ErrHelpDisabled`, sem tocar `HintRequest` em si.
- Decision: (b).
- Rationale: isola a mudança do caminho já testado de
  assistance-hints-detours; a tool MCP é o lugar certo para essa decisão
  porque só ela sabe distinguir "bloqueado" de qualquer outro erro.
- Consequences: o evento de bloqueio usa a revisão atual do stream como
  `expectedRevision` (sempre sucede, pois é um fato de auditoria interno
  que o chamador não tinha como prever) e portanto avança a revisão real
  do stream — um chamador com uma revisão antiga recebe o mesmo
  `STATE_CONFLICT` retryable que qualquer escrita concorrente já produz
  (não um comportamento novo, apenas o protocolo de concorrência
  otimista já existente se aplicando a mais um tipo de evento).

## 6. Validation

### Strategy
Combinar clocks falsos, policy tests, transcripts adversariais e revisão de relatório.

### Deterministic checks
- Test: go test ./internal/session/... ./internal/assessment/...
- Lint: gofmt -l internal/session internal/assessment
- Typecheck: go vet ./internal/session/... ./internal/assessment/...
- Build: go build ./cmd/codinho
- Security / Contract: privacy review e leak tests.

### Execution log
- `go build ./...` → ok (2026-08-23).
- `gofmt -l .` → sem saída (2026-08-23).
- `go vet ./...` e `GOOS=windows go vet ./...` → sem diagnósticos (2026-08-23).
- `go test ./... -race` → `ok` em todos os pacotes, incluindo `internal/session` (interview.go), `internal/assessment` (interview_report.go) e `internal/mcpserver` (interview_status, interview_report, session_finish com reason, hint bloqueado registrado), sem data races (2026-08-23).
- `govulncheck ./...` → "No vulnerabilities found." (2026-08-23).

### Results summary
Simulação de entrevista genérica entregue: `session_start` aceita
`time_limit_seconds`; `interview_status` reporta tempo decorrido e
timeout a partir do `RecordedAt` real do evento `session_started`
(nunca um clock manipulável); pedidos de pista bloqueados pela política
são registrados sem revelar nada (`RecordBlockedHintAttempt`);
`session_finish` aceita `reason` (`explicit`/`timeout`);
`interview_report` monta um relatório descritivo (avaliações, pistas
concedidas/bloqueadas, reflexões, gaps, recomendações) que nunca reduz
o desempenho a um score único, sempre com `integrity_note` declarando
os limites da simulação (R8).

### Requirement trace
- R1 [satisfied] session.StartInput.TimeLimit + session_start (mode+time_limit_seconds+help) + test:TestStartAcceptsTimeLimit.
- R2 [satisfied] instruction_get já restringe a objective+scope; disclosure_max: 0 do modo interview bloqueia decomposição.
- R3 [satisfied] internal/session/interview.go (InterviewElapsed, InterviewStatus usando RecordedAt real) + test:TestInterviewStatusUsesRealRecordedStartTime.
- R4 [satisfied] Service.RecordBlockedHintAttempt + hook em assistance_tools.go + test:TestContractHintRequestBlockedByPolicyIsRecorded.
- R5 [satisfied] Service.FinishWithReason + tool session_finish(reason) + test:TestFinishWithReasonPersistsReason, test:TestContractSessionFinishAcceptsReasonAndInterviewReportReadsItBack.
- R6 [satisfied] assessment.BuildInterviewReport (Evaluations, Hints, Reflections) + test:TestBuildInterviewReportNeverCollapsesToASingleScore.
- R7 [satisfied] InterviewReport.Gaps (lista, nunca score) + Recommended (caller-supplied) + tool interview_report.
- R8 [satisfied] assessment.integrityNote presente em todo InterviewReport + test:TestBuildInterviewReportNeverCollapsesToASingleScore (assert IntegrityNote != "").

### Known gaps
- Validação humana de realismo e utilidade será parte do aceite (v1-integrated-acceptance).
- Conteúdo curado de entrevistas (`packs/go-interviews/`) pertence a go-interviews-pack.
- `pause`/`resume` durante entrevista reutilizam o lifecycle genérico já existente; nenhum comportamento especial de "pausa autorizada" foi adicionado além do que Pause/Resume já fazem — se um piloto revelar necessidade de distinguir pausa autorizada de pausa comum, isso é um follow-up.

## 7. Final Report

### Delivered scope
Simulação de entrevista técnica genérica: duração opcional, política de
auxílio fixada, bloqueio de pistas registrado, encerramento explícito ou
por timeout, e relatório final descritivo com aviso de integridade.

### Files and modules changed
- internal/session/interview.go (novo): InterviewElapsed, InterviewStatus, RecordBlockedHintAttempt, FinishWithReason, InterviewReport.
- internal/session/service.go: StartInput.TimeLimit; Finish delega a FinishWithReason.
- internal/assessment/interview_report.go (novo): InterviewReport pura, nunca reduz a um score.
- internal/application/interview.go (novo): aliases e wrappers para mcpserver.
- internal/mcpserver/session_tools.go: session_start(time_limit_seconds), session_finish(reason) dedicado, tools interview_status e interview_report.
- internal/mcpserver/assistance_tools.go: hint_request registra tentativa bloqueada.
- .agents/skills/codinho/references/{session-modes,mcp-tool-routing}.md, SKILL.md, assets/interview-report-template.md, testdata/host/interview-transcripts/.

### Validation executed
- Command: go test ./... -race
- Result: ok em todos os pacotes.
- Command: govulncheck ./...
- Result: No vulnerabilities found.

### Residual risks
- Variações de host podem afetar bloqueio conversacional (mitigado pela regra 14 e pelo bloqueio real no servidor, não apenas na skill).
- Sem piloto humano ainda para validar se 45min/no_hints é um default razoável — fácil de ajustar via `time_limit_seconds`/`help` explícitos.

### Follow-ups
- [covered: go-interviews-pack] Fornecer quatro simulações curadas.
