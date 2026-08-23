package session

import "github.com/oseiaspereira88/codinho/internal/learning"

// ModeDefaults holds the per-mode independent starting point every
// session dimension falls back to when session_start leaves it unset.
// PROJECT.md §8.3 is explicit that the seven session dimensions "não são
// inferidas umas das outras": a mode only supplies defaults, it never
// forces a dimension. An explicit StartInput field always overrides its
// mode's default (requirement R1, R2).
type ModeDefaults struct {
	Depth      learning.Depth
	Help       learning.HelpPolicyKind
	Disclosure learning.DisclosureLevel
	Evaluation learning.EvaluationPolicyKind
}

// DefaultsForMode returns mode's independent defaults, drawn from the
// behavior table in PROJECT.md §8.2. Any mode not explicitly listed
// (including an empty or unrecognized value) falls back to practice's
// defaults, matching Service.Start's own fallback for an unset Mode.
func DefaultsForMode(mode learning.PedagogicalMode) ModeDefaults {
	switch mode {
	case learning.ModeTeaching:
		// "Explicações contextuais, granularidade adaptável e pistas
		// progressivas": start wide (macro) so there is room to adapt
		// down, allow hints up through a concept/API pointer.
		return ModeDefaults{
			Depth: learning.DepthMacro, Help: learning.HelpProgressive,
			Disclosure: learning.DisclosureConceptOrAPI, Evaluation: learning.EvaluationOnDemand,
		}
	case learning.ModeReview:
		// "Problemas curtos de competências já estudadas": short (meso),
		// limited hinting since the competency was already introduced.
		return ModeDefaults{
			Depth: learning.DepthMeso, Help: learning.HelpLimited,
			Disclosure: learning.DisclosureGuidingQuestion, Evaluation: learning.EvaluationOnStepComplete,
		}
	case learning.ModeDebug:
		// "Investigação guiada sem revelar a causa antes da hipótese":
		// hints stop at a logical outline, never a pseudocode/skeleton
		// rung that would hand over the cause.
		return ModeDefaults{
			Depth: learning.DepthMacro, Help: learning.HelpProgressive,
			Disclosure: learning.DisclosureLogicalOutline, Evaluation: learning.EvaluationOnDemand,
		}
	case learning.ModeExploration:
		// "Feedback livre, sem obrigação de avanço ou avaliação": no
		// hint restriction; evaluation stays on-demand, never forced.
		return ModeDefaults{
			Depth: learning.DepthLayer, Help: learning.HelpFree,
			Disclosure: learning.DisclosureLogicalOutline, Evaluation: learning.EvaluationOnDemand,
		}
	case learning.ModeInterview:
		// "Briefing completo, ... auxílio bloqueado ou limitado e
		// revisão ao final" — the interview protocol itself belongs to
		// the interview-mode spec; this is only the dimension defaults a
		// session falls back to before that protocol configures it
		// further.
		return ModeDefaults{
			Depth: learning.DepthChallenge, Help: learning.HelpNoHints,
			Disclosure: learning.DisclosureNone, Evaluation: learning.EvaluationOnlyAtEnd,
		}
	default: // learning.ModePractice and any unrecognized value
		// "Menos contexto, pistas registradas e foco em execução":
		// narrow (micro), progressive hints (every use is recorded
		// regardless of mode).
		return ModeDefaults{
			Depth: learning.DepthMicro, Help: learning.HelpProgressive,
			Disclosure: learning.DisclosureGuidingQuestion, Evaluation: learning.EvaluationOnDemand,
		}
	}
}

// DebugProtocolStages is the ordered debugging micro-step taxonomy
// PROJECT.md §8.1 "Depuração" declares. It is exposed so the tutor skill
// and any future structured debug tooling present the same order instead
// of each re-deriving it (requirement R6).
func DebugProtocolStages() []string {
	return []string{
		"reproduce",
		"locate_first_divergence",
		"hypothesize",
		"observe",
		"confirm_or_reject",
		"apply_smallest_fix",
	}
}
