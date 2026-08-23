---
spec: reliability-observability-compatibility
category: added
breaking: false
refs:
---

`codinho doctor` agora diagnostica de verdade: catálogo, estado local
gravável, lock órfão de uma execução anterior encerrada abruptamente,
evidência e toolchain necessária — sempre sem expor ambiente completo,
código ou segredos. Toda chamada de tool MCP passa a gerar uma linha de
log estruturada (ferramenta, duração, status, tamanho) em stderr.
Startup e busca no catálogo têm orçamentos de desempenho documentados e
comprovados.
