---
title: "ListReviewBundles aborts entire scoped listing when an unrelated bundle file fails to load"
type: "bug"
module: ""
created_at: "2026-08-23T01:45:59Z"
status: staged
upstream: "https://github.com/oseiaspereira88/pose"
privacy: sanitized-synthetic
---

# Contribution Draft: ListReviewBundles aborts entire scoped listing when an unrelated bundle file fails to load

## Summary
`Store.ListReviewBundles(scope)` (internal/pose/review_bundle.go:1058) loads
every `.json` file under `.pose/review-bundles/` via `LoadReviewBundle`
*before* filtering by `scope`, and returns `nil, err` immediately if any
single file fails to load — even when that file's `Payload.Scope.Ref` does
not match the requested `scope`.

Consequence: a single corrupted/unreadable review bundle anywhere in the
directory poisons `ListReviewBundles` for every other scope that calls it.
Since `VerifyReviewBundle` → `ReviewCheck` → `GetCloseoutState` all go
through `ListReviewBundles`, `pose check --strict`'s review-closeout gate
reported the exact same "review bundle <id> digest mismatch" for 7
unrelated specs — all citing the alphabetically-first bundle filename in
the directory, none citing their own actual bundle.

## Reproduction (synthetic)
1. Seal review bundles for two independent specs, A and B (`pose review
   bundle spec:A --seal`, `pose review bundle spec:B --seal`).
2. Hand-edit spec A's sealed bundle JSON file (any byte change to a string
   value inside `payload`, e.g. renaming a path), so its stored
   `bundle_digest` no longer matches the recomputed payload digest.
3. Run `pose review verify spec:B` (B's own bundle is untouched and valid).

Expected: B's review closeout is verified/reported independently of A's
corrupted bundle.
Actual: `pose: review bundle <A's id> digest mismatch` — B's verification
fails citing A's bundle, purely because A's file sorts before B's in
`os.ReadDir` order.

## Root cause
```go
// internal/pose/review_bundle.go:1074-1085
for _, entry := range entries {
    if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
        continue
    }
    bundle, loadErr := s.LoadReviewBundle(strings.TrimSuffix(entry.Name(), ".json"))
    if loadErr != nil {
        return nil, loadErr   // <-- aborts for ALL scopes, not just the one asked for
    }
    if scope == "" || bundle.Payload.Scope.Ref == scope {
        bundles = append(bundles, bundle)
    }
}
```

## Suggested fix
Skip (optionally collecting as a warning) bundle files that fail to load
when `scope != ""` and the failure isn't confirmed to belong to that scope
— e.g. peek at the raw JSON's `payload.scope.ref` field before calling the
strict `LoadReviewBundle`, or continue past `loadErr` for entries whose
digest-independent scope ref doesn't match, only surfacing the load error
when the matching-scope bundle itself is unreadable.

## Context
Found while investigating a real (separate, non-engine) incident: a
project-wide rename script had string-replaced content inside two already
sealed bundle JSON files, correctly invalidating just those two bundles'
self-consistency. The invalidation was expected; the fact that it silently
masked the closeout status of five completely unrelated, untouched, valid
bundles was not.

## Upstream issue
https://github.com/oseiaspereira88/pose/issues/40
