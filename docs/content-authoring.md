# Autoria de conteúdo curricular

Guia para quem escreve packs (`packs/*.yaml`) — temas, conceitos,
competências, trilhas e desafios. O contrato formal vive em
`schemas/*.schema.json`; o loader em `internal/curriculum` é a
implementação de referência. Campos YAML desconhecidos são rejeitados. O corpus
`testdata/catalog-quality/schema-corpus.json` exercita schemas e loader com os
mesmos documentos. Regras de grafo, identidade e revisão entre pessoas são
semânticas adicionais do Go; rascunhos mínimos legados continuam carregáveis.

## Antes de escrever

1. Rode `codinho catalog validate` no seu pack isolado (ou combinado com
   os demais) para ver diagnósticos estruturais (schema, IDs, referências,
   ciclos, versões, paths de fixture) e achados editoriais (pistas,
   competências, relações) antes de pedir revisão.
2. Todo `id` é global ao catálogo inteiro, não só ao seu pack — escolha
   IDs namespaced (ex.: `go-debug.slice-off-by-one`) para evitar colisão.
3. `version` de pack e de desafio segue `major.minor.patch` estrito
   (`1.0.0`), sem sufixo de pré-release.

## Regras estruturais (bloqueantes, `internal/curriculum.Validate`)

- IDs únicos, referências resolvidas, sem ciclos em relações de
  precedência (`requires`, `recommended_before`, `deepens_into`).
- `fixture[].path` nunca é absoluto nem contém `..`.
- Todo critério (`criteria`) exige ao menos uma `evidence.strategies`.

## Regras editoriais (`internal/curriculum.RunEditorialChecks`)

