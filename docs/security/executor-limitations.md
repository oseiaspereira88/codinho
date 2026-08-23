# Limites de segurança do executor de checks

## O que o executor garante

- Nunca roda um comando livre: só `check_id`s já declarados no desafio
  fixo da sessão, resolvidos contra um allowlist fixo de runners
  (`go_test`, `go_test_race`, `go_vet`, `gofmt_check`, `go_build`,
  `go_benchmark`, `internal_ast`) — `internal/checks/registry.go`.
- Nunca usa shell: cada argumento é um elemento de `exec.Cmd.Args`,
  nunca interpolado numa string de shell.
- Ambiente restrito a um allowlist fixo de variáveis; rede negada por
  padrão (`GOPROXY=off`), só liberada quando o próprio desafio autora
  `checks[].network: true` — nunca por parâmetro de chamada.
- Saída capturada é redigida e limitada em tamanho antes de virar
  evidência.
- Processo roda em seu próprio grupo (onde suportado) e é morto por
  timeout — nunca vaza um processo órfão indefinidamente.

## O que o executor não garante (limites honestos, não escondidos)

- **Não é um sandbox de kernel.** Um check `go_test`/`go_build` roda com
  os mesmos privilégios do usuário local que executa `codinho serve`.
  Código Go deliberadamente hostil dentro do desafio ou da fixture do
  aluno tem o mesmo poder que teria se o aluno rodasse `go test`
  manualmente no terminal. Isolamento forte (container, VM, gVisor) é
  uma capacidade pós-V1 explicitamente fora do escopo desta spec e de
  safe-check-executor (Non-goal).
- **Confia no toolchain `go`/`gofmt` do host.** Não valida a integridade
  desses binários.
- **`internal_ast` é a única execução verdadeiramente sandboxed** (não
  spawna processo — só faz parse com `go/parser`), mas por isso mesmo só
  serve para verificação estrutural, nunca para rodar testes reais.
- **Timeout não é uma garantia de billing/recursos** — limita tempo de
  parede, não CPU/memória; um check pode consumir muita memória dentro
  do timeout.

## Modo entrevista: limites de integridade

Uma simulação de entrevista (`mode: interview`) roda na mesma máquina do
aluno, sem vigilância, gravação de tela ou proctoring — ver
`internal/assessment.integrityNote`, sempre repassado no
`interview_report`. Isso é uma escolha deliberada (interview-mode
Decision 1: "declarar limites e focar prática, não certificação"), não
uma lacuna a esconder: o resultado é evidência pedagógica para o
próprio aluno, nunca uma credencial verificável por terceiros.

## Testes reservados

Nenhum teste deste repositório executa contra a rede real ou against
uma máquina de terceiros. `govulncheck` e os testes de segurança
(`*_security_test.go` / `internal/security`) rodam totalmente offline e
são parte do gate de fechamento de toda spec (ver AGENTS.md).
