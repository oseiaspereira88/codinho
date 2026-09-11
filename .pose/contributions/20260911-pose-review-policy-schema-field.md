# POSE check: campo de política rejeitado pelo schema v2

## Observação

Em 2026-09-11, `PATH=/home/go/go/bin:$PATH pose check --strict` falhou na
estrutura das políticas de review de 30 specs com a mensagem:

`invalid schema-v2 review policy: json: unknown field "evidence_vocabulary_reconciled_at"`

A falha ocorre antes de qualquer validação específica da alteração editorial
que motivou a execução. O working tree contém políticas geradas pelo POSE que
incluem esse campo, enquanto o binário local rejeita o mesmo campo como
desconhecido.

## Reprodução sintética

1. Em uma instância POSE com uma policy schema-v2 contendo apenas campos
   suportados, execute `pose check --strict`.
2. Acrescente o campo de metadado `evidence_vocabulary_reconciled_at` à mesma
   policy, sem alterar os demais campos.
3. Execute novamente `pose check --strict`.

Resultado observado: a segunda execução falha no parsing da policy com
`unknown field`, embora o campo esteja presente nas políticas geradas pela
instância.

## Impacto

O check estrutural não consegue produzir um resultado verde para a instância
com essas políticas. A falha não impede os testes determinísticos do executor
editorial, mas bloqueia a evidência de `pose check --strict` até que o schema
ou as políticas sejam alinhados.

## Próxima ação sugerida

Alinhar a versão do binário POSE e o schema de políticas: aceitar o campo como
metadado conhecido ou removê-lo/renomeá-lo na geração. Reexecutar `pose check
--strict` depois da correção.
