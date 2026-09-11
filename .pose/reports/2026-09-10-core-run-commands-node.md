# Node de limpeza e continuidade após panic

Spec: go-foundations-packs, in-progress.

O desafio combinado de execução de comandos recebeu uma verificação que une
os efeitos observáveis da recuperação: o comando que panica ainda registra sua
limpeza, e o comando seguinte continua sendo executado. A fixture já cobre a
sequência normal, panic e continuação; o node a torna explícita no fluxo.

`go run ./cmd/codinho catalog validate --checks --json` passou com diagnostics
null e editorial null. `git diff --check` passou; revisão humana e playtest
permanecem pendentes.
