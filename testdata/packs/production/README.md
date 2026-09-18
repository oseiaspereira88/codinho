# Fixtures de produção

As fixtures executáveis de `go-production` e
`go-architecture-debugging` ficam embutidas nos packs
[`packs/go-production.yaml`](../../../packs/go-production.yaml) e
[`packs/go-architecture-debugging.yaml`](../../../packs/go-architecture-debugging.yaml).
Essa é a fonte canônica consumida por `catalog validate --checks`; mantê-la nos
packs evita duplicação divergente.

Esta pasta contém somente este README e não tem fixtures materializadas.
Quando necessário, `codinho workspace prepare <challenge-id> --dest <path>`
prepara uma cópia de trabalho no destino escolhido, sem transformar esta pasta
em uma segunda fonte de conteúdo.
