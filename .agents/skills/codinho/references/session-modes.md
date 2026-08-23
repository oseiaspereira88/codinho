# Modos pedagógicos e profundidade

## Modos (`session_start.mode`)

Os seis modos são independentes entre si — um especialista pode
escolher `teaching` em profundidade `micro` só para recordar sintaxe;
um iniciante pode pedir `exploration` num desafio inteiro.

| Modo | Quando o aluno pede | Defaults (`depth`/`help`/`disclosure_max`/`evaluation`) | Comportamento esperado da skill |
|---|---|---|---|
| `teaching` | Aprender algo novo do zero | `macro` / `progressive` / 2 (conceito ou API) / `on_demand` | Mais explicação proativa (`feedback_record` tipo `explanation`), pistas mais cedo na ladder. |
| `practice` (default) | Praticar o que já viu | `micro` / `progressive` / 1 (pergunta direcionadora) / `on_demand` | Fluxo padrão descrito em `mcp-tool-routing.md` — observar, avaliar, avançar. |
| `review` | Revisar algo já dominado, possivelmente citado por `review_due` | `meso` / `limited` / 1 / `on_step_complete` | Menos explicação, mais avaliação direta; boa hora para citar `mastery_evidence_record` com `dimension: retention`. |
| `debug` | "Meu código não funciona, me ajuda a achar o erro" | `macro` / `progressive` / 3 (estrutura lógica, nunca pseudocódigo) / `on_demand` | Nunca corrija diretamente (regra 3) nem revele a causa antes do aluno formular uma hipótese (PROJECT.md §8.2). Siga o protocolo de 6 estágios abaixo e use `check_run` real — o aluno encontra o bug, não você. Boa hora para `dimension: debugging` em `mastery_evidence_record`. |
| `exploration` | Explorar um desafio inteiro sem compromisso com conclusão | `layer` / `free` / 3 / `on_demand` | Sem obrigação de tentativa ou avanço (PROJECT.md §8.2) — feedback livre, menos pressão por `step_complete`. |
| `interview` | Simular uma entrevista técnica | `challenge` / `no_hints` / 0 / `only_at_end` | Comportamento específico de interview-mode (spec futura); estes são só os defaults de dimensão que a sessão usa até aquela spec configurar o protocolo completo (briefing, cronômetro, rubrica). |

Cada `session_start` pode sobrescrever qualquer campo explicitamente —
o modo só fornece o ponto de partida (PROJECT.md §8.3: as sete
dimensões nunca se inferem umas das outras).

### Protocolo de depuração (`mode: debug`)

PROJECT.md §8.1 declara seis estágios ordenados para o modo depuração.
Siga essa ordem ao guiar o aluno — nunca pule para a correção antes de
uma hipótese formulada pelo próprio aluno:

1. `reproduce` — reproduzir o problema relatado.
2. `locate_first_divergence` — localizar o primeiro ponto onde o
   comportamento observado diverge do esperado.
3. `hypothesize` — o aluno propõe uma hipótese (não você).
4. `observe` — coletar uma observação real (`check_run`,
   `workspace_observe`) que confirme ou refute a hipótese.
5. `confirm_or_reject` — confirmar ou rejeitar a hipótese com base na
   observação.
6. `apply_smallest_fix` — aplicar a menor correção possível — feita pelo
   aluno, nunca pela skill (regra 3).

Um desafio `kind: debug` traz seu código inicial (com o bug) declarado
em `fixture` (administrative-cli-fixtures); use `codinho workspace
prepare <challenge-id> --dest <path>` para materializá-lo antes de
iniciar a sessão — nunca escreva esse código você mesmo.

### Granularidade adaptativa e explicação obrigatória

PROJECT.md §8.5: o ajuste automático de granularidade só pode ocorrer
após evidência configurável (ex.: `progress_get` mostrando várias
evidências recentes sem ajuda, ou o oposto) e **deve ser informado ao
aluno**; ajuste manual do aluno sempre prevalece. Ao chamar
`granularity_adjust` — manual ou automaticamente — preencha sempre
`reason` com uma frase curta que você também repete ao aluno (ex.:
"você acertou os últimos 3 passos sem pista, vamos ampliar para
`meso`").

## Profundidades (`session_start.depth`, `granularity_adjust.depth`)

Do mais amplo ao mais granular: `challenge` → `layer` → `macro` →
`meso` → `micro`.

- **`micro`** (default): uma instrução pequena por vez, uma correção
  focal por interação (regra 11). Ideal para quem quer recordar sintaxe
  rápido.
- **`meso`/`macro`**: agrupam vários micropassos; menos interrupção,
  mais autonomia entre check-ins.
- **`layer`**: uma perspectiva inteira do desafio (ex.: só a camada de
  domínio, ignorando persistência).
- **`challenge`**: o desafio inteiro como uma unidade — para quem quer
  tentar de ponta a ponta antes de qualquer feedback.

Peça `granularity_adjust` quando o aluno disser algo como "pode ir mais
rápido" (aumentar profundidade) ou "quero passo a passo" (diminuir).
