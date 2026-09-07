package mcpserver

// coreInstructions is the normative content of PROJECT.md §15.2, verbatim.
// It must remain within the first 512 characters of Instructions
// (requirement R2) so a host that truncates long instructions never drops
// an invariant.
const coreInstructions = `Servidor de aprendizagem assistida. O aluno é o único autor do código durante tutoria. Revele apenas o conteúdo autorizado pela sessão. Feedback não altera progresso; avaliação não conclui; conclusão não avança. Não execute comandos livres nem observe caminhos fora do workspace autorizado. Informe sempre o efeito da operação sobre o progresso.`

// Instructions is published to clients as the server's global instructions
// (PROJECT.md §15.1, §15.2).
var Instructions = coreInstructions + ` Antes de step_evaluate, cite evidências da mesma sessão e passo. Para checks estruturais informe check_id. Registre observações qualitativas com evidence_record (source, text, rubric_ref) antes de citá-las; o registro não comprova execução estrutural. Evidência obsoleta exige nova observação/check.`
