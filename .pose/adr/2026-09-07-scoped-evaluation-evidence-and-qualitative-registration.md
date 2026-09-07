# ADR: Scoped evaluation evidence and qualitative registration

## Status
Accepted — 2026-09-07.

## Context
EvidenceGet restringe leitura, mas StepEvaluate consultava blobs globalmente.
Uma avaliação podia aprovar com evidência de outra sessão ou workspace alterado.
knowledge:adr-durable-session-replay-review registrou a lacuna. Sessão, aplicação,
MCP e executor compartilham um log; preservamos essas fronteiras e o envelope v1.

## Decision
- Injetar no núcleo uma interface de validação implementada na aplicação. Sem
  adaptador, qualquer citação nova falha fechada. Retry confirmado ocorre antes
  dessa validação, preservando a resposta histórica sem consultar arquivos atuais.
- Autorizar somente eventos confirmados do mesmo session_id e step_id. O stream
  fixa desafio/conteúdo; exigir que o check citado exista no desafio fixado e
  coincida com check_id explícito no critério. Conferir integridade do blob.
- Ler root/globs do evento de observação anterior ao produtor. Revalidar root
  canônica e fingerprint; ausência de escopo recuperável, blob ou root rejeita
  consumo. Repetir validação imediatamente antes do append, sob mutex da sessão;
  expected_revision detecta eventos concorrentes de workspace/check.
- Para checks novos, comparar fingerprint antes/depois da execução e rejeitar
  drift. Globs delimitam o escopo declarado; os registros não provam arquivos
  externos a esse escopo. Escritores externos não respeitam mutex: a garantia
  é amostragem de arquivos, não snapshot atômico nem detecção de alterações ABA.
- Derivar verdict estrutural somente do resultado de check autorizado. Baseline,
  diff ou nota qualitativa não demonstram uma afirmação estrutural arbitrária;
  retornar unverifiable. Sem ID estrutural, manter unverifiable.
- Adicionar evidence_record: source enum learner_explanation/tutor_observation/
  external_artifact, texto observado e rubric_ref obrigatórios. Persistir blob
  redigido com sessão, passo, desafio, origem e rubrica; não buscar URL ou executar
  conteúdo externo. O registro atesta que o tutor apresentou o texto, não sua
  veracidade. Só julgamento qualitativo pode consumi-lo com a mesma rubrica.
  Observações/checks autorizados também podem fundamentar julgamento qualitativo
  com rubrica; continuar verificando freshness nesses casos.
- Persistir evidence_lineage junto à avaliação para ligar critério ao evento,
  escopo e fingerprint amostrados. Não reescrever avaliações ou tentativas antigas.
  Adicionar campos opcionais omitempty à identidade JSON de requests para não
  invalidar retries que precedem o novo contrato.
- Preservar a política existente de seleção/cobertura de critérios pelo tutor;
  esta remediação não implementa cobertura completa da rubrica autorada.

## Alternatives
- Confiar em evidence_id: rejeitado; endereçamento por conteúdo não autoriza uso.
- Chamar EvidenceGet só no MCP: rejeitado; deixa bypass pela aplicação e ignora
  o nó/check e retries históricos.
- Buscar provas externas automaticamente: rejeitado; introduz rede/autoridade
  indevida e execução de conteúdo fora do escopo local.
- Bloquear filesystem externo: rejeitado sem snapshot/sandbox isolado; declarar
  a janela residual em vez de prometer atomicidade impossível com mutex local.

## Consequences
Clientes passam check_id e registram notas qualitativas antes de citá-las.
Erros estáveis não expõem paths/texto. Logs legados conservam replay; evidência
sem escopo recuperável não autoriza avaliação nova. Nenhum banco ou LLM novo.
Rever o contrato se for necessário avaliar snapshots imutáveis ou artefatos
externos autenticados. Validar a compatibilidade em processos MCP reais.

## References
- [Spec](../specs/2026-09-07-evaluation-evidence-lineage.md).
- [Fronteiras](2026-08-22-agent-mcp-and-core-boundaries.md).
- [Replay](2026-09-06-durable-session-replay-with-pinned-content.md).
