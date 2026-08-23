# Autoria de conteúdo curricular

Guia para quem escreve packs (`packs/*.yaml`) — temas, conceitos,
competências, trilhas e desafios. O contrato formal vive em
`schemas/*.schema.json`; o loader em `internal/curriculum` é a
implementação de referência que sempre vence em caso de divergência.

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

Se o desafio também declarar `checks`, rode
`codinho catalog validate --checks` para confirmar que a fixture forma um
workspace executável para cada check (ex.: um check `go_test` precisa de
`go.mod` dentro da própria fixture).

## Publicação (`publication:`)

Um desafio começa como rascunho (`publication.status` vazio ou `draft`) e
nunca é verificado por essas regras nesse estado. Para marcar
`status: published`, preencha:

```yaml
publication:
  status: published
  author: <seu nome ou handle>
  reviewed_by: <pessoa DIFERENTE de author>
  playtested: true   # só true depois de um playtest real
```

`playtested: true` é uma afirmação de honra — nenhuma automação consegue
confirmar que um playtest realmente aconteceu (Decision 1).

## Cobertura e gate V1

`codinho catalog validate --v1-gate` compara o catálogo carregado contra
os limiares do roadmap V1 (160 conceitos, 100 competências, 84 desafios,
12 trilhas, 500 step nodes). Ele fica desligado por padrão porque um
catálogo em progresso está, por definição, abaixo desses números — use-o
apenas ao se aproximar do aceite V1.