| Regra | Severidade | O que verifica |
|---|---|---|
| `hints_increasing_order` | bloqueante | Níveis de hint devem crescer estritamente. |
| `hint_reveals_solution_early` | bloqueante | Nenhum hint abaixo do nível 6 pode ter `kind: solution`. |
| `brief_contains_code_fence` | bloqueante | `brief`/`acceptance` nunca contém um bloco de código (```). |
| `missing_competency` | bloqueante | `competencies.primary` não pode ficar vazio. |
| `missing_acceptance` | bloqueante | `acceptance` não pode ficar vazio. |
| `published_without_author`/`published_without_review`/`published_same_reviewer`/`published_without_playtest` | bloqueante | Só se aplica quando `publication.status: published` — ver seção de publicação abaixo. |
| `compound_micro_instruction` | aviso | Objetivo de passo `micro` parece juntar dois verbos independentes. |
| `missing_hints_for_gated_step` | aviso | Passo que exige avaliação positiva sem nenhum hint autorado. |
| `missing_reflection` | aviso | Nenhuma pergunta de reflexão em toda a árvore do desafio. |
| `relation_isolated` | aviso | Item sem nenhuma `relation` além de `requires`/`prerequisites`. |
| `relation_not_reciprocal` | aviso | `contrasts_with`/`commonly_fails_with` sem a aresta recíproca. |

Achados de aviso nunca bloqueiam `catalog validate`; eles existem para
guiar revisão humana (Decision 1 de catalog-authoring-quality: automação
não finge compreender pedagogia).

## Fixtures (`fixture:`)

Um desafio `kind: debug` (ou qualquer outro que precise de código
inicial) declara `fixture: [{path, content}]`. Nunca escreva esse
conteúdo manualmente no workspace do aluno — só
`codinho workspace prepare <challenge-id> --dest <path>` materializa.

Se o desafio declarar checks, cada um precisa de expectativa para o código
inicial e para uma referência completa. Esses arquivos privados são usados por
`codinho catalog validate --checks`; nunca pelo workspace do aluno.

```yaml
validation:
  reference_fixture:
    - path: main.go
      content: "package main\n"
  expectations:
    - check_id: parses
      baseline: syntax_failure
      reference: pass
```

Esse exemplo pressupõe um check `internal_ast` chamado `parses` e uma fixture
inicial intencionalmente inválida. A referência deve formar um workspace completo,
incluindo `go.mod` e testes quando o runner exigir. Ela não sobrepõe a fixture
inicial. Para desafios sem código inicial distribuído, use `baseline_fixture`
dentro de validation e uma `justification` explicando essa alternativa
reproduzível. Um texto dizendo que o check passou em outro lugar não é prova.

Baseline aceita `pass`, `test_failure`, `compile_failure`, `format_failure`,
`analysis_failure` ou `syntax_failure`; referência exige `pass`. Falha não prevista,
timeout, skip, nenhum teste/benchmark executado e saída truncada não passam.
Testes usam JSON e `-count=1`; benchmarks usam também `-benchtime=1x` e precisam
produzir uma medição. Cada check/cenário roda em diretório temporário próprio.
O executor usa runners permitidos, sem shell, com timeout e captura limitada;
nega downloads de módulos. Isso não é isolamento de rede do código executado:
fixtures são código local confiável, como no contrato existente de `--checks`.

O relatório contém IDs/versões de pack e desafio, cenário, expectativa, resultado
e digest do check com seus arquivos. Também identifica revisão do binário,
modificações locais e versão Go quando disponíveis; revisão desconhecida é
indicada como `unknown`. CI associa o artefato ao commit da execução. Conteúdo da
referência e stdout/stderr do programa não entram no relatório, no catálogo
público ou nos eventos de sessão.

## Publicação (`publication:`)

Pack e desafio começam como rascunhos quando publication está ausente, vazio ou
com status draft. Para publicar, preencha metadados verdadeiros em ambos:

```yaml
publication:
  status: published
  author: <seu nome ou handle>
  reviewed_by: <pessoa DIFERENTE de author>
  playtested: true   # somente após playtest real
```

A publicação do pack governa temas, conceitos, competências e trilhas. Cada
desafio requer sua própria publicação e referências a conteúdo publicado. Uma
publicação inválida bloqueia a inicialização do MCP, inclusive em autoria.
Automação verifica declarações, mas não autentica pessoas nem confirma playtests.
A pré-revisão por outro agente não preenche esses metadados.

`codinho serve` usa somente o catálogo publicado para busca, recomendação e novas
sessões. `codinho serve --authoring` inclui rascunhos para autoria/playtest local.
Comandos administrativos list/show continuam permitindo inspecionar inventário.
Sessões anteriores retomam o conteúdo público que já estava fixado no histórico.
Nenhum rascunho é promovido automaticamente e nenhum trabalho autorado é apagado.

## Cobertura e gate V1

O relatório separa `inventory`, `drafts`, `published` e `eligible`. Nos desafios,
`canonical: true` é uma declaração explícita de curadoria; protótipos sem essa
marca e derivados com `variant_of: <id>` não contam como novos desafios elegíveis.
O array reservado `variants` também não aumenta a contagem. Revise equivalência e
repetição pedagogicamente antes de marcar um desafio canônico.

`codinho catalog validate --v1-gate` usa cobertura elegível (160 conceitos,
100 competências, 84 desafios, 12 trilhas, 500 nós), distribuição exata de
`packs/distribution.json` e execução de todos os checks publicados. A política
versionada define tipos globais e por pack ou grupo planejado. Depuração,
refatoração e revisão compartilham o alvo investigation; backend e produção ainda
usam seus totais de grupo, sem inventar uma divisão entre packs. Tipos inesperados
e packs publicados fora da política falham. Use `--distribution <arquivo.json>`
para validar uma política explicitamente; ela deve ter totais coerentes.

`--published-checks` executa todos os checks publicados e informa quantos foram
verificados. CI usa esse comando e arquiva `catalog-proof.json`. Zero publicados
produz zero verificados, sem afirmar prontidão V1. `--checks` cobre também
rascunhos e diagnostica metadados ausentes; `catalog validate` sem esses flags
continua utilizável durante autoria incremental. V1 permanece um gate explícito:
um catálogo numericamente suficiente de rascunhos deve falhar.


## Sequência e alternativas da árvore

Use IDs não vazios e únicos dentro de cada desafio, incluindo suas layers.
Declare kind macro, meso ou micro nos passos. Mantenha os filhos em ordem
obrigatória; omitir children_mode equivale a sequence.

Declare children_mode: choice somente quando os filhos forem alternativas
exclusivas. Inclua ao menos duas alternativas com IDs distintos. Numa sessão
mais detalhada que esse nó, sua instrução serve como checkpoint da decisão;
a seleção ocorre por step_advance com next_step_id, depois da conclusão
explícita do checkpoint. Não use choice para listar ações que precisam ser
executadas em sequência.

Uma sessão que conclua o próprio nó numa janela ampla cobre sua subárvore
sem fabricar conclusões individuais. Consulte [compatibilidade](compatibility.md#navegação-da-árvore-da-sessão)
para cursor, agrupamento e retomada.
