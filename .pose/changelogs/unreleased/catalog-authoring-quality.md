---
spec: catalog-authoring-quality
category: added
breaking: false
refs:
---

`codinho catalog validate` agora aplica um gate editorial completo além
da validação de schema: pistas em ordem crescente sem gabarito
prematuro, competência e critério de aceite obrigatórios, metadados de
publicação com revisor diferente do autor, e coerência do grafo de
relações. `--checks` executa de verdade os checks declarados contra a
fixture materializada; `--v1-gate` confere os limiares de conclusão da
V1 sob demanda.
