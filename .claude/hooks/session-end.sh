#!/usr/bin/env bash
# session-end.sh - Stop hook for session-end verification
set -euo pipefail

REPO_ROOT=$(git -C "$(dirname "$0")" rev-parse --show-toplevel 2>/dev/null) || exit 0

# Files modified vs HEAD (unstaged + staged, deduplicated)
MODIFIED=$(git -C "$REPO_ROOT" diff --name-only HEAD 2>/dev/null || true)
STAGED=$(git -C "$REPO_ROOT" diff --name-only --cached 2>/dev/null || true)
ALL_CHANGED=$(printf '%s\n%s\n' "$MODIFIED" "$STAGED" | sort -u | grep -v '^$' || true)

[ -z "$ALL_CHANGED" ] && exit 0

ISSUES=""

block() {
  printf '{"decision":"block","reason":"%s"}' "$1"
  exit 0
}

# Check for merge conflict markers
while IFS= read -r f; do
  FULL="$REPO_ROOT/$f"
  [ -f "$FULL" ] || continue
  if grep -qE '^(<{7}|>{7}|={7}) ' "$FULL" 2>/dev/null; then
    ISSUES="${ISSUES}\\n- Unresolved merge conflict in: $f"
  fi
done <<< "$ALL_CHANGED"

# Check for hardcoded secret patterns in modified non-test Go files
GO_FILES=$(printf '%s\n' "$ALL_CHANGED" | grep '\.go$' | grep -v '_test\.go$' || true)
SECRET_RE='(password|passwd|api_key|apikey|secret_key|auth_token|access_key)\s*[:=]\s*"[^"]+'
while IFS= read -r f; do
  [ -z "$f" ] && continue
  FULL="$REPO_ROOT/$f"
  [ -f "$FULL" ] || continue
  if grep -iqE "$SECRET_RE" "$FULL" 2>/dev/null; then
    ISSUES="${ISSUES}\\n- Possible hardcoded secret in: $f"
  fi
done <<< "$GO_FILES"

# Check for stray debug prints in non-test Go files
while IFS= read -r f; do
  [ -z "$f" ] && continue
  FULL="$REPO_ROOT/$f"
  [ -f "$FULL" ] || continue
  if grep -qE 'fmt\.(Print|Println|Printf)\(' "$FULL" 2>/dev/null; then
    ISSUES="${ISSUES}\\n- Stray fmt.Print* in: $f (use tracing.Log instead)"
  fi
done <<< "$GO_FILES"

[ -z "$ISSUES" ] && exit 0

block "Session-end verification found issues:${ISSUES}\\n\\nPlease review and resolve before finishing."
