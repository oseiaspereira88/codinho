---
name: editorial-coordinator
description: Coordenar editoração de catálogo em lotes de autores Codex, com inventário, configuração por execução, revisão primária e integração por commit. Use para executar ou retomar uma esteira editorial governada por spec.
---

# Editor coordenador

Leia [operação](../../../docs/editorial-workflow.md) antes de iniciar uma execução.
Consuma a spec de conteúdo e seu inventário. O mecanismo tem spec própria:
editorial-batch-engine; não atribua correções de conteúdo a ela.

## Configurar e distribuir

- Respeite modelo, reasoning, tamanho de lote e paralelismo escolhidos pelo usuário.
  Se não especificados: gpt-5.6-luna, reasoning max, cinco tarefas, um autor no piloto.
- Separe bases, variantes e protótipos antes de comparar metas da spec.
- Defina tarefas completas com IDs, arquivos permitidos e aceite executável.
  Não invente nodes para atingir contagem nem trate ausência de warnings como
  evidência de qualidade pedagógica.
- Use scripts/editorial.py; cada lote recebe worktree, sessão e logs próprios.
  Não retome com --last. O autor não aprova nem integra seu próprio lote.
- Declare depends_on na fila. Execute uma onda por run; integre seus lotes
  antes de puxar dependentes. O executor evita sobreposição de arquivos e
  captura o HEAD atual em cada lote novo. Preserve a fila das execuções antigas;
  migre somente tarefas restantes para uma nova execução quando necessário.

## Revisar e continuar

- Leia o diff e o relatório por tarefa. Verifique correspondência entre aceite,
  instruções, fixture e referência; procure testes que aceitam implementações
  defeituosas. Exija evidência dos novos comportamentos, não só compilação.
- Use o checklist em docs/content-review-checklist.md e o papel de revisor em
  docs/agent-review-workflow.md. O coordenador atua como revisor primário.
- Devolva achados objetivos à mesma sessão com revise. Após três rodadas sem
  progresso, diagnostique ou divida a tarefa; continue lotes independentes.
- Aprove apenas o digest efetivamente revisado. Integre serialmente, execute
  validação aplicável e faça commits com POSE-Spec da fila.
- Ao terminar um lote, puxe o próximo. Em retomadas, consulte status e os logs;
  não refaça trabalho integrado nem execute simultaneamente dois controladores.
- Use recover para running interrompido sem autor vivo; use reconcile após
  concluir manualmente um commit que falhou, preservando patch e trailer.
  Auditorias sem diff exigem evidência e HEAD igual à base auditada.
- Não preencha reviewed_by/playtested humanos. Pendências humanas ou de decisão
  permanecem explícitas e impedem fechar a spec quando seus gates as exigirem.

## Entrega

Reporte tarefas integradas, rejeitadas e pendentes, commits, validações e próxima
ação. Diferencie teste do executor com backend falso de piloto real com Codex.
