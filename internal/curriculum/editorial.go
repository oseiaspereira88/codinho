package curriculum

import (
	"fmt"
	"sort"
	"strings"
)

// EditorialSeverity distinguishes objective, mechanically-enforceable
// findings from qualitative ones only a human can judge
// (catalog-authoring-quality Decision 1: "bloquear regras objetivas e
// exigir evidência humana para qualidade semântica").
type EditorialSeverity string

const (
	SeverityBlocking EditorialSeverity = "blocking"
	SeverityWarning  EditorialSeverity = "warning"
)

// EditorialRuleID names one catalog-authoring-quality rule, stable
// across runs so CI can pin or suppress a specific rule without changing
// results of every other rule silently (Compatibility requirement).
type EditorialRuleID string

const (
	RuleHintsIncreasingOrder     EditorialRuleID = "hints_increasing_order"
	RuleHintRevealsSolutionEarly EditorialRuleID = "hint_reveals_solution_early"
	RuleBriefContainsCodeFence   EditorialRuleID = "brief_contains_code_fence"
	RuleCompoundMicroInstruction EditorialRuleID = "compound_micro_instruction"
	RuleMissingCompetency        EditorialRuleID = "missing_competency"
	RuleMissingAcceptance        EditorialRuleID = "missing_acceptance"
	RuleMissingHintsForGatedStep EditorialRuleID = "missing_hints_for_gated_step"
	RuleMissingReflection        EditorialRuleID = "missing_reflection"
	RuleRelationIsolated         EditorialRuleID = "relation_isolated"
	RuleRelationNotReciprocal    EditorialRuleID = "relation_not_reciprocal"
	// RuleFixtureNotReproducible, RuleCheckNotResolvable and
	// RuleCheckExecutionError are reported by the CLI layer (not this
	// package: fixture materialization and execution belong to the
	// administrative command's explicit code-execution boundary) after actually
	// materializing a challenge's fixture and running its declared
	// checks against it (requirement R5).
	RuleFixtureNotReproducible EditorialRuleID = "fixture_not_reproducible"
	RuleCheckNotResolvable     EditorialRuleID = "check_not_resolvable"
	RuleCheckExecutionError    EditorialRuleID = "check_execution_error"
)

// EditorialFinding is one catalog-authoring-quality finding (requirement
// R2, R3, R4, R6, R7; non-functional: "Mensagens devem apontar pack,
// arquivo, item, regra e correção possível").
type EditorialFinding struct {
	File       string
	Item       string
	Rule       EditorialRuleID
	Severity   EditorialSeverity
	Detail     string
	Suggestion string
}

// RunEditorialChecks runs every catalog-authoring-quality rule over
// packs, in addition to (never instead of) Validate's structural checks:
// it never re-implements schema, ID, reference or cycle validation
// (requirement R1 stays Validate's alone).
func RunEditorialChecks(packs []Pack) []EditorialFinding {
	var findings []EditorialFinding
	for _, p := range packs {
		for _, ch := range p.Challenges {
			findings = append(findings, checkHintOrder(p.File, ch.ID, ch.Layers)...)
			findings = append(findings, checkNoEmbeddedSolutionText(p.File, ch)...)
			findings = append(findings, checkCompoundIntent(p.File, ch.ID, ch.Layers)...)
			findings = append(findings, checkCoverageGaps(p.File, ch)...)
			findings = append(findings, checkPublicationMetadata(p.File, ch)...)
		}
	}
	findings = append(findings, checkRelationCoherence(packs)...)
	return findings
}

func walkLayers(layers []LayerAuthoring, fn func(StepAuthoring)) {
	var walk func(steps []StepAuthoring)
	walk = func(steps []StepAuthoring) {
		for _, s := range steps {
			fn(s)
			walk(s.Children)
		}
	}
	for _, l := range layers {
		walk(l.MacroSteps)
	}
}

