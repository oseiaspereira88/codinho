package assessment

import (
	"testing"

	"github.com/oseiaspereira88/codinho/internal/curriculum"
)

func TestPrepareFeedbackAssemblesPacketFromStep(t *testing.T) {
	step := curriculum.StepAuthoring{
		ID: "model.declare-user-struct",
		Instruction: curriculum.InstructionAuthoring{
			Objective: "Declarar o tipo User.",
			Scope:     "Somente a declaração.",
		},
	}
	packet := PrepareFeedback(step, "why exported fields?")
	if packet.StepID != step.ID || packet.Objective != step.Instruction.Objective || packet.Scope != step.Instruction.Scope {
		t.Fatalf("unexpected packet: %+v", packet)
	}
	if packet.Question != "why exported fields?" {
		t.Fatalf("question = %q, want the caller's question", packet.Question)
	}
	if len(packet.RubricRefs) == 0 {
		t.Fatal("expected at least one rubric reference")
	}
}

func TestPrepareFeedbackRubricRefsAreIndependentCopies(t *testing.T) {
	first := PrepareFeedback(curriculum.StepAuthoring{}, "")
	first.RubricRefs[0] = "mutated"
	second := PrepareFeedback(curriculum.StepAuthoring{}, "")
	if second.RubricRefs[0] == "mutated" {
		t.Fatal("PrepareFeedback must not share a mutable backing array across calls")
	}
}
