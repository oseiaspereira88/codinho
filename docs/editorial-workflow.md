---
title: Esteira editorial em lotes
doc_type: howto
---

# Esteira editorial

O coordenador inventaria, delega, revisa e integra. Cada autor recebe um lote
em worktree independente. O backend usa Codex CLI já instalado; não exige SDK
ou uma nova chave API. A conta deve permitir o modelo escolhido.

## Iniciar

```sh
python3 scripts/editorial.py init \
  --queue docs/editorial/go-foundations-queue.json \
  --run-dir /tmp/codinho-editorial-run \
  --model gpt-5.6-luna --reasoning max --parallelism 3 --batch-size 5 \
  --completed-task pilot-load-value-empty \
  --completed-task pilot-finder-zero
python3 scripts/editorial.py status --run-dir /tmp/codinho-editorial-run
python3 scripts/editorial.py run --run-dir /tmp/codinho-editorial-run --limit-batches 1
```

Escolha a configuração antes da execução. O perfil padrão usa o modelo
`gpt-5.6-luna` com reasoning `max` e um autor.
Use `--completed-task <id>` para iniciar uma nova execução após integrar tarefas
de uma fila anterior; o ID permanece no histórico, mas libera dependências da
nova execução sem redisparar conteúdo já integrado.
Use `--task <id>` uma vez por tarefa para formar uma subfila em ordem explícita.
Isso permite colocar lotes de packs distintos na mesma onda, sem cruzar seus
arquivos permitidos.
Para o piloto, acrescente `--limit-batches 1` ao comando run.
Até oito autores são permitidos pelo executor, sujeitos à capacidade da conta,
memória e permissões locais. O limite dos subagentes nativos da sessão não é o
mesmo que o de subprocessos CLI. Não há fallback silencioso de modelo/reasoning.
Use um diretório persistente externo ao repo para trabalhos que precisam
sobreviver à limpeza de /tmp. Não versione credenciais ou logs de sessões.

## Revisar e integrar

```sh
python3 scripts/editorial.py review --run-dir /tmp/codinho-editorial-run \
  --batch batch-001 --decision request-changes --feedback 'Descreva o achado e seu aceite.'
python3 scripts/editorial.py revise --run-dir /tmp/codinho-editorial-run \
  --batch batch-001 --feedback 'Corrija os achados e execute a regressão.'
python3 scripts/editorial.py review --run-dir /tmp/codinho-editorial-run \
  --batch batch-001 --decision approve --feedback 'Registre checks e revisão de cada tarefa.'
python3 scripts/editorial.py integrate --run-dir /tmp/codinho-editorial-run --batch batch-001
```

Aprovação é um ato do coordenador após ler diff e evidência, nunca consequência
automática de exit code zero. Mudanças posteriores invalidam o digest aprovado.
Revisão por comando é controle operacional, não autenticação contra um agente
malicioso. Aplique o checklist e não execute review dentro do autor.

Cada fila contém spec e tasks com id, title, paths, acceptance e, opcionalmente,
depends_on (IDs anteriores). Preserve IDs entre rodadas. Paths são arquivos
exatos relativos ao repo. Cada run despacha uma onda de até parallelism lotes
prontos; aguarde revisão e integração antes de chamar run novamente.
Dependências exigem integração dos pré-requisitos. Lotes que compartilham arquivos
aguardam os anteriores. O tamanho de lote é um máximo: dependências também
separam lotes. Cada lote novo captura o HEAD atual; sua base fica fixa nas revisões.
Execuções antigas preservam sua fila original: inicialize uma subfila atualizada
com tarefas restantes para adotar dependências adicionadas posteriormente.
Não marque gates humanos como satisfeitos por ordem de execução.

Auditorias sem alterações podem ser aprovadas com evidência. Sua integração
registra no_changes e o commit auditado, sem criar commit vazio, e exige que
o HEAD continue igual à base auditada.

Falha/timeout preserva worktree e logs. Retome o lote pela sessão registrada;
se não houver sessão utilizável, investigue a falha antes de iniciar outra.
Conflitos de integração não autorizam sobrescrita: reconcilie, valide e revise
o novo patch. Integração exige árvore principal limpa.
Se a reconciliação exigir mudar a base, ou uma auditoria sem diff ficar
obsoleta, inicialize uma nova subfila dessas tarefas no HEAD atual e faça nova
revisão. revise mantém a base original; não faça rebase nem altere state.json
manualmente para preservar aprovação. Interrupções anteriores ao registro do
PID também exigem diagnóstico manual antes de criar a execução substituta.
Use revise em um lote approved para revogar a aprovação e pedir correções.
Após interrupção, use recover --run-dir ... --batch ...: ele recusa recuperação
se o PID do autor ainda existir e recupera a sessão dos logs antes de permitir
revise. Execuções legadas sem PID exigem diagnóstico manual. Falha de commit
preserva o índice em integration-failed. Corrija a causa do hook/commit, conclua
o commit com o trailer da fila e use reconcile --run-dir ... --batch ...:
o comando verifica pai, digest exato e trailer antes de registrar integrated.

## Validação e piloto

```sh
python3 -m unittest discover -s scripts -p test_editorial.py
```

Os testes usam backend falso e repositórios temporários. O piloto histórico usou
cinco tarefas e luna/high; novas execuções usam o perfil padrão luna/max. O
coordenador verifica cada uma e registra o resultado na spec. Execução iniciada
não significa execução concluída: aguarde exit code e
evidência final. Matriz completa após integração é um gate adicional.

Consulte [inventário](editorial/go-foundations-inventory.md) e
[papéis de autor/revisor](agent-review-workflow.md). O script legado
agent-review.sh continua restrito ao uso sequencial; use o novo executor para
sessões explícitas e paralelismo. CLI verificada contra `codex exec --help` e
[documentação oficial](https://developers.openai.com/codex/noninteractive).
