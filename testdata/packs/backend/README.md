# Fixtures backend

As fixtures executáveis do catálogo backend são mantidas dentro de
`packs/go-concurrency-context.yaml`, `packs/go-http.yaml` e
`packs/go-database.yaml`, onde `catalog validate --checks` as carrega e
verifica com confinamento de paths.

Este diretório é o ponto reservado para cópias de workspace preparadas por
`codinho workspace prepare`. Ele não duplica a fonte canônica nem contém
credenciais, banco local persistente ou dependência de rede. A extração de
fixtures para uso manual permanece uma pendência editorial separada do gate
determinístico dos packs.
