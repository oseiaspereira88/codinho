# Matriz de compatibilidade

Ver [`docs/install.md`](install.md) para pré-requisitos de instalação e
[`.github/workflows/ci.yml`](../.github/workflows/ci.yml) para os gates
que rodam a cada pull request contra esta matriz (requirement R4/R5).

## Go

- **Mínimo declarado**: `go 1.25.0` (`go.mod`).
- **Testado nesta sessão**: `go1.26.6 linux/amd64` — `go build ./...`,
  `go vet ./...` e `go test ./... -race` passam limpos.
- `GOOS=windows go vet ./...` roda sem diagnósticos a cada spec fechada
  desde o início do projeto (checagem estática cross-compile; não é
  execução real em Windows).

## Sistema operacional

| OS | Status | Evidência |
|---|---|---|
| Linux | **Suportado, testado** | Toda a suíte de testes (`go test ./... -race`) roda nesta plataforma a cada spec. |
| macOS | **Suportado, não testado nesta sessão** | Código não usa nenhuma primitiva específica de Linux (paths, syscalls e process-group handling em `internal/checks/process_unix.go` usam a API POSIX padrão do pacote `syscall`, compatível com Darwin). Sem execução real registrada. |
| Windows | **Não declarado suportado até haver evidência real** | `GOOS=windows go vet` passa (checagem estática), e `internal/checks/process_windows.go`/`internal/diagnostics/process_windows.go` implementam os pontos de divergência de plataforma (grupo de processo, verificação de PID). Sem execução real registrada — não reivindicar suporte além do que a Compatibility desta spec permite ("documentar Windows como suportado somente após evidência"). |

## SDK MCP

- `github.com/modelcontextprotocol/go-sdk` — versão fixada em `go.mod`
  (`v1.7.0`). Atualizações passam pela suíte de contrato completa
  (`internal/mcpserver/*_contract_test.go`) antes de merge.

## Schema de catálogo e de evento

- Catálogo: `schema_version: 1` (`internal/curriculum.SchemaVersion`).
  Mudança de versão exige uma spec própria (não coberta por packs
  automaticamente).
- Evento: `schema_version: 1` (`internal/eventstore.CurrentEventSchemaVersion`).
  Compatibilidade futura é via `UpcasterRegistry` — um evento mais antigo
  é migrado evento a evento durante o replay; um evento mais novo sem
  upcaster registrado falha explicitamente
  (`ErrIncompatibleEventVersion`), nunca é silenciosamente ignorado.

## Atualização de pack sem quebrar sessão fixada

Reinicie `codinho serve` para aplicar edições do catálogo a sessões novas;
não há hot-reload. Sessões iniciadas com `recovery_version: 1` preservam
política completa e cópia do desafio com SHA-256 no evento de início.
O replay restaura estado, nó, pistas, detours, avaliação e revisão sem
reconsultar o desafio atual. Remover o desafio de um catálogo válido não
remove o conteúdo fixado de sessões existentes.

Repita uma mutação confirmada com o mesmo `request_id` e entrada original
para recuperar seu resultado, mesmo após reinício/progresso posterior.
Reutilização incompatível retorna `STATE_CONFLICT`. IDs antigos permanecem
reservados. Sessões legadas sem política/conteúdo recuperável retornam
`SESSION_RECOVERY_UNAVAILABLE`; preserve o histórico e inicie outra sessão.

Baselines, roots canônicas, globs e vínculo sessão–evidência sobrevivem ao
reinício. `evidence_get` verifica atualidade no escopo persistido; root
indisponível/substituída nunca demonstra freshness. Não troque root/globs
de uma baseline existente. `step_evaluate` também verifica a origem e
atualidade no consumo, conforme o contrato abaixo.

Use `codinho doctor` após interrupção abrupta. `SIGKILL` pode deixar lock
órfão: confirme que nenhum servidor usa o workspace antes de removê-lo,
conforme o diagnóstico. O writer com lock preserva bytes de cauda incompleta
em `events.jsonl.recovery-*` (permissão 0600) antes de truncar somente essa
cauda. Corrupção em linha terminada ou revisão inválida bloqueia startup;
inspeções da CLI são somente leitura e nunca reparam o log. Inclua backups
de recuperação na política de retenção/exclusão do estado.

Mantenha manifest e catálogo válidos para iniciar o servidor. O mecanismo
não migra payloads futuros nem corrige conteúdo corrompido por inferência.
Rollback para binário anterior perde projeção de sessão e proteção contra
colisões; não o execute como writer sobre estado novo sem estratégia de
migração. Consulte o [ADR](../.pose/adr/2026-09-06-durable-session-replay-with-pinned-content.md).

Evidência de 2026-09-07: `go test ./cmd/codinho -run TestSessionRecoveryOverRealStdio`
(update, remove, resposta perdida/cauda interrompida e legado/corrupção),
além das regressões em `internal/session`, `internal/application` e `internal/eventstore`.

## Evidência usada na avaliação

Use o evidence_id produzido por workspace_observe, check_run ou evidence_record
na mesma sessão e passo. Um blob existente em disco, um ID de outra sessão ou
um produtor de outro passo não autoriza avaliação. Para critério estrutural
baseado em check, envie também check_id, igual ao check executado do catálogo
fixado. Sem check comprovado, evidência de baseline/diff/nota resulta em
unverifiable; a presença de conteúdo não demonstra uma afirmação estrutural.

