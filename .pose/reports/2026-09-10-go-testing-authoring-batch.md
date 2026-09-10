# Lote go-testing — 2026-09-10 UTC

Spec: go-foundations-packs (in-progress). Um desafio combinado draft,
oito conceitos, cinco competências e sete nodes. Autoria codex; não há
reviewed_by ou playtested humano declarado.

A suíte testa fronteiras temporais, relógio ausente, chamada única inclusive
com entrada vazia, mesmo instante em fusos diferentes, ordem preservada e
independência da saída. Os casos são determinísticos e não usam esperas reais.
Os micropassos intermediários usam evidência qualitativa com rubrica; o check
completo fica no último passo. O novo caso solicitado ao aluno precisa de
revisão qualitativa, pois a suíte não comprova sua autoria ou qualidade.

Validação: pose validate --strict executou 22 checks, todos aprovados, sem
skips. O gate editorial verificou 41/41 checks do catálogo, com test_failure
na baseline e pass na referência de go-testing. Docs-check: zero erros/warnings.
Tech-debt: zero marcadores descobertos. Sem nova dependência, rede ou arquivos
externos nas fixtures. Não houve pré-revisão pedagógica independente neste lote.

O catálogo fundamental contém agora 43 desafios (33 A / 9 C / 1 F),
62 conceitos, 45 competências e 104 nodes. Restam duas autorias em go-errors,
variantes, aprofundamento, trilhas e revisão/playtest. A divergência da meta
44 versus 45 continua documentada; nenhum conteúdo foi removido/reclassificado.
