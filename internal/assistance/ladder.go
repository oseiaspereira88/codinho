// Package assistance holds the pure, catalog-backed decisions behind the
// assistance ladder, syntax recall and pedagogical detours (PROJECT.md
// §8.4, §8.8). It never mutates session state or writes to the event
// store: internal/session owns the ladder ratchet, hint events and detour
// lifecycle under its own lock, calling into this package for the
// decisions themselves (non-goal: "redigir explicações abertas dentro do
// MCP" — this package classifies content, it never authors prose).
package assistance

import (
	"errors"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

// KindSyntaxRecall is the authored HintAuthoring.Kind value PROJECT.md
// §14.9 uses to mark a rung as a syntax reminder. Syntax recall need not
// consume a hint the way other rungs do (RF-022; PROJECT.md §8.4 "custo
// menor no modo ensino").
const KindSyntaxRecall = "syntax_recall"

// defaultKind is the generic content classification for a disclosure
// level per PROJECT.md §8.4's table, used when the step's authoring does
// not override it with a more specific HintAuthoring.Kind.
var defaultKind = map[learning.DisclosureLevel]string{
	learning.DisclosureGuidingQuestion: "guiding_question",
	learning.DisclosureConceptOrAPI:    "concept_or_api",
	learning.DisclosureLogicalOutline:  "logical_outline",
	learning.DisclosurePseudocode:      "pseudocode",
	learning.DisclosureSkeleton:        "skeleton",
	learning.DisclosureSolution:        "solution",
}

// ErrSolutionRequiresConfirmation is returned when the next rung would
// reveal the solution but the caller did not explicitly confirm it
// (PROJECT.md §8.4: "nível 6 exige pedido explícito e confirmação do
// efeito pedagógico").
var ErrSolutionRequiresConfirmation = errors.New("assistance: revealing the solution requires explicit confirmation")

// ErrHelpDisabled is returned when the session's help policy forbids
// hints entirely (HelpNoHints).
var ErrHelpDisabled = errors.New("assistance: session help policy disallows hints")

// ErrNoHintAuthored is returned when a step has no authored syntax-recall
// rung for SyntaxRecallLevel to find.
var ErrNoHintAuthored = errors.New("assistance: no syntax-recall hint authored for this step")

// KindFor resolves the content classification for level on step: the
// step's authored override when present, otherwise the level's generic
// default from PROJECT.md §8.4's table.
func KindFor(step curriculum.StepAuthoring, level learning.DisclosureLevel) string {
	for _, h := range step.Hints {
		if learning.DisclosureLevel(h.Level) == level && h.Kind != "" {
			return h.Kind
		}
	}
	return defaultKind[level]
}

// NextHint decides the next disclosure level hint_request would grant for
// a step currently at current, without mutating anything: it climbs
// exactly one rung per call by construction (PROJECT.md §8.4). The
// session's disclosure cap and hint-level bookkeeping are enforced
// separately by learning.StepProgress.GrantHint.
func NextHint(step curriculum.StepAuthoring, current learning.DisclosureLevel, help learning.HelpPolicyKind, confirmSolution bool) (learning.DisclosureLevel, string, error) {
	if help == learning.HelpNoHints {
		return 0, "", ErrHelpDisabled
	}
	next := current + 1
	if next == learning.DisclosureSolution && !confirmSolution {
		return 0, "", ErrSolutionRequiresConfirmation
	}
	return next, KindFor(step, next), nil
}

// SyntaxRecallLevel finds the step's authored syntax-recall rung, if any
// (RF-022): a targeted lookup independent of the caller's current ladder
// position.
func SyntaxRecallLevel(step curriculum.StepAuthoring) (learning.DisclosureLevel, bool) {
	for _, h := range step.Hints {
		if h.Kind == KindSyntaxRecall {
			return learning.DisclosureLevel(h.Level), true
		}
	}
	return 0, false
}
