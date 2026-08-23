<!--
Modelo para resumir uma sessão ao final (após session_finish). Preencha
só com dados que vieram de uma tool MCP real — nunca invente números.
Remova qualquer seção sem dado real em vez de deixar um placeholder.
-->

## Resumo da sessão

**Desafio:** {título do catalog_get} · **Modo:** {mode} · **Profundidade:** {depth}

### O que foi praticado

- {passos concluídos, de step_advance/step_complete ao longo da sessão}

### Evidência e avaliação

- {resumo de criteria/verdict de cada step_evaluate com submission_intent: true}
- {achados blocking/important_non_blocking relevantes, sem repetir cada advisory}

### Domínio (progress_get, se consultado)

- {dimensão → estado, só para as competências tocadas nesta sessão}

### Próxima revisão (review_due, se consultado)

- {competência e prazo, se houver alguma vencida ou próxima}

### Sugestão de próximo passo

- {learning_path_recommend, se o aluno pedir uma recomendação — sempre consultivo, nunca inicie sozinho}
