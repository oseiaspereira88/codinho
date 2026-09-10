# Revisão do lote de aprofundamento do catálogo

Escopo inspecionado: commits `0a8c964^..13e9dc1`, com foco nos novos nodes de
go-core, go-io e go-testing e nas fixtures que sustentam seus critérios.
Spec: `go-foundations-packs`, ainda in-progress.

## Achados e correções

- Média: três critérios não tinham cobertura discriminante. Uma implementação
  que perde bytes retornados junto de EOF, outra que retorna duração positiva
  com ErrMissingClock e outra que retorna slice alocado para entrada vazia
  passavam nos testes anteriores. Foram adicionadas fixtures de regressão ao
  exercício e à referência e ampliados os filtros dos checks para executá-las.
- Média: “slice vazio comparável” em testable-expiry induzia uma propriedade
  incorreta de Go. A instrução agora exige vazio não nil com relógio válido.
- Baixa: check-ticket-boundary descrevia entrada em lote apesar de Valid receber
  dois instantes. O escopo agora pede tabela de chamadas individuais.
- Baixa: o node de newline apontava para erros de Writer apesar de existir
  io-json-encode. Referência corrigida. A propagação de erro agora pede errors.Is,
  como a fixture, em vez de exigir igualdade de instâncias além do contrato.
- Baixa: a frase “nenhuma entrada mantém o retorno nil” negava o resultado
  esperado no desafio de defer. Corrigida para entradas nil e vazias.

## Evidência e limites

Os três mutantes passaram nos testes antigos e falharam nas novas regressões.
`catalog validate --checks --json` retornou diagnostics null após as correções.
As fixtures verificam comportamento da solução; não demonstram que o aluno
escreveu bons testes nem substituem playtest humano.

O catálogo também emite avisos editoriais `missing_hints_for_gated_step` nos
nodes recentes. Portanto, diagnostics null não deve ser reportado como ausência
de toda pendência editorial. A próxima etapa é escrever pistas específicas para
esses passos e revisar redundâncias, antes de acrescentar nodes à contagem.

`go test -race -count=1 ./...` passou. A primeira matriz encontrou bloqueio de
escrita no cache de módulos durante build e foi repetida com permissão de acesso.
O executável govulncheck instalado em `/home/go/go/bin` também precisou ser
incluído no PATH para a matriz localizar o scanner obrigatório.
Resultado final: `PATH="/home/go/go/bin:$PATH" pose validate --strict` passou
com `Result: SUCCESS`; govulncheck não encontrou vulnerabilidades.
`pose check --strict` e `git diff --check` passaram. Artifact-check identificou
avisos globais preexistentes de proveniência em dois arquivos de curriculum,
fora deste diff; não são evidência de fechamento da spec.

O bundle `rvb-e43e5212ea41ad2d` (digest
`sha256:e43e5212ea41ad2d9f352f5236b4c346bbf223b54999d133bb2e545f3370c1f3`)
resolveu a spec inteira, abrangendo mudanças anteriores a esta revisão, e
informou ausência de evidência estruturada atribuída. Não foi selado nem
atestado: esta revisão de lote não aprova o fechamento da spec inteira.

## Rules applied during review

- Tipo: correção de conteúdo e testes de regressão.
- Workflows: review e bugfix; skills pose-review e pose-bugfix.
- Security: exemplos sintéticos, sem novas dependências ou dados privados.
- Documentation-style e delivery-evidence: correções verificáveis, sem declarar
  nodes adicionais como prova de entrega pedagógica.
- Knowledge consumido: knowledge:go-foundations-io-batch, especialmente a
  restrição contra inflar a árvore com passos redundantes.
- Tech-debt: zero marcadores descobertos; recurrence-check de 14 dias: zero
  chaves acima do limiar. Integração MCP não se aplica a este diff de catálogo.
- Ferramentas de fechamento permanecem diferidas; revisão humana e metas da
  spec continuam pendentes. Não há mudança de assinatura pública.

Decisão sobre o lote original: mudanças solicitadas, corrigidas neste lote.
Continuação implementada: cobertura executável das três fronteiras acima.
