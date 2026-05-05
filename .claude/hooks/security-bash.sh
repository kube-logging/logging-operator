#!/usr/bin/env bash
# security-bash.sh - PreToolUse security hook for Bash tool
set -euo pipefail

INPUT=$(cat)
COMMAND=$(printf '%s' "$INPUT" | jq -r '.tool_input.command // ""')

if [ -z "$COMMAND" ]; then exit 0; fi

deny() {
  printf '{"hookSpecificOutput":{"hookEventName":"PreToolUse","permissionDecision":"deny","permissionDecisionReason":"%s"}}' "$1"
  exit 0
}

# git reset --hard - irreversible local state destruction
if printf '%s' "$COMMAND" | grep -qE '(^|[|;&[:space:]])git\s+reset\s+--hard'; then
  deny "Blocked: git reset --hard discards changes irreversibly. Use git stash or checkout specific files instead."
fi

# git checkout -- (overwrites working tree files)
if printf '%s' "$COMMAND" | grep -qE '(^|[|;&[:space:]])git\s+checkout\s+--\s'; then
  deny "Blocked: git checkout -- overwrites working tree files irreversibly. Stage or stash changes first."
fi

# git clean -f (deletes untracked files)
if printf '%s' "$COMMAND" | grep -qE '(^|[|;&[:space:]])git\s+clean\s+(-[a-zA-Z]*f|.*-f)'; then
  deny "Blocked: git clean -f deletes untracked files permanently. Review with git status first."
fi

# Piping network content to a shell interpreter
if printf '%s' "$COMMAND" | grep -qE '(curl|wget).*(bash|sh|zsh|fish)\b|\|\s*(bash|sh|zsh|fish)\b'; then
  deny "Blocked: Piping network content directly to a shell interpreter is a supply chain attack vector."
fi

# Shell -c with embedded rm -rf (obfuscation bypass)
if printf '%s' "$COMMAND" | grep -qE '(bash|sh|zsh)\s+-c\s+.*rm\s+-[rRfF]'; then
  deny "Blocked: Shell -c invocation with embedded rm -rf detected."
fi

# Redirect write to .env or secrets/
if printf '%s' "$COMMAND" | grep -qE '>\s*(\.env|\.env\.[a-zA-Z0-9_-]+|secrets/)'; then
  deny "Blocked: Redirect write to .env or secrets/ directory is prohibited."
fi

# Redirect write to credential/key files
if printf '%s' "$COMMAND" | grep -qE '>\s*[^[:space:]]*(\.pem|\.key|\.pfx|\.p12|credentials|id_rsa|id_ed25519)'; then
  deny "Blocked: Redirect write to credential/key file detected."
fi

exit 0
