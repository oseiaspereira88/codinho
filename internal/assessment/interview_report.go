package assessment

import (
	"time"

	"github.com/oseiaspereira88/codinho/internal/learning"
)

// InterviewEvaluation is one step's recorded evaluation, decoded by the
// caller from an evaluation_recorded event. It is never re-derived here:
// Resolve (sealed) already produced these verdicts.
type InterviewEvaluation struct {
	StepID   string
	Criteria []learning.CriterionResult
}

// InterviewHintUsage summarizes hint activity across an interview
// session, decoded by the caller from hint_requested events (requirement
// R4: "bloquear ou registrar pedidos de pista").
type InterviewHintUsage struct {
	Granted int
	Blocked int
}

// InterviewReflection is one reflection prompt/answer pair recorded
// during the session, decoded from a reflection_recorded event.
type InterviewReflection struct {
	StepID string
	Prompt string
	Answer string
}

// InterviewReport is the descriptive, non-scalar summary produced at the
// end of an interview session (PROJECT.md §9.5: "avalia código,
// raciocínio, testes, complexidade e comunicação ao final"; requirement
// R6, R7; non-functional requirement: "Relatório não deve reduzir
// desempenho a um score único").
type InterviewReport struct {
	ChallengeID   string
	Elapsed       time.Duration
	TimedOut      bool
	FinishReason  string
	Evaluations   []InterviewEvaluation
	Hints         InterviewHintUsage
	Reflections   []InterviewReflection
	Gaps          []string
	Recommended   []string
	IntegrityNote string
}

// integrityNote is PROJECT.md §22's honest disclosure of this
// simulation's limits (requirement R8): no proctoring, no vendor
// affiliation, no credential.
const integrityNote = "Simulação local sem vigilância, gravação ou vínculo com processo seletivo real; o resultado é evidência pedagógica, não uma credencial ou aprovação."

// BuildInterviewReport assembles a descriptive report from already-
// recorded evidence. It never re-derives a verdict and never collapses
// the evaluations into a single score: gaps lists every non-met
// criterion by step, for the reader to judge in context.
func BuildInterviewReport(
	challengeID string,
	elapsed time.Duration,
	timedOut bool,
	finishReason string,
	evaluations []InterviewEvaluation,
	hints InterviewHintUsage,
	reflections []InterviewReflection,
	recommended []string,
) InterviewReport {
	var gaps []string
	for _, e := range evaluations {
		for _, c := range e.Criteria {
			if c.Verdict != learning.VerdictMet {
				gaps = append(gaps, e.StepID+"/"+c.Name)
			}
		}
	}
	return InterviewReport{
		ChallengeID:   challengeID,
		Elapsed:       elapsed,
		TimedOut:      timedOut,
		FinishReason:  finishReason,
		Evaluations:   evaluations,
		Hints:         hints,
		Reflections:   reflections,
		Gaps:          gaps,
		Recommended:   recommended,
		IntegrityNote: integrityNote,
	}
}
