# ADR: Durable session replay with pinned content

## Status
Accepted — 2026-09-07 UTC.

## Context
O log JSONL sobrevive ao processo, mas session.New recria mapas/contador vazios.
Eventos legados session_started só têm challenge_id/mode e não permitem
reconstruir consentimentos ou versão de conteúdo com segurança.
knowledge:planning-audit-2026-09 reproduziu perda de sessão e colisão de IDs.

## Decision
- Preservar o envelope de evento schema_version 1 e estender session_started
  com recovery_version 1, entrada original, política completa e cópia JSON
  do desafio resolvido com digest SHA-256. A cópia serve somente à sessão;
  o YAML continua fonte de autoria. Limitar o tamanho ao orçamento do log.
- Reaplicar eventos confirmados por um redutor determinístico. Validar a
  mesma transição sobre cópia do estado antes de anexar um novo evento.
  Não depender de snapshots externos gravados depois de responder.
- Recuperar o contador de IDs e cache de início pelo histórico, incluindo
  streams legadas. Reutilização incompatível de request_id é conflito.
- Persistir digest da entrada das mutações; reconstruir a resposta no estado
  da revisão confirmada para retries históricos. Validar tipo/payload dentro
  do lock do store, incluindo concorrência entre serviços.
- Tratar avaliação com submission_intent como intenção durável; correlacionar
  a tentativa secundária pelo ID do evento e completar append interrompido
  somente após validar toda a projeção. Não modificar história incompatível.
- Declarar sessão legada/incompatível irrecuperável com erro próprio,
  preservando o log e permitindo sessões novas. Não inferir política
  ausente a partir dos defaults ou do catálogo atual.
- Persistir baseline, root real, globs e referência de evidência em eventos
  de observação. Revalidar o root recuperado antes de usá-lo; não seguir uma
  substituição por symlink nem reutilizar baseline em outro escopo.
- Corrigir append após cauda JSONL incompleta preservando os bytes rejeitados
  em arquivo de recuperação privado. Falha ou corrupção intermediária não
  permite append silencioso.
- Separar OpenReadOnly para inspeção administrativa sem lock: somente o
  writer sob lock pode reparar cauda. Manter resolução manual de lock órfão.

## Alternatives
- Snapshot posterior ao append: rejeitado como autoridade porque deixa uma
  janela de perda entre confirmação e snapshot.
- Reconsultar o pack atual por ID: rejeitado porque troca checks/instruções
  sem consentimento após update.
- Migrar para SQLite: rejeitado; não resolve a ausência de payload de sessão
  e contraria o ADR local-versioned-catalog-and-event-state sem medição.

## Consequences
O log cresce uma vez por sessão pelo desafio fixado. Eventos antigos seguem
legíveis para auditoria, mas nem todos podem voltar a ser sessões ativas.
Conteúdo recuperado nunca aparece como gabarito público; as projeções
existentes continuam aplicando disclosure. Retenção/exclusão do estado inclui
os backups de cauda truncada.

## Review trigger
Reavaliar se volume medido comprometer startup ou se o schema de sessão
precisar de migração. Um novo payload exige teste de replay e compatibilidade.

## References
- [Spec](../specs/2026-09-07-session-recovery-version-pinning.md).
- [Armazenamento local](2026-08-22-local-versioned-catalog-and-event-state.md).
