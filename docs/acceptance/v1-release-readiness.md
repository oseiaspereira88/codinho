---
title: Prontidão de entrega V1
doc_type: reference
---

# Prontidão de entrega V1

Estado em 2026-09-08: gates implementados e CI nativa conferida. Não há aceite V1 nesta rodada.
Consulte a [spec](../../.pose/specs/2026-09-07-v1-delivery-ci-assurance.md) e
os [critérios do roadmap](../../.pose/roadmaps/codinho-v1.md).

| Requisito | Evidência necessária | Estado |
|---|---|---|
| R1 | Instalação verificada e gates POSE estritos | Ferramentas verificadas; gates locais passaram, incluindo validate. |
| R2 | Matriz única, race, formatação, contratos e fixtures | 19 de 19 checks passam; 38 checks editoriais verificados em 76 cenários. |
| R3 | Jobs nativos Linux e macOS no commit candidato | Run 34283684515 passou nas duas plataformas no commit b9b8868. |
| R4 | Actions por commit, scanner verificado, release sem publish | Testes de pins e rejeição de labels inválidos passam localmente. |
| R5 | Proveniência reconciliada, roots e contratos | artifact-check estrito e surface-check passaram no b9b8868; avisos históricos de órfãos permanecem explícitos. |
| R6 | C1–C9 atuais, negativos de evidência | Negativos locais passam; roll-up do roadmap pendente. |
| R7 | Manifest e comandos/links válidos | Manifest com 13 documentos; docs-check passa sem erros/avisos. |
| R8 | Relatórios por commit/plataforma/host e revisão independente | Relatórios e logs remotos conferidos; decisão formal no registro de revisão da spec. |

Registre em cada candidato o commit, resultado, OS, arquitetura, versão de Go,
host e run/attempt. A CI arquiva `delivery-validation.json` e `native-ci.json`
por plataforma. Exija resultados atuais de todos os checks da matriz; relatório
ausente, antigo, com skip ou check desconhecido deve falhar.

Execute o gate quantitativo/editorial do catálogo somente no candidato completo.
Não substitua playtest humano, validação de hosts de tutor ou revisão independente
do threat model por testes unitários ou por configuração de workflow.

Validação local de 2026-09-08: `make check` com POSE 1.7.12 e scanner 1.6.0,
Go go1.26.5-X:nodwarf5, Linux/amd64, host local, base Git
`dbd848e927ff636e65c2a38220f029e29a07057b` com alterações no worktree:
19 checks pass. As provas privadas verificaram os 38 checks nos 76 cenários.
Isso não é uma atestação de um commit imutável. O relatório nativo inclui
source_modified; a CI exige false para os arquivos autorados do candidato.

`pose roadmap-check codinho-v1 --strict` reconheceu C1–C9 e retornou
terminal=false. C3, C5 e C8 têm relatórios de aceite ausentes; as specs ainda
não terminais também bloqueiam o roadmap. Os demais critérios não receberam
aceite independente nesta execução. `surface-check` e `artifact-check`
devem ser reexecutados sobre o commit final e sua evidência atual.

## Execução nativa conferida

[GitHub Actions 34283684515](https://github.com/oseiaspereira88/codinho/actions/runs/34283684515)
executou o commit `b9b886895023ffa4e831962c5620269c2b84dcdd` em 2026-09-08, attempt 1.
Os dois artefatos foram baixados e conferidos; ambos declaram source_modified=false.

| Plataforma | Host | Go | Checks | Resultado |
|---|---|---|---|---|
| Linux/amd64 | github-actions | go1.27.1 | 19 | pass |
| macOS/arm64 | github-actions | go1.27.1 | 19 | pass |

Essa execução inclui baseline/referência dos 38 checks editoriais, protocolos,
startup e smoke-install. Não comprova Linux/arm64 ou Windows nativos, nem
publicação, aceitação pedagógica ou integração humana com os hosts de tutor.
Para uma revisão posterior, use os artefatos da execução do seu próprio commit;
a tabela identifica exclusivamente a execução acima.

Também foram conferidos `pose-check.log`, `pose-validate.latest.log` e o relatório
`2026-09-08-standard-native-ci-validation.md` em ambas as plataformas.

| Arquivo | Plataforma | SHA-256 |
|---|---|---|
| native-ci.json | Linux/amd64 | 15b7d7b16808130b10b10be7b8b4daa9b459fb08746a623f6163d9d4f8312ebc |
| delivery-validation.json | Linux/amd64 | 4457abbb06aac6ee7d0c1984da8e0c6bd3966f62ad9c086631a70628f95e4ceb |
| native-ci.json | macOS/arm64 | a512e66c2b0a125d039aecdcd401a34745a5b0d197113176fbf0220d5ca603f6 |
| delivery-validation.json | macOS/arm64 | 44eff3ea1643b2e6a5cc057516e062bb903978b37c74c209f27755e5b2744571 |

## Critérios C1–C9

Use os blockers do roadmap e os documentos de aceite, além dos resultados dos
checks. Um campo `passed` isolado não concede aceite ao critério composto.

| Critério | Evidência disponível | Pendência de aceite V1 |
|---|---|---|
| C1 | cli-reachability e cli-e2e passam no candidato nativo | Reconciliar no candidato integrado. |
| C2 | session-recovery e session-tree-progression passam | Reexecutar com o catálogo final. |
| C3 | catalog-publication-integrity e provas privadas passam | Catálogo curado e matriz de requisitos ainda pendentes. |
| C4 | tutor-skill-routing passa por stdio real | Conteúdo canônico e validação dos hosts ainda pendentes. |
| C5 | Fluxos automatizados disponíveis | Relatório de piloto dos dois hosts ausente. |
| C6 | CI nativa e threat model revisados nesta spec | Exigir revisão formal válida e evidência do candidato integrado. |
| C7 | MCP, stdout, startup e smoke-install passam | Reconciliar no candidato integrado. |
| C8 | Este relatório registra plataformas e limitações | Matriz de requisitos e piloto ainda ausentes. |
| C9 | evaluation-evidence passa por stdio real | Reexecutar no candidato integrado. |
