// Package curriculum loads, validates and indexes authored YAML packs into
// an immutable, queryable catalog. It never generates content, never
// recommends a path, and never runs checks; it only turns authored YAML
// into validated Go values (PROJECT.md §14, non-goals of this spec).
package curriculum

// SchemaVersion is the pack authoring schema version this loader
// understands. Loading a pack authored against a different version fails
// with DiagIncompatibleSchemaVersion before the catalog is materialized
// (requirement R7).
const SchemaVersion = 1

// ItemKind discriminates the queryable catalog entity types (requirement
// R1).
type ItemKind string

const (
	KindTheme      ItemKind = "theme"
	KindConcept    ItemKind = "concept"
	KindCompetency ItemKind = "competency"
	KindTrack      ItemKind = "track"
	KindChallenge  ItemKind = "challenge"
)

// PackManifest lists the pack files that make up a catalog
// (packs/manifest.yaml).
type PackManifest struct {
	SchemaVersion int      `yaml:"schema_version"`
	Packs         []string `yaml:"packs"`
}

// Pack is one authored YAML file: a self-contained set of themes, concepts,
// competencies, tracks, challenges and relations, versioned as a unit.
type Pack struct {
	SchemaVersion int                   `yaml:"schema_version"`
	ID            string                `yaml:"id"`
	Version       string                `yaml:"version"`
	Themes        []ThemeAuthoring      `yaml:"themes"`
	Concepts      []ConceptAuthoring    `yaml:"concepts"`
	Competencies  []CompetencyAuthoring `yaml:"competencies"`
	Tracks        []TrackAuthoring      `yaml:"tracks"`
	Challenges    []ChallengeAuthoring  `yaml:"challenges"`
	// Relations is additive to Prerequisites (curriculum-graph-path-
	// recommendation Decision 2): it connects any two catalog items by one
	// of the eight declared kinds, independent of the older, challenge-only
	// Prerequisites mechanism.
	Relations   []RelationAuthoring  `yaml:"relations"`
	Publication PublicationAuthoring `yaml:"publication,omitempty"`

	// File is the source path, set by the loader for diagnostics. It is not
	// part of the authored YAML shape.
	File string `yaml:"-"`
}

// RelationAuthoring is one authored edge between two catalog items,
// identified by ID regardless of kind (curriculum-graph-path-
// recommendation, requirement R1). Kind must be one of the eight declared
// RelationKind values; the loader rejects anything else explicitly rather
// than ignoring it (Compatibility).
type RelationAuthoring struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
	Kind string `yaml:"kind"`
}

// ThemeAuthoring is a knowledge area (PROJECT.md §7.1, §14.2).
type ThemeAuthoring struct {
	ID    string `yaml:"id"`
	Title string `yaml:"title"`
}

// ConceptAuthoring is something a learner can understand (PROJECT.md §7.2).
type ConceptAuthoring struct {
	ID      string                   `yaml:"id"`
	Title   string                   `yaml:"title"`
	Content *ConceptContentAuthoring `yaml:"content,omitempty"`
}

// ConceptContentAuthoring is public explanatory content, never challenge data.
type ConceptContentAuthoring struct {
	Version      int                        `yaml:"version"`
	Explanation  string                     `yaml:"explanation"`
	Example      ConceptExampleAuthoring    `yaml:"example"`
	Analogy      string                     `yaml:"analogy,omitempty"`
	RelationRefs []ConceptRelationReference `yaml:"relation_refs,omitempty"`
}

type ConceptExampleAuthoring struct {
	Context     string `yaml:"context" json:"context"`
	Code        string `yaml:"code" json:"code"`
	Explanation string `yaml:"explanation" json:"explanation"`
}

type ConceptRelationReference struct {
	Kind      string `yaml:"kind"`
	ConceptID string `yaml:"concept_id"`
}

// CompetencyAuthoring is something a learner can demonstrate (PROJECT.md
// §7.2).
type CompetencyAuthoring struct {
	ID    string `yaml:"id"`
	Title string `yaml:"title"`
}

// TrackAuthoring is an ordered journey toward a learning goal (PROJECT.md
// §7.1, §14.6).
type TrackAuthoring struct {
	ID     string   `yaml:"id"`
	Title  string   `yaml:"title"`
	Themes []string `yaml:"themes"`
}

// ChallengeAuthoring mirrors the challenge anatomy of PROJECT.md §14.8.
type ChallengeAuthoring struct {
	SchemaVersion    int              `yaml:"schema_version"`
	ID               string           `yaml:"id"`
	Version          string           `yaml:"version"`
	Title            string           `yaml:"title"`
	Kind             string           `yaml:"kind"`
	Difficulty       string           `yaml:"difficulty"`
	EstimatedMinutes int              `yaml:"estimated_minutes"`
	Themes           []string         `yaml:"themes"`
	Competencies     CompetencyRefs   `yaml:"competencies"`
	Prerequisites    []string         `yaml:"prerequisites"`
	Brief            string           `yaml:"brief"`
	Constraints      []string         `yaml:"constraints"`
	Acceptance       []string         `yaml:"acceptance"`
	Layers           []LayerAuthoring `yaml:"layers"`
	Checks           []CheckAuthoring `yaml:"checks"`
	// Variants hold alternate phrasings reserved for reuse (e.g. interview
	// mode) and are excluded from default catalog queries (requirement R6).
	Variants []string `yaml:"variants"`
	// Fixture lists starter files this challenge ships (e.g. buggy code for
	// a debug challenge). It is authored inline, like every other field, so
	// it is versioned and reviewed the same way as the rest of the
	// challenge; the administrative CLI's workspace-prepare command is the
	// only thing that ever materializes it onto disk (administrative-cli-
	// fixtures requirement R6).
	Fixture []FixtureFileAuthoring `yaml:"fixture"`
	// Publication declares editorial metadata gating status: published
	// (catalog-authoring-quality requirement R7). A zero value (empty
	// Status) is an ordinary, unpublished authoring draft and is never
	// flagged by the editorial gate.
	Publication PublicationAuthoring `yaml:"publication"`
	Canonical   bool                 `yaml:"canonical,omitempty" json:"Canonical,omitempty"`
	VariantOf   string               `yaml:"variant_of,omitempty" json:"VariantOf,omitempty"`
	// Validation is private editorial input; never serialize into public sessions.
	Validation *CheckValidationAuthoring `yaml:"validation,omitempty" json:"-"`
}

