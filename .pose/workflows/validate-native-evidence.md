# Workflow: Native validation evidence

## Objective

Produzir evidência reproduzível para `pose validate --strict` e distinguir uma
falha do produto de uma limitação do ambiente de execução.

## Scope

Aplicar este workflow a revisões que dependam da matriz nativa, especialmente
quando Go, caches de módulos ou `govulncheck` estiverem envolvidos. Ele não
substitui os checks específicos declarados pela spec.

## Preconditions

- Identificar a spec, o commit avaliado e os checks exigidos.
- Usar um diretório temporário gravável fora do repositório.
- Confirmar que o cache de módulos é gravável e está pré-aquecido quando a
  execução for offline.
- Registrar se a execução tem acesso à rede para checks de segurança.

## Execution checklist

1. Criar um diretório temporário para caches, JSON e log bruto.
2. Definir `GOCACHE` e `GOMODCACHE` para caminhos graváveis desse diretório.
3. Definir `GOFLAGS=-mod=readonly -buildvcs=false` para impedir alterações
   implícitas em dependências ou metadados de build.
4. Executar `pose validate --strict --report --json <arquivo>` e preservar
   stdout, stderr, exit code e o JSON produzido.
5. Executar os checks focados da spec no mesmo ambiente e anexar seus resultados.
6. Classificar separadamente falha de check, timeout, cache não gravável e
   indisponibilidade de rede/DNS; não converter falha em aprovação manual.
7. Referenciar os logs e o JSON na spec ou na revisão, incluindo o commit
   avaliado e qualquer limitação residual.

## Outputs

- Log bruto e JSON da validação preservados fora do repositório ou anexados à
  evidência da revisão.
- Resultado registrado como `pass`, `fail` ou `environment-limited` com causa.
- Revisão bloqueada quando um check obrigatório falhar sem uma mitigação
  reproduzível.

## Recurrence pilot

Medir `validate-native` por 45 dias após cada intervenção registrada com
`pose recurrence-effect --register`. Não promover um check opcional ou alterar
a matriz antes de observar pelo menos quatro semanas com taxa de sucesso de
95% e evidência completa.

## References

- [review workflow](review.md)
- [recurrence escalation](recurrence-escalation.md)
