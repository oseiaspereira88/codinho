package session

import (
	"fmt"

	"github.com/oseiaspereira88/codinho/internal/learning"
	"github.com/oseiaspereira88/codinho/internal/mastery"
)

// EvidenceThreshold is how many accumulated signals SuggestGranularity
// requires before proposing a change — the "evidência configurável"
// PROJECT.md §8.5 requires. It is a variable, not a constant, so a future
// spec can tune it without changing this package's exported surface.
var EvidenceThreshold = 3

// depthLadder is PROJECT.md §8.5's granularity ladder, coarsest first:
// "dificuldade observada: macro → meso → micro → pista" moves toward the
// end of this slice; "facilidade repetida: micro → meso → macro →
// desafio autônomo" moves toward its start.
var depthLadder = []learning.Depth{
	learning.DepthChallenge, learning.DepthLayer, learning.DepthMacro, learning.DepthMeso, learning.DepthMicro,
}

// SuggestGranularity proposes a coarser or finer depth than current based
// on proj's accumulated mastery evidence, and explains why (PROJECT.md
// §8.5: "deve ser informado ao aluno"). It never mutates anything: the
// caller decides whether to apply the suggestion via GranularityAdjust,
// and an explicit choice from the learner always wins (Constraint:
// "Ajuste automático nunca supera escolha manual"). ok is false when
// proj's evidence has not yet crossed EvidenceThreshold, or current has
// no coarser/finer neighbor left to move to.
func SuggestGranularity(current learning.Depth, proj mastery.DimensionProjection) (suggested learning.Depth, reason string, ok bool) {
	if proj.EvidenceCount < EvidenceThreshold {
		return current, "", false
	}
	switch proj.State {
	case mastery.StateDemonstratesWithoutHelp, mastery.StateRetained, mastery.StateTransferred:
		if coarser, has := coarserDepth(current); has {
			return coarser, fmt.Sprintf("%d evidências recentes sem ajuda: ampliando de %s para %s", proj.EvidenceCount, current, coarser), true
		}
	case mastery.StateNotObserved, mastery.StateIntroduced, mastery.StateDemonstratesWithHelp:
		if finer, has := finerDepth(current); has {
			return finer, fmt.Sprintf("%d evidências ainda com dificuldade: reduzindo de %s para %s", proj.EvidenceCount, current, finer), true
		}
	}
	return current, "", false
}

func coarserDepth(d learning.Depth) (learning.Depth, bool) {
	for i, cur := range depthLadder {
		if cur == d {
			if i == 0 {
				return "", false
			}
			return depthLadder[i-1], true
		}
	}
	return "", false
}

func finerDepth(d learning.Depth) (learning.Depth, bool) {
	for i, cur := range depthLadder {
		if cur == d {
			if i == len(depthLadder)-1 {
				return "", false
			}
			return depthLadder[i+1], true
		}
	}
	return "", false
}
