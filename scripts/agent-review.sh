#!/usr/bin/env bash
# agent-review.sh — thin, agent-agnostic invocation primitive for the
# multi-agent author/review loop documented in
# docs/agent-review-workflow.md.
#
# This script does NOT know about "author" or "reviewer" roles — that
# semantics (what to read, what never to edit, how to phrase the prompt)
# lives in the doc/prompt templates, not here. This script only knows how
# to start or resume one non-interactive, sandboxed session with a given
# backend agent CLI and capture its final message.
#
# Usage:
#   scripts/agent-review.sh new    <agent> <output-file> <prompt-text|@prompt-file>
#   scripts/agent-review.sh resume <agent> <output-file> <prompt-text|@prompt-file>
#
# <agent> selects the backend command template. Only "codex" is
# implemented today; add a case below to support another CLI (another
# Claude instance invoked non-interactively, an "agy"/Antigravity CLI,
# etc.) without touching any caller.
#
# Examples:
#   scripts/agent-review.sh new codex /tmp/review-round1.md @/tmp/round1-prompt.txt
#   scripts/agent-review.sh resume codex /tmp/review-round2.md "Only X changed, re-check that."
set -euo pipefail

if [[ $# -lt 4 ]]; then
  echo "usage: $0 <new|resume> <agent> <output-file> <prompt-text|@prompt-file>" >&2
  exit 1
fi

mode="$1"
agent="$2"
output="$3"
shift 3
prompt="$*"
if [[ "$prompt" == @* ]]; then
  prompt="$(cat "${prompt:1}")"
fi

case "$agent" in
  codex)
    # Model pinned explicitly rather than left to ~/.codex/config.toml's
    # default: this backend is used as the project's reviewer, and an
    # implicit dependency on whatever the local global config happens to
    # say would make the "adversarial, independent model" property of the
    # review silently fragile to unrelated config changes.
    codex_model=(-c model="gpt-5.6-luna" -c model_reasoning_effort="high")
    case "$mode" in
      new)
        # -s workspace-write: lets the agent run build/test/validate commands.
        # --skip-git-repo-check: this script may run from a worktree/subdir.
        # Non-interactive by construction (codex exec never prompts).
        codex exec "${codex_model[@]}" -s workspace-write --skip-git-repo-check -o "$output" "$prompt"
        ;;
      resume)
        # Resumes the most recent session for this cwd — no session id
        # bookkeeping needed as long as rounds run sequentially from the
        # same working directory. The resumed session already has every
        # file it read in earlier rounds in context: do not re-paste
        # file contents in $prompt, only describe what changed.
        codex exec "${codex_model[@]}" resume --last -o "$output" -- "$prompt"
        ;;
      *)
        echo "unknown mode: $mode (use new|resume)" >&2
        exit 1
        ;;
    esac
    ;;
  *)
    echo "unknown agent backend: '$agent' — add a case in $0" >&2
    exit 1
    ;;
esac

cat "$output"
