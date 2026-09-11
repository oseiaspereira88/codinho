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
  --model gpt-5.6-luna --reasoning high --parallelism 3 --batch-size 5
python3 scripts/editorial.py status --run-dir /tmp/codinho-editorial-run
python3 scripts/editorial.py run --run-dir /tmp/codinho-editorial-run
```

Escolha configuração antes da execução; defaults conservadores usam um autor.
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

Cada fila contém spec e tasks com id, title, paths e acceptance. Preserve IDs
entre rodadas. Paths são arquivos exatos relativos ao repo. Tarefas que dependem
de mudanças integradas devem entrar em uma execução posterior baseada no novo
HEAD. Não execute a fila inteira antes de triar dependências e gates humanos.

Falha/timeout preserva worktree e logs. Retome o lote pela sessão registrada;
se não houver sessão utilizável, investigue a falha antes de iniciar outra.
Conflitos de integração não autorizam sobrescrita: reconcilie, valide e revise
o novo patch. Integração exige árvore principal limpa.

## Validação e piloto

```sh
python3 -m unittest discover -s scripts -p test_editorial.py
```

Os testes usam backend falso e repositórios temporários. O piloto real usa cinco
tarefas e luna/high; o coordenador verifica cada uma e registra o resultado na
spec. Execução iniciada não significa execução concluída: aguarde exit code e
evidência final. Matriz completa após integração é um gate adicional.

Consulte [inventário](editorial/go-foundations-inventory.md) e
[papéis de autor/revisor](agent-review-workflow.md). O script legado
agent-review.sh continua restrito ao uso sequencial; use o novo executor para
sessões explícitas e paralelismo. CLI verificada contra `codex exec --help` e
[documentação oficial](https://developers.openai.com/codex/noninteractive).
