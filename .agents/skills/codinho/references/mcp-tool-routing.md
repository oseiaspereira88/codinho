# Roteamento de tools MCP (contrato V1, 29 tools)

Toda resposta é um Envelope com `status` (`ok`/`error`), `progress_effect`,
e (em erro) `error.code`/`error.message`/`error.retryable`. Toda tool que
muda estado de sessão exige `expected_revision` (pegue de `session_get`
antes) e aceita `request_id` opcional para retry idempotente — sempre
reuse o mesmo `request_id` ao repetir uma chamada que pode ter falhado
por timeout, nunca ao repetir uma ação genuinamente nova.

## 1. Descoberta de catálogo (somente leitura, sem sessão)

| Tool | Quando usar |
|---|---|
| `catalog_search` | Filtrar por `kind`, `theme`, `text`, `competency`, `difficulty`, `challenge_kind`, `max_minutes`, `prerequisite`. Ajuda o aluno a escolher um desafio. |
| `catalog_get` | Detalhar um item por `id` antes de `session_start`. |
| `concept_relations_get` | Ver `requires`/`recommended_before`/`relates_to`/`contrasts_with`/`commonly_fails_with`/`applies_in`/`deepens_into`/`evidences` de um item — útil para explicar por que um desafio é sugerido. |
| `learning_path_recommend` | Recomendar o próximo desafio a partir de `competency_id`/`theme_id`/`time_budget_minutes`/`completed` (lista de challenge_id que VOCÊ sabe que o aluno já terminou — o servidor não guarda isso). Cada recomendação vem com `explanation` (evidence/gap/dependency/cost_minutes) — sempre repasse essa explicação ao aluno, nunca só o `challenge_id`. **Consultivo apenas: nunca chame `session_start` sozinho a partir de uma recomendação sem o aluno aceitar.** |

## 2. Início e leitura de sessão

| Tool | Quando usar |
|---|---|
| `session_start` | Uma vez por desafio escolhido. Guarde `session_id` e `revision` retornados. |
| `session_get` | Antes de qualquer inferência de estado (regra 5); sempre que você não tiver certeza de `revision`/passo ativo/política. |
| `instruction_get` | Para reler a instrução ativa (objetivo+escopo) sem mudar nada. |
| `session_configure` | Aluno pede para mudar política de ajuda ou avaliação no meio da sessão. |
| `session_pause` / `session_resume` / `session_finish` | Aluno pausa, retoma ou encerra explicitamente. Nunca infira encerramento por silêncio. |
| `granularity_adjust` | Aluno pede mais/menos granularidade (`challenge`→`micro`). |

## 3. Assistência progressiva

| Tool | Quando usar |
|---|---|
| `hint_request` | Aluno pede ajuda. Sobe UM nível por chamada. `confirm_solution: true` só quando o próximo nível é a solução (nível 6) E o aluno confirmou que quer isso. |
| `syntax_recall_get` | Aluno esqueceu sintaxe específica (não é uma pista pedagógica sobre a solução — é lembrete de sintaxe, livre mesmo em políticas restritivas quando aplicável). |
| `concept_content_get` | Aluno pede para entender um conceito citado na instrução, por `concept_id`. Use exemplos diferentes da solução ativa (regra 13). |
| `learning_detour_start` / `learning_detour_finish` | Aluno faz uma pergunta conceitual fora do fluxo do passo atual. Abra o desvio, responda, feche com `outcome: resolved` ou `abandoned`, e retome o passo original — o desvio nunca substitui o passo ativo. |

## 4. Observação e verificação de código

| Tool | Quando usar |
|---|---|
| `workspace_observe` | SEMPRE antes de `step_evaluate` com evidência real (regra 6). Primeira chamada por `(session, step)` estabelece baseline; chamadas seguintes reportam diff. Precisa de `root` (caminho real autorizado) e, opcionalmente, `globs` (arquivos do desafio). |
| `check_run` | Depois de `workspace_observe` no mesmo passo. Só executa `check_id` já declarado no desafio (nunca um comando livre). Reusa o `root`/`globs` da última `workspace_observe` — não repita esses parâmetros. Retorna `outcome: pass|fail|error|skipped` e `evidence_id`. |
| `evidence_get` | Para reler o conteúdo de uma evidência já produzida (por `workspace_observe` ou `check_run`), checando `stale` (o workspace mudou desde a coleta). |

## 5. Feedback e avaliação

| Tool | Quando usar |
|---|---|
| `feedback_prepare` | Antes de redigir feedback: monta objetivo, escopo, pergunta e `rubric_refs` (ver `feedback-rubric.md`). Nunca gera o texto por você. |
| `feedback_record` | Registra o feedback que VOCÊ redigiu, com `type` (`confirmation`, `question`, `explanation`, `relation`, `suggestion`, `idiom`, `risk`, `violation`, `error`). Nunca conclui nem avança (regra 8). |
| `step_evaluate` | Com `criteria` (cada um com `name`, `kind`, `severity`, e — para critérios não-`structural` — `evidence_id`+`rubric_ref`+seu `verdict` julgado). `submission_intent: true` só com confirmação inequívoca do aluno (regra 7). Um critério `structural` citando evidência de `check_run` recebe verdict REAL (pass→met, fail→not_met); citando evidência sem check real, o verdict é `unverifiable` a menos que haja presença de evidência (`met`) — nunca simule um resultado. |
| `reflection_record` | Aluno responde a uma pergunta de reflexão. Cite `competency_id` quando souber qual. |
| `step_complete` | Só depois de avaliação satisfazer a política do passo, ou com `override: true` explícito (nunca proponha override sem o aluno pedir). |
| `step_advance` | Chamada separada de `step_complete` (regra 8, invariante 5). Pode retornar `branches` (ramificação — apresente as opções) ou `done: true` (fim da árvore). |

## 6. Domínio e progresso (não são por sessão)

| Tool | Quando usar |
|---|---|
| `progress_get` | Ver domínio atual por competência/dimensão. Sem `session_id` — é global ao aluno local. Use `revision` retornado como `expected_revision` da próxima `mastery_evidence_record`. |
| `review_due` | Ver revisões espaçadas vencidas, já ordenadas por prioridade. Nunca bloqueia início de sessão livre (regra R6 de mastery-review-scheduling) — é só uma sugestão. |
| `mastery_evidence_record` | Depois de uma avaliação real (não invente): cite `competency_id`, `dimension` (uma das 8: `understanding`, `syntax_recall`, `guided_implementation`, `autonomous_implementation`, `debugging`, `explanation`, `retention`, `transfer`), `evidence_id` já existente, `help_used`, `solution_revealed`, `success`. O servidor decide a promoção — você nunca define o estado diretamente. |

## Ordem típica de uma sessão de prática

```
catalog_search / learning_path_recommend  (aluno escolhe)
session_start
instruction_get
[hint_request]*                    (conforme pedido)
workspace_observe                  (baseline)
  ... aluno escreve código ...
workspace_observe                  (diff)
check_run                          (se o desafio declarar checks)
feedback_prepare → feedback_record (opcional, consultivo)
step_evaluate                      (com submission_intent quando o aluno confirmar)
reflection_record                  (se o passo pedir reflexão)
step_complete
step_advance
mastery_evidence_record            (se a evidência justificar)
```
