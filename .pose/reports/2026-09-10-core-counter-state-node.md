# Node de estado independente em Counter

Spec: go-foundations-packs, in-progress.

O desafio de receptor de ponteiro recebeu uma verificação explícita de que
instâncias separadas de Counter não compartilham estado. A fixture já executa
duas instâncias em sequências distintas; o novo node torna esse contrato
observável no percurso pedagógico.

`go run ./cmd/codinho catalog validate --checks --json` passou com diagnostics
null e editorial null. `git diff --check` passou. O contrato e a fixture não
foram alterados; revisão humana e playtest permanecem pendentes.
