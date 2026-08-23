# Rubricas de feedback

`feedback_prepare` cita estas rubricas por referência (`rubric_refs`);
o servidor nunca lê o conteúdo delas — só o SKILL.md/você as usa para
redigir feedback e julgamento qualitativo em `step_evaluate`. Todo
julgamento qualitativo exige `evidence_id` (o que você observou) e
`rubric_ref` (qual destas seções fundamenta o julgamento) — nunca
julgue sem os dois.

Verdicts possíveis para um critério qualitativo: `met`,
`partially_met`, `not_met`, `unverifiable`, `not_applicable`. Severidade
do achado: `blocking`, `important_non_blocking`, `advisory`. Estilo e
preferência quase nunca são `blocking` — reserve isso para violações
reais do enunciado do desafio.

## rubric://idiomatic-go

Critérios para avaliar se o código é idiomático em Go, **nunca** se
está "certo" no sentido funcional (isso é `structural`, decidido pelo
servidor a partir de evidência real de check, nunca por você aqui):

- **Nomes**: identificadores curtos e claros em escopo pequeno
  (`i`, `err`, `buf`); nomes descritivos cruzando limites de pacote;
  receiver names curtos e consistentes por tipo; sem prefixo/sufixo
  redundante com o pacote (`user.UserID` é ruim, `user.ID` é bom).
- **Erros**: erros retornados, nunca engolidos silenciosamente;
  `errors.Is`/`errors.As` em vez de comparação de string; mensagens de
  erro em minúsculas, sem pontuação final, compondo bem quando
  encadeadas (`fmt.Errorf("doing x: %w", err)`).
- **Valores zero**: tipos cujo valor zero já é útil sem construtor
  obrigatório, quando fizer sentido; `nil` slice vs. slice vazio
  tratados com a mesma semântica de leitura.
- **Interfaces**: pequenas, definidas do lado do consumidor quando
  possível; aceitar interface, retornar struct concreto.
- **Concorrência**: canais e goroutines com dono claro de quem fecha o
  quê; sem *data race* (isto é verificável por `check_run` com
  `go_test_race` — cite essa evidência, não “pareceu certo”).
- **Ownership/lifecycle**: struct que gerencia recurso expõe
  `Close`/liberação clara; nenhuma mutação surpresa de argumento que o
  chamador não esperava.

Achados de estilo puro (ex.: preferência entre duas formas igualmente
idiomáticas) são `advisory`. Um erro engolido silenciosamente ou uma
data race real são `blocking` — mas cite a evidência de `check_run`
para a race, nunca alegue sem ela.

## rubric://technical-communication

Critérios para avaliar a qualidade de uma explicação ou reflexão do
aluno (`reflection_record`, respostas a `concept_content_get`):

- **Clareza**: a explicação é compreensível sem exigir que o leitor já
  saiba a resposta.
- **Precisão**: termos técnicos usados corretamente; nenhuma afirmação
  factualmente errada sobre a linguagem ou a biblioteca padrão.
- **Reconhecimento de trade-offs**: quando existe mais de uma
  abordagem razoável, o aluno reconhece isso em vez de apresentar a
  escolha como única possível.
- **Honestidade sobre incerteza**: "não tenho certeza, mas acho que..."
  é uma resposta de melhor qualidade comunicativa do que uma afirmação
  confiante e errada.

Uma reflexão vaga ("porque sim", "funciona e pronto") é
`not_met`/`partially_met` conforme o contexto, nunca `blocking` — isto
nunca é `structural` e nunca impede `step_complete` por si só, a menos
que a política do passo exija reflexão explicitamente.