Registre explicações e observações externas com evidence_record antes de citá-las:
source deve ser learner_explanation, tutor_observation ou external_artifact;
text e rubric_ref são obrigatórios. O registro preserva origem e rubrica,
redige segredos reconhecidos e limita texto a 64 KiB. A avaliação qualitativa
cita o ID retornado e a mesma rubrica. O registro contém uma declaração do tutor;
não autentica uma origem externa, não acessa URLs e não substitui execução real.
Registros qualitativos podem ser lidos por evidence_get no escopo da sessão.

EVALUATION_EVIDENCE_INVALID indica blob ausente/corrompido, escopo ou check/rubrica
incompatível. EVALUATION_EVIDENCE_STALE indica que o root/fingerprint não pode
mais confirmar os arquivos observados; observe e execute novamente. Essas
rejeições não registram avaliação/tentativa nem alteram revisão. Check que
observe drift entre início/fim também rejeita seu resultado.

A avaliação amostra fingerprints no escopo dos globs antes da resolução e
imediatamente antes do append. Isso não congela escritores externos, cobre
arquivos fora dos globs ou detecta alterações que ocorram e sejam revertidas
entre amostragens. evidence_lineage no evento de avaliação identifica o produtor
e o fingerprint usados; a cobertura completa dos critérios autorados continua
sendo uma política distinta da autorização das evidências citadas.

Retries confirmados mantêm entrada/resultado originais, inclusive os anteriores
a check_id, sem reavaliar os arquivos atuais. Avaliações legadas permanecem
no histórico; evidência antiga sem escopo recuperável não autoriza avaliações
novas. O novo evento evidence_recorded exige este leitor; evite escritores
anteriores no estado atualizado. Consulte o
[ADR](../.pose/adr/2026-09-07-scoped-evaluation-evidence-and-qualitative-registration.md).

Valide com go test ./cmd/codinho -run TestEvaluationEvidenceOverRealStdio e
com os testes TestEvaluationEvidence de internal/application/internal/session.

## Navegação da árvore da sessão

Inicie com depth challenge, layer, macro, meso ou micro. session_start,
session_get e instruction_get retornam o ID e kind da janela real. Quando
não houver a profundidade pedida, use o nó mais profundo autorado. Challenge
exibe title/brief; layer exibe seu ID e o título do desafio, sem expor filhos.

Percorra todas as layers com step_complete e step_advance separados. Irmãos
são obrigatórios em ordem; children_mode: choice declara alternativas exclusivas.
Ao chegar a uma escolha numa profundidade mais fina, conclua o checkpoint
exibido e consulte step_advance. Envie next_step_id igual a um dos IDs diretos
retornados em branches. O servidor ativa a janela na profundidade pedida dentro
desse ramo. Um ID descendente, de outro ramo ou fora da árvore é rejeitado.

Ajuste a granularidade para exibir um ancestral ou retornar ao cursor mais fino.
Avaliações, pistas e filhos concluídos permanecem associados aos seus nós;
a troca de janela não avança para um irmão. Concluir uma janela ampla cobre
sua subárvore para navegação, sem produzir eventos de conclusão dos filhos.
Não refine uma janela já concluída nem reabra um ancestral cujos filhos foram
percorridos para obter crédito novamente: avance.
Repetir a conclusão do mesmo nó não cria outro evento. Retornar done em
step_advance não encerra a sessão; use session_finish explicitamente. A primeira
confirmação de done é registrada para preservar a auditoria de overrides.

Use a revisão atual também ao consultar branches ou done. Overrides permitem
pular a janela atual, ficam registrados e não criam evidência de conclusão.
Navegação exige sessão ativa e nenhum desvio conceitual aberto. Hints continuam
sujeitos à política da sessão, inclusive em entrevista e após reinício.

Preserve logs antigos: inícios sem navigation_version mantêm o macro original,
e retries recuperam os resultados da época. Novos eventos usam
navigation_version: 1; versões desconhecidas falham com diagnóstico de
recuperação. Não use escritores antigos num estado que já recebeu esses
novos eventos: eles não preservam cursor, escolhas e progresso por nó.

Valide com go test ./cmd/codinho -run TestTreeProgressionOverRealStdio.
Consulte o [ADR do cursor](../.pose/adr/2026-09-07-depth-aware-session-cursor-and-explicit-branch-selection.md)
para limites de agregação e integração futura com trilhas.

## Limites de recursos aplicados

| Recurso | Limite | Onde |
|---|---|---|
| Tamanho de arquivo YAML de pack | 1 MiB | `curriculum.DefaultLimits.MaxFileBytes` |
| Profundidade de nó YAML | 32 | `curriculum.DefaultLimits.MaxNodeDepth` |
| Aliases YAML | 64 | `curriculum.DefaultLimits.MaxAliases` |
| Saída de check (stdout/stderr) | 256 KiB, com sufixo de truncamento | `internal/checks/limits.go` |
| Checks concorrentes | 4 | `internal/checks/limits.go` (`maxParallel`) |
| Timeout de check | por check, validado por `checks.Resolve` | `internal/checks/registry.go` |
| Resultado de busca no catálogo | 100 itens, com `Truncated: true` explícito | `internal/curriculum/search.go` |

## Desempenho (informativo, não um SLO de serviço cloud)

- Busca no catálogo em escala V1 (84 desafios sintéticos): p95 bem
  abaixo de 100ms — `TestSearchP95BudgetAtV1Scale`.
- Startup do `codinho serve` até a primeira resposta MCP real: bem
  abaixo de 1s — `TestServeStartupRespondsWithinOneSecond`.

Esses números são medidos no hardware desta sessão de desenvolvimento;
são thresholds informativos, não uma garantia de SLO (Non-functional:
"Benchmarks são informativos com thresholds documentados").
