# Fluxo de autor/revisor entre agentes

Padrão para quando dois processos de IA independentes colaboram na autoria
de conteúdo do catálogo (ou em qualquer outra tarefa que se beneficie de
uma segunda opinião automatizada antes da revisão humana): um agente
autora, outro revisa de forma adversarial, sem editar nada. Nasceu em
`go-foundations-packs` Decision 3 (Codex como pré-revisor padrão de
`catalog-authoring-quality`) mas é agnóstico ao par de agentes e à tarefa —
documentado aqui para não se perder entre sessões.

## Papéis são posições, não identidades fixas

Nada aqui amarra "Claude autora, Codex revisa". Os papéis são:

- **Autor**: produz ou modifica arquivos. Sandbox com escrita liberada
  (`workspace-write`), permissão explícita do que pode editar.
- **Revisor**: só lê e roda comandos de validação (build, test, lint,
  `catalog validate`). Sandbox também pode ser `workspace-write` (porque
  rodar testes exige escrever em cache/tmp), mas a restrição de "não editar
  nada versionado" é feita **no prompt**, não no sandbox — e verificada
  depois com `git status --short` (deve continuar limpo).

Qualquer CLI de agente não-interativa pode ocupar qualquer papel:
Codex autora e Claude revisa, duas instâncias de Codex, duas instâncias de
Claude (uma delas via este mesmo Claude Code em modo headless, se
disponível), uma futura CLI da Antigravity (`agy`) etc. Trocar quem ocupa
qual papel é só trocar qual comando o script invoca — ver
[`scripts/agent-review.sh`](../scripts/agent-review.sh).

## O script

`scripts/agent-review.sh <new|resume> <agent> <output-file> <prompt|@arquivo>`
é a primitiva de invocação: não sabe o que é "autor" ou "revisor", só sabe
iniciar ou retomar uma sessão não-interativa e sandboxed de um backend
(hoje: `codex`, via `codex exec`). Adicionar outro backend é um `case` novo
no script, sem tocar quem o chama.

O backend `codex` fixa explicitamente `model="gpt-5.6-luna"` e
`model_reasoning_effort="high"` (via `-c` no `codex exec`) em vez de
depender do default de `~/.codex/config.toml` — esse é o modelo padrão do
revisor neste projeto. Fixar no script em vez de confiar no config global
evita que a propriedade "modelo/fornecedor independente do autor" da
revisão quebre silenciosamente se o config global mudar por outro motivo.
Para usar outro modelo numa rodada pontual, edite o array `codex_model` no
script (não há flag de override por chamada hoje).

## Por que retomar sessão em vez de começar do zero a cada rodada

Uma chamada nova (`codex exec ...`) sem contexto prévio precisa reler os
arquivos relevantes (checklist, docs de autoria, o pack inteiro) toda vez —
isso apareceu nesta própria sessão: cada rodada nova de revisão consumiu
~60–100 mil tokens, a maior parte só relendo arquivo que a rodada anterior
já tinha lido. `codex exec resume --last "<o que mudou>"` continua a mesma
sessão: o agente já tem os arquivos e o veredito anterior em contexto, e só
precisa avaliar o diff/a mudança relatada — muito mais barato e mais rápido
por rodada, e reduz o risco do revisor "esquecer" um achado da rodada
anterior que ele mesmo levantou.

Regra prática: use `new` na primeira rodada de um lote; use `resume` em
toda rodada seguinte do mesmo lote, descrevendo só o que mudou desde o
veredito anterior — nunca colando de novo o conteúdo dos arquivos que o
agente já leu.

## Template de prompt — revisor

```
Você é um revisor de conteúdo independente do autor, no repositório
<repo>. NÃO edite nenhum arquivo — comandos read-only (leitura, build,
test, lint, validate) são permitidos. Não faça perguntas, produza só o
relatório final.

Contexto: <o que foi adicionado/mudado e onde>.

Faça, nesta ordem:
1. Leia <checklist/guia de autoria relevante>.
2. Leia <os arquivos/trechos específicos a revisar>.
3. Rode <comando de validação determinística> e reporte o resultado.
4. Avalie <critérios do checklist>, marcando o que exige uma pessoa real
   (playtest, julgamento pedagógico) que você não pode substituir.
5. Aponte vazamento de solução, critério não verificável, instrução com
   mais de uma intenção, e qualquer inconsistência entre constraints e a
   solução esperada.
6. Dê um veredito final por item revisado: aprovado sem ressalvas /
   aprovado com ressalvas menores (liste) / rejeitado (liste motivos
   bloqueantes), deixando claro que é uma pré-revisão automatizada e não
   substitui revisão/evidência humana exigida pelo gate do projeto.
```

Rodadas seguintes (`resume`): repita só os itens 3–6 focados no que mudou,
citando exatamente o que foi corrigido desde o veredito anterior.

## Template de prompt — autor

```
Você é o autor de <o que precisa ser criado/corrigido>, no repositório
<repo>. Siga <guia de autoria relevante>. Pode editar somente
<arquivos/paths explicitamente liberados>. Ao final, rode <validação
determinística> e reporte o resultado; não afirme sucesso sem rodar.
```

