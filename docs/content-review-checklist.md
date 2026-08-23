# Checklist de revisão de conteúdo

Para quem revisa um desafio antes de `publication.status: published`.
**Deve ser uma pessoa diferente do autor** (Constraint de
catalog-authoring-quality) — `published_same_reviewer` bloqueia a
publicação se `reviewed_by` repetir `author`.

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
