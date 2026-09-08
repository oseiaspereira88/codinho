---
title: Normalize root module metadata keys during indexing
type: bug
created_at: 2026-09-08T16:00:00Z
status: staged
privacy: sanitized-synthetic
---

# Normalize root module metadata keys during indexing

## Observation

POSE 1.8.1 indexes a root Go module with path "" but looks up metadata by
that exact string. Metadata declared under "." is ignored and the projection
reports defaulted/incomplete despite all required fields being present.
Assessment discovery normalizes both paths to the root directory.

## Synthetic reproduction

Create a root go.mod with module example.invalid/demo and declare metadata
for "." in .pose/indexes/module-metadata.json with owner @team, criticality
high, domain runtime and validationProfile baseline. Run pose index.
The root package in repo-map.json has path "", metadataStatus.source
"defaulted", and all four fields reported missing. Declaring an identical
entry under "" makes that projection complete. This is a compatibility
workaround, not two distinct components.

## Proposed solution

Normalize empty and dot root paths before indexing metadata and add regression
coverage shared with assess/review. Reject conflicting aliases explicitly.
Keep the public index path format compatible. No upstream submission made.