// checkHintOrder enforces PROJECT.md §8.4's ladder: each step's hints
// must strictly increase in level, and nothing below level 6 may be a
// solution-kind hint (requirement R4).
func checkHintOrder(file, challengeID string, layers []LayerAuthoring) []EditorialFinding {
	var out []EditorialFinding
	walkLayers(layers, func(s StepAuthoring) {
		prev := -1
		for _, h := range s.Hints {
			if h.Level <= prev {
				out = append(out, EditorialFinding{
					File: file, Item: challengeID + "/" + s.ID, Rule: RuleHintsIncreasingOrder, Severity: SeverityBlocking,
					Detail:     fmt.Sprintf("hint level %d does not increase after %d", h.Level, prev),
					Suggestion: "ordene os hints por level estritamente crescente, sem repetir nível",
				})
			}
			prev = h.Level
			if h.Level < 6 && h.Kind == "solution" {
				out = append(out, EditorialFinding{
					File: file, Item: challengeID + "/" + s.ID, Rule: RuleHintRevealsSolutionEarly, Severity: SeverityBlocking,
					Detail:     fmt.Sprintf("hint de nível %d tem kind=solution", h.Level),
					Suggestion: "gabarito (kind=solution) só é permitido no nível 6",
				})
			}
		}
	})
	return out
}

// checkNoEmbeddedSolutionText flags a code fence in brief or acceptance
// text, a common way a gabarito leaks before any hint is even requested
// (requirement R4).
func checkNoEmbeddedSolutionText(file string, ch ChallengeAuthoring) []EditorialFinding {
	var out []EditorialFinding
	if strings.Contains(ch.Brief, "```") {
		out = append(out, EditorialFinding{
			File: file, Item: ch.ID, Rule: RuleBriefContainsCodeFence, Severity: SeverityBlocking,
			Detail: "brief contém um bloco de código", Suggestion: "código-fonte pertence a fixture ou solução, nunca ao brief",
		})
	}
	for _, a := range ch.Acceptance {
		if strings.Contains(a, "```") {
			out = append(out, EditorialFinding{
				File: file, Item: ch.ID, Rule: RuleBriefContainsCodeFence, Severity: SeverityBlocking,
				Detail: "um critério de acceptance contém um bloco de código", Suggestion: "critérios de aceite descrevem comportamento observável, não código",
			})
		}
	}
	return out
}

// checkCompoundIntent flags a micro step whose objective joins more than
// one independent clause — a conservative heuristic for the single-
// intent test (requirement R2). It never blocks: compound phrasing is
// often a false positive, so this is always advisory.
func checkCompoundIntent(file, challengeID string, layers []LayerAuthoring) []EditorialFinding {
	var out []EditorialFinding
	walkLayers(layers, func(s StepAuthoring) {
		if s.Kind != "micro" {
			return
		}
		if hasCompoundIntent(s.Instruction.Objective) {
			out = append(out, EditorialFinding{
				File: file, Item: challengeID + "/" + s.ID, Rule: RuleCompoundMicroInstruction, Severity: SeverityWarning,
				Detail:     s.Instruction.Objective,
				Suggestion: "divida em dois micropassos — um verbo independente por passo",
			})
		}
	})
	return out
}

var compoundIntentMarkers = []string{" e depois ", " e então ", " e também ", " e ainda ", "; "}

