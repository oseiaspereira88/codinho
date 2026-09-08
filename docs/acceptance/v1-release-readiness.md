---
title: Prontidão de entrega V1
doc_type: reference
---

# Prontidão de entrega V1

Estado em 2026-09-08: implementação em andamento. Não há aceite V1 nesta rodada.
Consulte a [spec](../../.pose/specs/2026-09-07-v1-delivery-ci-assurance.md) e
os [critérios do roadmap](../../.pose/roadmaps/codinho-v1.md).

| Requisito | Evidência necessária | Estado |
|---|---|---|
| R1 | Instalação verificada e gates POSE estritos | Ferramentas verificadas; gates locais passaram, incluindo validate. |
| R2 | Matriz única, race, formatação, contratos e fixtures | 19 de 19 checks passam; 38 checks editoriais verificados em 76 cenários. |
| R3 | Jobs nativos Linux e macOS no commit candidato | Run 34273134508 passou nas duas plataformas; confira o commit dos artefatos. |
| R4 | Actions por commit, scanner verificado, release sem publish | Testes de pins e rejeição de labels inválidos passam localmente. |
| R5 | Proveniência reconciliada, roots e contratos | Policies e mapas testados; primeiro commit passou artifact-check. |
| R6 | C1–C9 atuais, negativos de evidência | Negativos locais passam; roll-up do roadmap pendente. |
| R7 | Manifest e comandos/links válidos | Manifest com 13 documentos; docs-check passa sem erros/avisos. |
| R8 | Relatórios por commit/plataforma/host e revisão independente | Relatórios remotos conferidos; revisão independente pendente. |

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

[GitHub Actions 34273134508](https://github.com/oseiaspereira88/codinho/actions/runs/34273134508)
executou o commit `649725999449ebbdb77dea0a7880fb60fe471f54` em 2026-09-08.
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
