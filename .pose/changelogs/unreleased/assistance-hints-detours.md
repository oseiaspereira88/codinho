---
spec: assistance-hints-detours
category: added
breaking: false
refs:
---

`codinho serve` now supports the assistance ladder and conceptual detours:
`hint_request` delivers progressively deeper hints (up to a commented
solution, which requires explicit confirmation), `syntax_recall_get`
offers a lighter-cost syntax reminder, `concept_content_get` returns
canonical concept records, and `learning_detour_start`/
`learning_detour_finish` let the learner ask a side question without
losing their place.
