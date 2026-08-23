// Package assessment holds the pure, catalog-backed decisions behind
// feedback preparation and criterion resolution (PROJECT.md §15.9, §21).
// Like internal/assistance, it never mutates session state or writes to
// the event store: internal/session owns evaluation, feedback, reflection,
// completion and advance under its own lock, calling into this package for
// the decisions themselves.
package assessment

import "github.com/oseiaspereira88/codinho/internal/curriculum"

// knownRubrics are the fixed rubric resource identifiers PROJECT.md §15.11
// documents. Resources are optional in the MCP protocol and no required V1
// flow depends on them (§15.11); feedback_prepare only *names* them so the
// tutor skill can pull the ones it needs from its own reference material
// (.agents/skills/codinho/references/feedback-rubric.md). This server does
// not author feedback text itself (non-goal).
var knownRubrics = []string{"rubric://idiomatic-go", "rubric://technical-communication"}

// FeedbackPacket is the read-only context feedback_prepare hands the tutor
// agent so it authors the actual feedback (requirement R1).
type FeedbackPacket struct {
	StepID     string
	Objective  string
	Scope      string
	Question   string
	RubricRefs []string
}

// PrepareFeedback assembles the packet from step's authored instruction and
// the learner's own question, if any.
func PrepareFeedback(step curriculum.StepAuthoring, question string) FeedbackPacket {
	return FeedbackPacket{
		StepID:     step.ID,
		Objective:  step.Instruction.Objective,
		Scope:      step.Instruction.Scope,
		Question:   question,
		RubricRefs: append([]string{}, knownRubrics...),
	}
}