func hasCompoundIntent(objective string) bool {
	lower := strings.ToLower(objective)
	for _, marker := range compoundIntentMarkers {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

// checkCoverageGaps flags a challenge missing a primary competency or
// any acceptance criterion, a gated micro step with no authored hint,
// and a challenge with no reflection prompt anywhere in its tree
// (requirement R3). Only the first two are structural enough to block;
// the rest need a human's judgment of whether reflection was actually
// warranted.
func checkCoverageGaps(file string, ch ChallengeAuthoring) []EditorialFinding {
	var out []EditorialFinding
	if len(ch.Competencies.Primary) == 0 {
		out = append(out, EditorialFinding{
			File: file, Item: ch.ID, Rule: RuleMissingCompetency, Severity: SeverityBlocking,
			Detail: "nenhuma competência primária declarada", Suggestion: "declare ao menos um id em competencies.primary",
		})
	}
	if len(ch.Acceptance) == 0 {
		out = append(out, EditorialFinding{
			File: file, Item: ch.ID, Rule: RuleMissingAcceptance, Severity: SeverityBlocking,
			Detail: "nenhum critério de aceite declarado", Suggestion: "declare ao menos um item em acceptance",
		})
	}
	hasReflection := false
	walkLayers(ch.Layers, func(s StepAuthoring) {
		if len(s.Reflection.Optional) > 0 {
			hasReflection = true
		}
		if s.Kind == "micro" && s.Completion.RequiresPositiveEvaluation && len(s.Hints) == 0 {
			out = append(out, EditorialFinding{
				File: file, Item: ch.ID + "/" + s.ID, Rule: RuleMissingHintsForGatedStep, Severity: SeverityWarning,
				Detail: "passo exige avaliação positiva mas não declara nenhum hint", Suggestion: "declare ao menos um hint de nível 1",
			})
		}
	})
	if !hasReflection {
		out = append(out, EditorialFinding{
			File: file, Item: ch.ID, Rule: RuleMissingReflection, Severity: SeverityWarning,
			Detail: "nenhuma pergunta de reflexão em toda a árvore do desafio", Suggestion: "adicione reflection.optional em ao menos um passo relevante",
		})
	}
	return out
}

// checkRelationCoherence flags an item with no relation beyond
// requires/prerequisites, and a contrasts_with/commonly_fails_with edge
// with no reciprocal counterpart (requirement R6). Both are advisory: a
// missing reciprocal edge can be intentionally one-directional and
// justified outside the schema, which is exactly the human judgment call
// Decision 1 reserves.
func checkRelationCoherence(packs []Pack) []EditorialFinding {
	type edgeKey struct {
		from, to string
		kind     RelationKind
	}

	itemFile := map[string]string{}
	hasNonPrecedenceRelation := map[string]bool{}
	seen := map[edgeKey]bool{}

	for _, p := range packs {
		for _, c := range p.Concepts {
			itemFile[c.ID] = p.File
		}
		for _, c := range p.Competencies {
			itemFile[c.ID] = p.File
		}
		for _, ch := range p.Challenges {
			itemFile[ch.ID] = p.File
		}
		for _, r := range p.Relations {
			kind := RelationKind(r.Kind)
			seen[edgeKey{r.From, r.To, kind}] = true
			if kind != RelationRequires {
				hasNonPrecedenceRelation[r.From] = true
				hasNonPrecedenceRelation[r.To] = true
			}
		}
	}

	var out []EditorialFinding
	ids := make([]string, 0, len(itemFile))
	for id := range itemFile {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if !hasNonPrecedenceRelation[id] {
			out = append(out, EditorialFinding{
				File: itemFile[id], Item: id, Rule: RuleRelationIsolated, Severity: SeverityWarning,
				Detail:     "item sem nenhuma relation além de requires/prerequisites",
				Suggestion: "declare ao menos uma relation (relates_to, applies_in, deepens_into, contrasts_with, commonly_fails_with, evidences)",
			})
		}
	}

	for _, p := range packs {
		for _, r := range p.Relations {
			kind := RelationKind(r.Kind)
			if kind != RelationContrastsWith && kind != RelationCommonlyFailsWith {
				continue
			}
			if !seen[edgeKey{r.To, r.From, kind}] {
				out = append(out, EditorialFinding{
					File: p.File, Item: r.From + "->" + r.To, Rule: RuleRelationNotReciprocal, Severity: SeverityWarning,
					Detail:     fmt.Sprintf("%s de %s para %s não tem a recíproca %s de %s para %s", r.Kind, r.From, r.To, r.Kind, r.To, r.From),
					Suggestion: "adicione a relação recíproca, ou documente a justificativa da assimetria na revisão",
				})
			}
		}
	}
	return out
}
