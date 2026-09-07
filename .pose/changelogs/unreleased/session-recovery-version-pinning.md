---
spec: session-recovery-version-pinning
category: fixed
breaking: false
refs:
---

Retome sessões após reiniciar o servidor, preservando conteúdo, políticas,
progresso, baselines e repetições de requests já confirmados. Atualizações do
catálogo não alteram sessões existentes; históricos sem dados suficientes
recebem um diagnóstico explícito de recuperação indisponível.
