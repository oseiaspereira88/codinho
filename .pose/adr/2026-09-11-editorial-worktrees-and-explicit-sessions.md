# ADR: Worktrees editoriais e sessões explícitas

## Status
Accepted — 2026-09-11. Spec: editorial-batch-engine.

## Context
resume --last não identifica lotes paralelos; autores no mesmo diretório podem
sobrescrever o mesmo pack YAML.

## Decision
Isolar cada lote em worktree e persistir ID da sessão e configuração. O
coordenador revisa o patch por digest e integra serialmente em árvore limpa.
Codex CLI recebe modelo e reasoning explícitos. Não adicionar SDK/API própria.

## Consequences
Preservar logs e trabalhos em falhas. Conflitos exigem nova revisão. Até oito
subprocessos são configuráveis, sujeitos aos recursos e permissões locais.
Não usar --last, aprovação automática ou integração cega.

## Alternatives
Diretório compartilhado descartado por colisões. API própria adiada porque
a CLI já fornece sessões e execução não interativa para o piloto.