// PublicationAuthoring is the editorial metadata catalog-authoring-
// quality requires before a challenge may declare `status: published`:
// who authored it, who reviewed it (a different person, per Constraint
// "revisão... por pessoa diferente do autor"), and whether it was
// playtested.
type PublicationAuthoring struct {
	Status     string `yaml:"status"`
	Author     string `yaml:"author"`
	ReviewedBy string `yaml:"reviewed_by"`
	Playtested bool   `yaml:"playtested"`
}

// StatusPublished is the only PublicationAuthoring.Status value the
// editorial gate enforces metadata for.
const StatusPublished = "published"

// FixtureFileAuthoring is one file a workspace-prepare command materializes
// verbatim into the learner's chosen destination.
type FixtureFileAuthoring struct {
	Path    string `yaml:"path"`
	Content string `yaml:"content"`
}

// CompetencyRefs distinguishes the competency a challenge primarily
// exercises from ones it merely touches.
type CompetencyRefs struct {
	Primary   []string `yaml:"primary"`
	Secondary []string `yaml:"secondary"`
}

// LayerAuthoring is a logical perspective of a challenge holding its
// macro-step tree (PROJECT.md §7.1, §14.8).
type LayerAuthoring struct {
	ID         string          `yaml:"id"`
	MacroSteps []StepAuthoring `yaml:"macro_steps"`
}

// StepAuthoring mirrors the micropasso anatomy of PROJECT.md §14.9. It is
// recursive: a macro or meso step nests its children under the same shape.
type StepAuthoring struct {
	ID           string               `yaml:"id"`
	Kind         string               `yaml:"kind"`
	Action       string               `yaml:"action"`
	Target       string               `yaml:"target"`
	Title        string               `yaml:"title"`
	Instruction  InstructionAuthoring `yaml:"instruction"`
	Concepts     []string             `yaml:"concepts"`
	Evidence     EvidenceAuthoring    `yaml:"evidence"`
	Criteria     []CriterionAuthoring `yaml:"criteria"`
	Completion   CompletionAuthoring  `yaml:"completion"`
	Hints        []HintAuthoring      `yaml:"hints"`
	Reflection   ReflectionAuthoring  `yaml:"reflection"`
	Children     []StepAuthoring      `yaml:"children"`
	ChildrenMode string               `yaml:"children_mode,omitempty" json:"ChildrenMode,omitempty"`
}

// InstructionAuthoring is the single instruction disclosed for a step.
type InstructionAuthoring struct {
	Objective   string   `yaml:"objective"`
	Scope       string   `yaml:"scope"`
	Constraints []string `yaml:"constraints"`
}

// EvidenceAuthoring declares how a step's completion can be observed.
type EvidenceAuthoring struct {
	Strategies []string `yaml:"strategies"`
}

// CriterionAuthoring is one evaluable criterion for a step.
type CriterionAuthoring struct {
	ID   string `yaml:"id"`
	Kind string `yaml:"kind"`
}

// CompletionAuthoring is the completion policy declared for a step.
type CompletionAuthoring struct {
	RequiresPositiveEvaluation bool `yaml:"requires_positive_evaluation"`
	RequiresUserConfirmation   bool `yaml:"requires_user_confirmation"`
}

// HintAuthoring is one rung of the assistance ladder for a step (PROJECT.md
// §8.4).
type HintAuthoring struct {
	Level int    `yaml:"level"`
	Kind  string `yaml:"kind"`
}

// ReflectionAuthoring lists optional reflection prompts for a step
// (PROJECT.md §8.9).
type ReflectionAuthoring struct {
	Optional []string `yaml:"optional"`
}

// CheckAuthoring declares a check by ID; the executor that resolves and
// runs it belongs to safe-check-executor. Network defaults to denied
// (safe-check-executor R5): only a check the challenge explicitly marks
// true may reach the network, never something a caller decides at
// check_run time (check_run takes no free parameter beyond check_id).
type CheckAuthoring struct {
	ID          string `yaml:"id"`
	Runner      string `yaml:"runner"`
	Package     string `yaml:"package"`
	TestPattern string `yaml:"test_pattern"`
	Timeout     string `yaml:"timeout"`
	Network     bool   `yaml:"network"`
}

// CheckValidationAuthoring supplies reproducible, private editorial scenarios.
// Each fixture is complete; reference files never overlay the starter workspace.
type CheckValidationAuthoring struct {
	BaselineFixture  []FixtureFileAuthoring      `yaml:"baseline_fixture,omitempty"`
	Justification    string                      `yaml:"justification,omitempty"`
	ReferenceFixture []FixtureFileAuthoring      `yaml:"reference_fixture"`
	Expectations     []CheckExpectationAuthoring `yaml:"expectations"`
}

type CheckExpectationAuthoring struct {
	CheckID   string `yaml:"check_id"`
	Baseline  string `yaml:"baseline"`
	Reference string `yaml:"reference"`
}
