---
title: "Retry transient editorial worktree creation failures"
type: "bug"
module: ""
created_at: "2026-09-16T00:56:30Z"
status: staged
upstream: "https://github.com/oseiaspereira88/pose"
privacy: sanitized-synthetic
---

# Contribution Draft: Retry transient editorial worktree creation failures

## Observation

The editorial batch executor marked a batch as failed before starting its
author because `git worktree add --detach <run>/batch-N/worktree <head>`
returned exit status 128. The target directory did not exist and the source
checkout was clean. Running the same command manually immediately afterward
created the detached worktree successfully, so the batch could be resumed
without changing its requested paths or content.

This happened on multiple sequential batches in one run. The failure consumed
an attempt and required manual worktree creation followed by `editorial.py
revise`; no author output or content error explained the original failure.

## Synthetic reproduction

1. Start an editorial run from a clean repository with a pending batch.
2. Let the executor dispatch a batch whose worktree path is a fresh directory
   under `/tmp/example-editorial/batch-001/worktree`.
3. Observe `git worktree add --detach` return status 128 and the batch become
   `failed` before an author session starts.
4. Run the same `git worktree add --detach` command manually with the same
   commit and fresh path; it succeeds, after which the batch can be resumed.

The reproduction uses only synthetic paths, a local Git repository, and a
detached commit; it does not require project code, private hosts, credentials,
or customer data.

## Proposed improvement

Make worktree dispatch resilient and diagnosable: retry transient
`git worktree add` failures with bounded backoff, verify the target path and
Git worktree metadata before retrying, and capture the command's stderr in the
batch state. If retries are exhausted, preserve a distinct infrastructure
failure classification without consuming an author attempt. Add an executor
regression test for a fresh-path retry and keep recovery idempotent.

This report is local feedback only; no upstream submission was made and no
POSE gate was bypassed.
