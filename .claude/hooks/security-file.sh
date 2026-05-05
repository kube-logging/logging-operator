#!/usr/bin/env bash
# security-file.sh - PreToolUse security hook for Write and Edit tools
set -euo pipefail

INPUT=$(cat)
FILE_PATH=$(printf '%s' "$INPUT" | jq -r '.tool_input.file_path // ""')

if [ -z "$FILE_PATH" ]; then exit 0; fi

deny() {
  printf '{"hookSpecificOutput":{"hookEventName":"PreToolUse","permissionDecision":"deny","permissionDecisionReason":"%s"}}' "$1"
  exit 0
}

# .env files
if printf '%s' "$FILE_PATH" | grep -qE '(^|/)\.env(\.[a-zA-Z0-9_-]+)?$'; then
  deny "Blocked: Writing to .env file is prohibited. Manage secrets outside the codebase."
fi

# secrets/ directories
if printf '%s' "$FILE_PATH" | grep -qE '(^|/)secrets/'; then
  deny "Blocked: Writing to secrets/ directory is prohibited."
fi

# Cryptographic key/certificate extensions
if printf '%s' "$FILE_PATH" | grep -qE '\.(pem|key|pfx|p12|jks|keystore)$'; then
  deny "Blocked: Writing to cryptographic key/certificate file is prohibited."
fi

# Credential file name patterns
if printf '%s' "$FILE_PATH" | grep -qE '(^|/)(credentials|id_rsa|id_ed25519|id_ecdsa|\.htpasswd)(\..*)?$'; then
  deny "Blocked: Writing to credential file is prohibited."
fi

# Files with 'secret' or 'token' in name
if printf '%s' "$FILE_PATH" | grep -qE '(^|/)[^/]*(secret|token)[^/]*$'; then
  deny "Blocked: Writing to file with '"'"'secret'"'"' or '"'"'token'"'"' in name is prohibited."
fi

# .npmrc (may contain auth tokens)
if printf '%s' "$FILE_PATH" | grep -qE '(^|/)\.npmrc$'; then
  deny "Blocked: Writing to .npmrc is prohibited (may contain registry auth tokens)."
fi

# .git internals
if printf '%s' "$FILE_PATH" | grep -qE '(^|/)\.git/(config|credentials|objects|FETCH_HEAD|packed-refs)'; then
  deny "Blocked: Direct writes to .git internals are prohibited. Use git commands instead."
fi

# SQLite databases
if printf '%s' "$FILE_PATH" | grep -qE '\.sqlite(3)?$'; then
  deny "Blocked: Direct writes to SQLite database files are prohibited."
fi

exit 0
