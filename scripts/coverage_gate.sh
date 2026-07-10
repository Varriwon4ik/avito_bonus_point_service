#!/usr/bin/env bash
#
# Per-module line-coverage gate (QRT-003, verifying QR-003 Testability).
#
# Reads a Go coverage profile and fails if any critical module is below the
# required line-coverage threshold (default 30%, overridable via
# COVERAGE_THRESHOLD). Critical modules are defined in docs/testing.md.
#
# Usage: bash scripts/coverage_gate.sh [coverage.out]
#
set -euo pipefail

PROFILE="${1:-coverage.out}"
THRESHOLD="${COVERAGE_THRESHOLD:-30}"

# Critical modules (package import-path prefixes) that must meet the threshold.
CRITICAL_MODULES=(
  "bonus-ledger/internal/data"
  "bonus-ledger/internal/api"
)

if [[ ! -f "$PROFILE" ]]; then
  echo "::error::coverage profile '$PROFILE' not found"
  exit 1
fi

# Aggregate covered/total statements per package from the coverage profile.
# Profile line format:
#   <import/path/file.go>:<startLine.col>,<endLine.col> <numStmts> <count>
# With -coverpkg, every test binary in the run emits its own copy of each
# instrumented block, so the same block appears once per test binary. Count
# each unique block once and treat it as covered if any binary executed it.
BY_PKG="$(mktemp)"
trap 'rm -f "$BY_PKG"' EXIT

awk '
  NR == 1 { next }                      # skip "mode:" header line
  {
    block_stmts[$1] = $2
    if ($3 + 0 > 0) block_hit[$1] = 1
  }
  END {
    for (b in block_stmts) {
      file = b; sub(/:.*/, "", file)    # strip ":line.col,line.col"
      pkg  = file; sub(/\/[^/]+$/, "", pkg)  # package import path = dir of file
      stmts[pkg] += block_stmts[b]
      if (b in block_hit) covered[pkg] += block_stmts[b]
    }
    for (p in stmts) {
      pct = (stmts[p] > 0) ? (covered[p] / stmts[p]) * 100 : 0
      printf "%s %d %d %.1f\n", p, covered[p] + 0, stmts[p], pct
    }
  }
' "$PROFILE" | sort > "$BY_PKG"

echo "Per-package line coverage:"
printf '  %-45s %-14s %s\n' "package" "covered/total" "coverage"
while read -r pkg cov tot pct; do
  printf '  %-45s %-14s %s%%\n' "$pkg" "${cov}/${tot}" "$pct"
done < "$BY_PKG"
echo

fail=0
for mod in "${CRITICAL_MODULES[@]}"; do
  read -r cov tot < <(awk -v m="$mod" '
    index($1, m) == 1 { c += $2; s += $3 }
    END { print c + 0, s + 0 }
  ' "$BY_PKG")

  if [[ "${tot:-0}" -eq 0 ]]; then
    echo "::error::critical module '$mod' has no measured statements"
    fail=1
    continue
  fi

  pct="$(awk -v c="$cov" -v s="$tot" 'BEGIN { printf "%.1f", (c / s) * 100 }')"
  meets="$(awk -v p="$pct" -v t="$THRESHOLD" 'BEGIN { print (p + 0 >= t + 0) ? "OK" : "FAIL" }')"
  printf 'module %-30s %s/%s = %s%% (threshold %s%%) -> %s\n' "$mod" "$cov" "$tot" "$pct" "$THRESHOLD" "$meets"

  if [[ "$meets" == "FAIL" ]]; then
    echo "::error::critical module '$mod' line coverage ${pct}% is below required ${THRESHOLD}%"
    fail=1
  fi
done

if [[ "$fail" -ne 0 ]]; then
  echo "Coverage gate FAILED."
  exit 1
fi
echo "Coverage gate passed: all critical modules >= ${THRESHOLD}%."
