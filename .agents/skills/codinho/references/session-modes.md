# Modos pedagógicos e profundidade

## Modos (`session_start.mode`)

Os seis modos são independentes entre si — um especialista pode
escolher `teaching` em profundidade `micro` só para recordar sintaxe;
um iniciante pode pedir `exploration` num desafio inteiro.

| Modo | Quando o aluno pede | Comportamento esperado da skill |
|---|---|---|
| `teaching` | Aprender algo novo do zero | Mais explicação proativa (`feedback_record` tipo `explanation`), pistas mais cedo na ladder. |
| `practice` (default) | Praticar o que já viu | Fluxo padrão descrito em `mcp-tool-routing.md` — observar, avaliar, avançar. |
| `review` | Revisar algo já dominado, possivelmente citado por `review_due` | Menos explicação, mais avaliação direta; boa hora para citar `mastery_evidence_record` com `dimension: retention`. |
| `debug` | "Meu código não funciona, me ajuda a achar o erro" | Nunca corrija diretamente (regra 3). Guie por perguntas e `check_run` real — o aluno encontra o bug, não você. Boa hora para `dimension: debugging` em `mastery_evidence_record`. |
| `exploration` | Explorar um desafio inteiro sem compromisso com conclusão | Granularidade mais alta (`challenge`/`layer`), menos pressão por `step_complete`. |
| `interview` | Simular uma entrevista técnica | Comportamento específico de interview-mode (spec futura); até lá, trate como `practice` com feedback mais escasso e avaliação concentrada no fim. |

Mecânica detalhada de `debug` e `interview` (perguntas socráticas
estruturadas, cronômetro, rubrica de entrevista) pertence às specs
`learning-practice-debug-modes` e `interview-mode` — quando essas specs
forem implementadas, esta referência será atualizada com os detalhes
mecânicos específicos. Até lá, os seis valores de `mode` já são aceitos
por `session_start` e a sessão se comporta com as políticas genéricas
(ajuda, avaliação, divulgação) descritas acima.

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
