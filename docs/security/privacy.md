# Privacidade e dados locais

## O que é coletado

Nada sai da máquina local por padrão (Constraint: "zero telemetria e
zero upload por padrão"). O único estado que o `codinho` grava é local,
em `<workspace>/.codinho/state/`:

- `events.jsonl` — o log de eventos de sessão/domínio/progresso
  (`internal/eventstore`), append-only.
- `evidence/` — blobs de evidência endereçados por conteúdo
  (`internal/evidence`), imutáveis (permissão `0o400` após escrita).

Nenhum dos dois é enviado a lugar nenhum pelo `codinho` em si; o que o
host (Codex CLI, extensão de IDE) faz com a conversa é responsabilidade
do host, fora deste escopo.

## O que nunca é coletado

- Arquivos excluídos por padrão (`.env`, `id_rsa`, `credentials.*`,
  `.aws`, `.netrc`, `.npmrc`, chaves `.pem`/`.key`/`.pfx`/`.p12`) —
  ver `internal/workspace/redaction.go`.
- Qualquer trecho que combine com um padrão de segredo conhecido
  (API key, senha, token, chave privada PEM, token do GitHub) — redigido
  antes de virar evidência, mesmo dentro de um arquivo legítimo.
- Áudio, vídeo, dados pessoais ou atividade fora do workspace declarado
  (nunca implementado; não há microfone/câmera/screen-capture em nenhum
  código deste repositório).

## Retenção

O log de eventos é intencionalmente **append-only e imutável** — é o
que permite replay determinístico de progresso/mastery
(mastery-review-scheduling). Isso significa que a retenção não é uma
expiração automática por item (apagar um evento no meio quebraria o
replay); a retenção é **explícita e no nível do workspace inteiro**:

- `codinho privacy export --dest <path>` copia o estado local
  (eventos + evidências) para onde o aluno quiser, sem tocar o original.
- `codinho privacy purge --confirm` remove `.codinho/state/` por
  completo — nunca o código do aluno, nunca `.git`, nunca nada fora
  dessa pasta. Sem `--confirm`, o comando recusa (fail-closed).

## Export e remoção (requirement R6)

```sh
codinho privacy export --dest ~/backups/codinho-2026-08-23
codinho privacy purge --confirm
```

`purge` é irreversível para o histórico de progresso local — não há
"lixeira". Exporte antes se quiser manter uma cópia.
