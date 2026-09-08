# ADR: Native CI evidence and prospective delivery governance

## Status
Accepted — 2026-09-08; implementation governed by v1-delivery-ci-assurance.

## Context
The current CI skips POSE when absent, selects mutable actions and an unpinned
scanner, and interpolates release input into shell code. Cross-compilation does
not establish native runtime compatibility. Historical delivery records also
refer to cmd/ailearn before its recorded rename to cmd/codinho.

## Decision
Use one required validation matrix from Makefile and native Linux/macOS CI.
Install POSE 1.7.12 from official release archives with repository-pinned SHA-256;
pin actions to full commits and govulncheck to a checksummed module version.
Select a patched Go toolchain explicitly in CI, preserving the module language
minimum. Pass release labels through environment variables and validate a closed
version grammar before building. Keep release jobs limited to build artifacts.

Validate structured results against the current matrix and Git commit, including
required check identity, command, outcome, freshness and matrix digest. Attach
native OS, architecture and Go version only from the executing process. Missing,
skipped or stale evidence blocks this gate. A successful CI run does not prove
human editorial acceptance or a complete V1 catalog.

Adopt artifact and delivery policies prospectively with repository-owned roots.
Reconcile historical rename provenance with the actual Git change set; preserve
immutable review bundles and attestations. Govern current user documentation
through an explicit manifest and record pending roadmap criteria honestly.

## Consequences
Tool updates require reviewed pin changes and native verification. Network/tool
failure is a failed gate. Linux ARM cross-builds remain build evidence until a
native job executes them. Historical records remain auditable, while current
entrypoints and contract mappings must be backed by actual files and tests.

Rejected: conditional POSE skipping hides missing governance. Mutable tool tags
make executions irreproducible. Rewriting historical attestations would erase
the source of evidence. Declaring macOS support from cross-builds proves less
than the contract requires.

References: [spec](../specs/2026-09-07-v1-delivery-ci-assurance.md),
[GitHub secure use](https://docs.github.com/en/actions/reference/security/secure-use)
(accessed 2026-09-08). Review pins, platform claims and reconciliation at the next
POSE or Go upgrade, or within 90 days.
