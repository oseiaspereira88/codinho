---
title: "Make knowledge-usage count citations in specs"
type: "bug"
module: ""
created_at: "2026-09-15T18:29:42Z"
status: staged
upstream: "https://github.com/oseiaspereira88/pose"
privacy: sanitized-synthetic
---

# Contribution Draft: Make knowledge-usage count citations in specs

## Observation

`pose knowledge-check --strict` accepts valid knowledge artifacts and reports
zero reference failures, while a spec containing the required literal form
`knowledge:<slug>` is present in the repository. `pose knowledge-usage` still
reports `citations:0` for every artifact and `usage.cited=0`.

This makes the usage projection disagree with the citation contract documented
by the feature and bugfix workflows. It can also make a consulted artifact look
unused and approach TTL expiry even when the spec contains the citation.

## Synthetic reproduction

1. Create a repository with a valid knowledge artifact whose slug is
   `sample-handoff`.
2. Create a spec containing the literal `knowledge:sample-handoff` in its
   execution log or another prose section.
3. Run `pose knowledge-check --strict` and then `pose knowledge-usage`.

The first command should pass, and the second should project at least one
citation for `sample-handoff`. No project code, private host, credential, or
customer data is needed.

## Proposed improvement

Make `knowledge-usage` scan the same canonical spec roots and literal citation
form accepted by `knowledge-check`, and add regression coverage for citations
in execution logs and requirement text. Keep dangling references as errors and
avoid counting arbitrary prose that does not match `knowledge:<slug>`.

This is local feedback only; no upstream submission was made and no POSE gate
was bypassed.
