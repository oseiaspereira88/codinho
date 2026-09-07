---
title: "Installed CLI can retain root-module evidence mismatch after source correction"
type: bug
created_at: "2026-09-07T14:35:00Z"
status: staged
privacy: sanitized-synthetic
---

# Observation

An installed CLI identifying as 1.7.10 rejected review evidence from module
`.` for a delivery in `cmd/example`, while strict validation and surface-check
passed. The installed executable embeds an older, modified VCS revision.
A clean committed source checkout contains the ancestor-module matching fix;
its existing focused regression passes and its compiled CLI accepts the same
evidence without changing the project's delivery contract or policies.

# Synthetic reproduction

Use a repository with one root go.mod, a `cmd/example/main.go` entrypoint,
an attributed spec targeting module:cmd/example and passed root validation.
Compare review bundle preparation between the older installed executable and
a build containing TestReviewBundleMatchesRootModuleValidationEvidenceForSubdirectoryTargets.
The former reports no attributed evidence; the latter recognizes root checks.
Renaming the target to module:. is insufficient in the older executable:
review component resolution rejects the root path. Do not alter project
architecture or fabricate validation modules to work around the mismatch.

# Proposed improvement

Expose build revision and dirty status in `pose version`; add release/install
smoke coverage for the root-module/subdirectory delivery regression. Link this
observation to the existing module-evidence contribution rather than creating
a duplicate upstream issue. This report contains only synthetic paths and
tooling behavior. No upstream submission was made.
