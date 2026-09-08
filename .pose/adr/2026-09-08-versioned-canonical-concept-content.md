# ADR: Versioned canonical concept content

## Status
Accepted — 2026-09-08. Spec: concept-content-authoring.

## Context
ConceptAuthoring currently contains only id/title. The read-only MCP tool must
serve canonical explanations without generating text or exposing challenge
solutions. Existing packs and session snapshots must remain readable.

## Decision
Add optional content to ConceptAuthoring under pack schema 1, with its own
integer version (initially 1). Omitted/null means missing; a present object must
be complete and bounded. Require explanation and one example with distinct
context, code and explanation. Allow an analogy and up to eight references to
existing outgoing concept-to-concept graph relations. Resolve only public
concept identities; never traverse to challenge/fixture/hint records.

Expose id/title unchanged plus content_status and content through the existing
concept_content_get tool. Sort relation projections by kind then concept_id.
Copy nested data on catalog insertion/query. Do not mutate sessions, emit
learning events, execute example code or contact a model/network service.

Limit text by Unicode code points, matching JSON Schema maxLength: explanation
8000; example context 200, code 4000, explanation 2000; analogy 2000. Require
non-whitespace text for mandatory fields. Reject unsupported versions and
invalid, duplicate or non-concept relation references at load time.

## Consequences
Legacy concepts explicitly report missing; the server does not invent fallback
explanations. New content changes require a pack version increment. Public
concept content is not a solution-disclosure capability. Authors must review
example context and semantic overlap: structural validation cannot prove that
arbitrary prose is not a copied answer. Samples and regression tests cover the
supported path; publication/playtesting remain separate pack gates.

Rejected: embedding an LLM duplicates tutor responsibility; reading challenge
solutions to explain a concept bypasses disclosure; requiring content on every
legacy concept would invalidate existing packs without improving their content.

Review when introducing content version 2, external references, personalization
or session-specific concept filtering. See the knowledge decision log.
