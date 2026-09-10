# Formas de chamada variádica

Spec: go-foundations-packs, in-progress.

O desafio go-core.sum-variadic-numbers testava apenas Sum(c.nums...), apesar
de o aceite mostrar chamadas diretas. Foi adicionado um passo para comparar
argumentos diretos com expansão de slice. A reflexão diferencia o local das
reticências na declaração e na chamada; os hints declaram pergunta orientadora
e recordação de sintaxe.

call_forms_test.go foi incluído nas fixtures do exercício e da referência.
Os casos chamam Sum(), Sum(5), Sum(1, 2, 3) e Sum(-1, 5, -4), além das versões
com expansão. Cada resultado é comparado ao valor esperado, evitando aceitar
duas implementações igualmente incorretas por mera igualdade entre resultados.
O filtro do check agora inclui TestSumCallForms.

Conhecimento consumido: knowledge:go-foundations-io-batch. A expansão mantém
assinatura, aceite e publicação; não declara revisão humana ou playtest.
Assessment de packs executado antes da edição. A spec permanece aberta.

Validação: catalog validate --checks retornou diagnostics null e editorial
null; `PATH="/home/go/go/bin:$PATH" pose validate --strict` terminou com
Result: SUCCESS. `pose check --strict` e `git diff --check` passaram.
