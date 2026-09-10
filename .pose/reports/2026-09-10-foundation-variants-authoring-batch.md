# Variantes fundamentais — 2026-09-10 UTC

Spec: go-foundations-packs, permanece in-progress.

Oito variantes atômicas draft: quatro em go-io (leitura fragmentada, campos
JSON, nomes CSV, falha de Writer) e quatro em go-testing (duração restante,
fronteira temporal, isolamento de snapshot, contagem de chamadas). Cada uma
usa variant_of e canonical false, com uma competência compartilhada com sua
origem. Os packs passam a 1.1.0. As 45 bases e todo seu conteúdo foram
preservados por comparação dos objetos YAML com HEAD anterior ao lote.

As variantes acrescentam 32 nodes próprios. Não aumentam a contagem de bases
nem compensam os 182 nodes ainda ausentes na decomposição dessas bases.
A meta de 44 não foi alterada. Não houve publicação, playtest nem revisão
humana; a equivalência pedagógica ainda exige avaliação.

Validação: matriz strict completa 23/23, zero skips; fixtures baseline/reference
incluem as oito variantes. TestFoundationVariantOrigins verifica vínculo,
classificação e competência de origem nos packs reais. Docs-check e ready-check
passaram. A revisão local dos contratos e fixtures ocorreu separadamente da
geração, pelo mesmo agente; não constitui uma avaliação humana independente.

Os passos intermediários pedem observação qualitativa. O check final verifica
o comportamento da solução, sem certificar a qualidade dos testes do aluno.
Não foram enviados reports ou issues externos.
