<!--
Modelo para o relatório de uma sessão mode:interview, após interview_report.
Preencha só com dados que vieram da tool real — nunca invente um número ou
verdict. Remova qualquer seção sem dado real em vez de deixar um placeholder.
Nunca reduza o resultado a um score único (non-functional requirement).
-->

## Relatório da simulação de entrevista

**Desafio:** {challenge_id, do interview_report} · **Duração:** {elapsed_seconds
formatado como minutos} · **Encerramento:** {finish_reason: explicit ou timeout}

### Avaliação por critério (nunca um score único)

- {para cada evaluations[].criteria[]: nome, kind, verdict}

### Lacunas (gaps)

- {cada entrada de gaps, no formato step_id/nome do critério}

### Uso de auxílio

- Pistas concedidas: {hints.Granted} · Pistas bloqueadas pela política: {hints.Blocked}

### Reflexões registradas

- {para cada reflections[]: prompt e answer}

### Recomendações para praticar depois

- {recommended, se fornecido — ex.: saída de learning_path_recommend}

### Aviso de integridade

{integrity_note, repetido na íntegra — nunca resumido ou omitido}
