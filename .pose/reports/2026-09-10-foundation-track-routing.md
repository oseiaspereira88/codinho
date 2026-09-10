# Roteamento das trilhas fundamentais — 2026-09-10 UTC

Spec: go-foundations-packs, permanece in-progress.

Cinco sequências autoradas em go-first-steps 1.2.0: Go do zero (11), transição
orientada a objetos (10), fluência prática (9), dados e tipos (10), testes e
design (7). O ID histórico go-from-zero foi preservado. As 47 posições
reutilizam desafios e não aumentam o inventário canônico. Comparação dos
objetos YAML confirmou a preservação integral dos desafios, conceitos,
competências e relações existentes.

TestFoundationTracksOverRealStdio copia os packs reais para uma instância
isolada, usa serve --authoring, lista exatamente cinco trilhas, inicia cada
uma e alcança suas instruções em ordem, verificando cursor, identidade e
proveniência draft. A travessia usa override explícito e não constitui
playtest, resolução dos exercícios ou aprovação pedagógica da ordem.

Validação: matriz completa 23/23 checks aprovados, zero skips, incluindo o
check registrado foundation-track-routing; 43/43 checks editoriais verificados.
Docs-check, check --strict e ready-check passaram. Nenhum conteúdo publicado,
metadata de revisor humano ou playtest foi declarado. Não houve alteração
do contrato MCP: as sequências usam os seletores já implementados.
