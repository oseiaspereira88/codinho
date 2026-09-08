package curriculum

import "strings"

// RulePublishedWithoutAuthor, RulePublishedWithoutReview,
// RulePublishedSameReviewer and RulePublishedWithoutPlaytest enforce the
// Constraint "revisão... por pessoa diferente do autor" for any
// challenge declaring publication.status: published (requirement R7).
const (
	RulePublishedWithoutAuthor   EditorialRuleID = "published_without_author"
	RulePublishedWithoutReview   EditorialRuleID = "published_without_review"
	RulePublishedSameReviewer    EditorialRuleID = "published_same_reviewer"
	RulePublishedWithoutPlaytest EditorialRuleID = "published_without_playtest"
)

// checkPublicationMetadata is a no-op for any challenge that has not
// declared publication.status: published — an ordinary authoring draft
// is never flagged. Once published is declared, every field this checks
// must be filled honestly by a human: this function only confirms the
// fields are non-empty and that the reviewer differs from the author; it
// can never confirm a playtest actually happened, which is exactly the
// qualitative judgment Decision 1 reserves for people, not automation.
func checkPublicationMetadata(file string, ch ChallengeAuthoring) []EditorialFinding {
	ch.Publication.Author = strings.TrimSpace(ch.Publication.Author)
	ch.Publication.ReviewedBy = strings.TrimSpace(ch.Publication.ReviewedBy)
	if ch.Publication.Status != StatusPublished {
		return nil
	}
	var out []EditorialFinding
	if ch.Publication.Author == "" {
		out = append(out, EditorialFinding{
			File: file, Item: ch.ID, Rule: RulePublishedWithoutAuthor, Severity: SeverityBlocking,
			Detail: "status published sem publication.author", Suggestion: "preencha publication.author",
		})
	}
	if ch.Publication.ReviewedBy == "" {
		out = append(out, EditorialFinding{
			File: file, Item: ch.ID, Rule: RulePublishedWithoutReview, Severity: SeverityBlocking,
			Detail: "status published sem publication.reviewed_by", Suggestion: "preencha publication.reviewed_by com uma pessoa diferente do autor",
		})
	} else if ch.Publication.Author != "" && strings.EqualFold(ch.Publication.ReviewedBy, ch.Publication.Author) {
		out = append(out, EditorialFinding{
			File: file, Item: ch.ID, Rule: RulePublishedSameReviewer, Severity: SeverityBlocking,
			Detail: "publication.reviewed_by é igual a publication.author", Suggestion: "a revisão deve ser feita por uma pessoa diferente do autor",
		})
	}
	if !ch.Publication.Playtested {
		out = append(out, EditorialFinding{
			File: file, Item: ch.ID, Rule: RulePublishedWithoutPlaytest, Severity: SeverityBlocking,
			Detail: "status published sem publication.playtested", Suggestion: "defina publication.playtested: true somente após um playtest real",
		})
	}
	return out
}
