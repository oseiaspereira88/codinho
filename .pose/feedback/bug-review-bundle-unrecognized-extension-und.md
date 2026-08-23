---
title: "review bundle: unrecognized non-root dotdir treated as unclassified blocker regardless of file extension"
kind: bug
engine_version: 1.7.10
reported_at: 2026-08-23T05:00:19-03:00
---

# POSE Engine Report: review bundle: unrecognized non-root dotdir treated as unclassified blocker regardless of file extension

## Description
POSE 1.7.9: `pose review bundle spec:<slug> --seal` fails with `[ERROR] unclassified review subject path` for any attributed file living under a brand-new top-level dotdir (e.g. `.codex/`) that is not one of the classifier's known roots (`.pose/`, `.agents/`, `.claude/`, `.github/`, etc.) and is not at the repo root either, so this is a distinct reproduction from the already-reported root-level-files issue (oseiaspereira88/pose#38): the component-skip-at-root heuristic described in #38 does not apply here since the path is not directly at `.`. Initial hypotheses about the file extension were ruled out experimentally: renaming to avoid a double extension (`config.toml.example` to `config.example.toml`) did not help, and converting the file entirely to `.md` with the same content in a fenced toml code block *also* failed identically with the same unclassified error — proving the extension is irrelevant. The actual root cause is that the classifier's path-prefix table has no entry for the `.codex/` directory itself (or presumably any other arbitrary top-level dotdir outside its known allowlist), and, same failure mode as #38, an unclassified attributed path becomes a hard `[ERROR]` blocker instead of being excluded from the subject. Workaround used: eliminated the new directory entirely and moved the same example content into a file under an already-recognized tree (`.agents/skills/codinho/references/codex-configuration.md`), which classifies fine (at most a non-fatal `[WARN] unmapped review component`). Suggested fix: either extend the classifier's known top-level-path allowlist to be extensible/configurable per-project (so a new dotdir a project introduces, e.g. `.codex/`, can be recognized), and/or (matching #38's suggestion) make unclassified always mean excluded from subject, never simultaneously blocking. Reproduction needs no private paths: `pose init`, create `.some-new-dotdir/anyfile.md` (extension does not matter), attribute it to a spec's declared Artifacts, run `pose review bundle --seal`.

---
### System Context (Auto-generated)
- **POSE Engine Version:** 1.7.10
- **OS/Arch:** linux/amd64
- **Go Version:** go1.26.5-X:nodwarf5
- **Reported At:** 2026-08-23T05:00:19-03:00
- **Kind:** bug

