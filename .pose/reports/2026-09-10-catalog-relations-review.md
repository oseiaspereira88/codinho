# Relações de aplicação do catálogo

Spec: go-foundations-packs, in-progress.

Os 32 avisos relation_isolated correspondiam a competências sem arestas
explícitas e aos desafios/conceito de filtragem e depuração de slices.
Foram acrescentadas 39 relações applies_in em cinco packs: debugging (5),
errors (6), first-steps (3), io (16) e testing (9).

Cada competência aponta aos desafios que já a declaram como primária ou
secundária. O conceito slice-bounds aponta ao desafio slice-off-by-one que
o utiliza. A relação descreve contexto de aplicação, não comprovação de
domínio pelo aluno. Não foram criados requisitos de precedência ou ciclos
de dependências, nem alteradas fixtures, autoria ou publicação.

A comparação das árvores YAML confirmou que apenas relations mudou.
Catalog validate retornou diagnostics null e editorial null: os 32 avisos
anteriores foram eliminados. Isso não substitui revisão humana ou playtest.
`PATH="/home/go/go/bin:$PATH" pose validate --strict` passou com Result: SUCCESS,
incluindo checks das fixtures, testes com race detector e scanner de
vulnerabilidades. `pose check --strict` e `git diff --check` passaram.

Conhecimento consumido: knowledge:go-foundations-io-batch. Regras aplicadas:
documentation-style (relações explícitas), security (sem dados sensíveis) e
delivery-evidence (separação entre validação e conclusão pedagógica).

O catálogo permanece com 71 conceitos, 60 competências, 54 desafios e 213
nodes totais. As metas e o fechamento de go-foundations-packs seguem pendentes.
