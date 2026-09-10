# Complemento go-errors — 2026-09-10 UTC

Spec: go-foundations-packs (in-progress). Adicionados classify-wrapped-errors
(combined) e maintain-memory-catalog (functional_slice), oito conceitos,
seis competências e 14 nodes. Pack 1.2.0, autoria codex, status draft.

Validação: pose validate --strict passou 22/22 checks, zero skips; gate
editorial 43/43 checks verificados. As duas novas baselines falham e suas
referências passam. Docs-check sem erros/warnings; ready-check passou.
Comparação dos objetos YAML antes/depois confirmou que todos os desafios,
conceitos, competências e relações anteriores foram preservados integralmente.

Revisão técnica: classificação usa identidade/tipo e prioridade definida sobre
árvores de erros. Catálogo preserva valor zero utilizável, estado em falhas,
independência de instâncias, ordenação e snapshots. Fixtures sem rede, dependência
externa ou caminhos externos. O catálogo de exercício não promete persistência,
concorrência ou segurança para receiver nil. Não houve revisão pedagógica
independente ou playtest humano; metadata correspondente não foi preenchida.

Inventário: 45 desafios fundamentais draft, 70 conceitos, 51 competências,
118 nodes. Restam decisão sobre distribuição, aprofundamento, variantes,
cinco trilhas e playtest. Não se alterou a meta nem se reclassificou conteúdo.
