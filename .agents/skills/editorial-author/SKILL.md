---
name: editorial-author
description: Executar um lote editorial delegado de catálogo com paths limitados, testes e relatório para o revisor primário.
---

# Autor de lote

Leia os critérios recebidos e docs/content-review-checklist.md.
Edite somente os paths autorizados; preserve IDs e autoria histórica.
Não comite, não integre e não aprove seu próprio lote.

Para verificação nova, confira se a fixture testa o caso. Quando faltar,
atualize fixture, reference_fixture e filtro do check. Execute ambos: baseline
deve falhar pelo motivo previsto, referência deve passar. Aguarde os processos
terminarem antes de afirmar sucesso.

Use uma intenção por micropasso, hints e conceitos reais. Não acrescente nodes
redundantes para aumentar contagem. Não preencha reviewed_by/playtested humanos.
Use o GOCACHE fornecido pelo executor; não remova caches compartilhados.

Entregue por tarefa: arquivos, comportamento antes/depois, checks executados,
resultados e pendências. Em revisão, resolva e relate cada achado por ID.
