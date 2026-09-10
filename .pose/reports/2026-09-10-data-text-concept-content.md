# Conteúdo conceitual de go-data-text — 2026-09-10 UTC

Spec: go-foundations-packs, in-progress.

O primeiro lote acrescenta conteúdo aos quatro conceitos de slices, UTF-8,
comma-ok e citação CSV. Cada exemplo usa apenas a biblioteca padrão e cobre
isolamento de backing array, limite por rune, distinção entre zero e ausência,
e escaping de campos CSV. Os exemplos foram preparados para execução isolada;
a validação do catálogo passou sem diagnósticos bloqueantes.

As relações conceituais foram adicionadas ao grafo para cada referência nova.
Desafios, competências e conteúdo anterior foram preservados. Nenhum conteúdo
foi publicado ou avaliado por humano.

Próximo passo: completar os quatro conceitos restantes de go-data-text
(conjuntos, strconv, capacidade de slice e preservação de espaçamento).
