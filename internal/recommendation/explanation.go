package recommendation

import (
	"fmt"
	"strings"
)

// explain builds the structured justification requirement R5 requires:
// evidence, gap, dependency and estimated cost — never a bare score.
// remediationFor is non-empty exactly when c is a remediation step ahead
// of that challenge ID (requirement R6, Decision 4).
func explain(c Candidate, m Mastery, completed map[string]bool, remediationFor string) Explanation {
	evidence := "sem evidência de domínio registrada ainda"
	if m.BestState != "" && m.BestState != "not_observed" {
		evidence = fmt.Sprintf("melhor estado de domínio observado: %s", m.BestState)
	}
	if m.ReviewOverdueBy > 0 {
		evidence += fmt.Sprintf("; revisão vencida há %d dia(s)", m.ReviewOverdueBy)
	}

	gap := "nenhum gap conhecido"
	if stateRank[m.BestState] < stateRank[remediationThreshold] {
		gap = "ainda não demonstra sem ajuda para a competência-alvo"
	}

	dependency := "sem pré-requisitos pendentes"
	if len(c.Prerequisites) > 0 {
		var unmet []string
		for _, p := range c.Prerequisites {
			if !completed[p] {
				unmet = append(unmet, p)
			}
		}
		if len(unmet) > 0 {
			dependency = "pré-requisitos pendentes: " + strings.Join(unmet, ", ")
		} else {
			dependency = "pré-requisitos satisfeitos: " + strings.Join(c.Prerequisites, ", ")
		}
	}

	return Explanation{
		Evidence: evidence, Gap: gap, Dependency: dependency,
		CostMinutes: c.EstimatedMinutes, RemediationFor: remediationFor,
	}
}
