---
spec: eventstore-idempotency-scope
category: fixed
breaking: false
refs:
---

Corrigido um bug de idempotência no event store: uma retentativa com o
mesmo `request_id` em duas sessões diferentes podia devolver o evento de
uma sessão para a outra, silenciosamente. Agora o `request_id` é único por
sessão, não globalmente — retentativas continuam idempotentes dentro da
mesma sessão, e sessões diferentes nunca mais colidem.
