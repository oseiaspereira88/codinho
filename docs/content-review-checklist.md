---
title: "Checklist de revisão de conteúdo"
doc_type: howto
---

# Checklist de revisão de conteúdo

Para quem revisa um desafio antes de `publication.status: published`.
**Deve ser uma pessoa diferente do autor** (Constraint de
catalog-authoring-quality) — `published_same_reviewer` bloqueia a
publicação se `reviewed_by` repetir `author`.

## Pré-revisão automatizada por outro agente

Antes da revisão humana, cada lote de `go-foundations-packs` (e specs de
pack irmãs que adotarem o mesmo processo, ver Decision 3 dessas specs)
passa por uma pré-revisão automatizada de um agente independente do
autor, em modo não-interativo, sandboxed e sem editar arquivos. Ela
aplica os mesmos itens deste checklist por leitura estática e roda
`codinho catalog validate`, achando cedo vazamentos de solução no texto
disclosado, micropassos com mais de uma intenção, critérios de aceite não
verificáveis e checks sem evidência executável real por trás.

Como invocar (papéis, template de prompt, por que retomar sessão em vez
de recomeçar a cada rodada): ver
[`docs/agent-review-workflow.md`](agent-review-workflow.md) e a skill
`agent-batch-review`. Mecânica de invocação:
[`scripts/agent-review.sh`](../scripts/agent-review.sh).

Isso **não substitui nada abaixo**: o veredito da pré-revisão nunca
preenche `reviewed_by` nem `playtested` — só reduz o que sobra para o
revisor humano avaliar. Os itens de "Playtest" abaixo continuam exigindo
uma pessoa real, sempre.

## Antes de revisar

- [ ] `codinho catalog validate` roda sem erro bloqueante para o pack.
- [ ] `codinho catalog validate --checks` roda sem erro (se o desafio
      declarar `fixture` + `checks`).
- [ ] Você **não** é a pessoa listada em `publication.author`.

## Conteúdo pedagógico (julgamento humano — nenhuma automação decide isto)

- [ ] O `brief` descreve o problema sem entregar a solução.
- [ ] Cada passo `micro` tem uma única intenção clara (mesmo quando o
      linter `compound_micro_instruction` não disparou — o linter é
      conservador e pode deixar passar casos reais).
- [ ] Os hints sobem em custo pedagógico real, não só em número — nível 1
      realmente ajuda menos que nível 3.
- [ ] Critérios de aceite descrevem comportamento observável, verificável
      por evidência real (não "o código está bom").
- [ ] Se o desafio pede reflexão, a pergunta exige explicar o *porquê*,
      não só descrever o *o quê*.
- [ ] Para `kind: debug`: a fixture reproduz um bug real e único; a
      correção mínima esperada não exige reescrever a função inteira.

## Playtest

- [ ] Você (ou outra pessoa que não o autor) executou o desafio de ponta
      a ponta — não apenas leu o YAML — usando `codinho workspace
      prepare` quando houver fixture, e uma sessão real via MCP/skill.
- [ ] O tempo estimado (`estimated_minutes`) é realista para quem já tem
      os pré-requisitos.
- [ ] Nenhuma pista, comentário de código ou nome de arquivo revela a
      solução antes do momento pretendido.

## Antes de marcar `publication.status: published`

- [ ] `author`, `reviewed_by` (diferente de `author`) e
      `playtested: true` estão preenchidos honestamente — `playtested:
      true` é uma afirmação de que o passo acima realmente aconteceu, não
      um checkbox de conveniência.