## Matar de verdade uma rodada travada antes de tentar de novo

`TaskStop` (ou o equivalente do seu harness) mata o wrapper de shell que
lançou `codex exec`, **não necessariamente o processo `codex` em si** —
ele pode continuar rodando, órfão, escrevendo na mesma sessão. Sintoma
observado nesta sessão: depois de matar uma rodada `new` travada há mais
de 1h30, uma nova rodada `new` no mesmo diretório ficou parada por mais
de 30 minutos sem produzir nada — o processo órfão anterior (confirmado
com `ps aux | grep codex`, ainda consumindo CPU) estava disputando lock
de sessão com a tentativa nova. Matar o PID do `codex` diretamente (`kill
<pid>` nos processos `node /usr/bin/codex` e o binário `codex` filho)
resolveu na hora.

Antes de iniciar uma rodada `new`/`resume` depois de matar uma travada,
confirme com `ps aux | grep -i "codex exec"` que não sobrou nada rodando
no mesmo diretório — senão a rodada nova pode travar por disputa de lock,
não por estar de fato processando algo.

## Cache do Go dentro do sandbox do revisor

Se o revisor precisar rodar `go test`/`go build` de verdade (para provar
um `checks:` real, não só ler o YAML), o sandbox `workspace-write` do
Codex só permite escrita em `workdir`, `/tmp` e `$TMPDIR` — `GOCACHE`
aponta por padrão para `~/.cache/go-build`, fora dessa allowlist, e o
comando falha com `read-only file system` antes mesmo de chegar ao teste.
Peça explicitamente no prompt para redirecionar: `GOCACHE=/tmp/<algo>
GOFLAGS=-mod=readonly GOPROXY=off go test ...`. Validado nesta sessão:
sem isso o revisor não conseguia confirmar evidência executável nenhuma;
com isso, rodou e confirmou o teste real em segundos.

## Cache de `go test` pode mascarar um mutante que deveria falhar

Ao verificar manualmente que uma implementação incorreta ("mutante")
realmente falha contra um teste — a técnica usada para dar evidência
comportamental real a critérios que só têm `source_inspection` por trás —
`go test ./...` sem `-count=1` pode devolver um resultado em cache mesmo
depois de trocar o `.go` da implementação no mesmo diretório de trabalho
(observado em `go-data-text` checkpoint 3: uma implementação de
`Dedupe` deliberadamente quebrada, que perde a ordem original ao montar
a saída iterando um map, "passou" na primeira rodada — só porque o
resultado de uma execução anterior, com um arquivo diferente, ainda
estava em cache). Rodar de novo com `-count=1` revelou a falha real, de
forma consistente em 5 repetições. Sempre que verificar um mutante
manualmente (autor ou revisor) contra um teste que envolve ordem de
iteração de map ou qualquer outra fonte de não-determinismo, use `go test
./... -count=1` (e rode algumas vezes) — nunca confie em uma passagem
sem `-count=1` depois de trocar o arquivo de implementação.

## Ganho real medido de retomar sessão

Nesta sessão, quatro rodadas de revisão do mesmo lote (`go-foundations-
packs` checkpoint 2): as duas primeiras foram `new` (sem contexto prévio)
e consumiram ~60–100 mil tokens cada, relendo os mesmos arquivos. Uma
terceira rodada `new` (pedindo pro revisor materializar e rodar teste do
zero) ficou presa mais de 1h30 sem concluir e precisou ser morta. Ao
`resume`-la e pedir só o passo específico que faltava (rodar dois
comandos com o GOCACHE corrigido), a rodada final custou ~6,6 mil tokens
e terminou em segundos. `resume` não é só mais barato — para uma tarefa
que trava, resumir e pedir algo mais específico também é como se
recupera sem perder o trabalho já feito na sessão travada.

## O que isso não substitui

Pré-revisão automatizada (por qualquer agente) nunca preenche metadados de
revisão humana exigidos por gates de qualidade do projeto (ex.:
`publication.reviewed_by`/`playtested` em `catalog-authoring-quality`) —
ela só reduz o que sobra para a pessoa avaliar. Ver Decision 1 de
`catalog-authoring-quality`: "automação não deve fingir compreender
pedagogia".

## Sobre A2A (Agent2Agent)

O protocolo A2A resolve um problema diferente do nosso: descoberta e
delegação de tarefas entre agentes de organizações/redes diferentes, via
HTTP/JSON-RPC, Agent Cards e ciclo de vida de task remoto. Aqui os dois
agentes rodam como subprocessos locais na mesma máquina, no mesmo
repositório, sem fronteira de rede ou de confiança entre eles — não há
descoberta a fazer nem identidade a negociar. Adotar A2A agora seria
construir um servidor/cliente A2A só para invocar um binário que já está
no PATH da própria máquina. Vale reconsiderar se este padrão precisar
orquestrar um revisor remoto/de outra organização (não mais um subprocesso
local) — até lá, `scripts/agent-review.sh` é a camada certa.
