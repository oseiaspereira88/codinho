---
spec: session-tree-progression
category: fixed
breaking: true
refs:
---

Sessões iniciam na profundidade solicitada e percorrem todas as layers,
preservando progresso ao mudar granularidade e reiniciar o servidor. Declare
children_mode: choice para alternativas e use next_step_id em step_advance
para selecionar um ramo permitido. Irmãos sem choice passam a seguir sequência
obrigatória; concluir uma janela ampla cobre sua subárvore sem duplicar eventos.
