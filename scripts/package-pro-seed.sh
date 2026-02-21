#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MANIFEST_PATH="${1:-$ROOT_DIR/pro/pro-seed-manifest.txt}"
OUTPUT_DIR="${2:-$ROOT_DIR/pro/archives}"

if [[ ! -f "$MANIFEST_PATH" ]]; then
  echo "manifest not found: $MANIFEST_PATH" >&2
  exit 1
fi

if ! command -v zip >/dev/null 2>&1; then
  echo "zip command not found" >&2
  exit 1
fi

mkdir -p "$OUTPUT_DIR"
STAMP="$(date +%Y%m%d_%H%M%S)"
STAGE_DIR="$(mktemp -d "${TMPDIR:-/tmp}/pro-seed.XXXXXX")"
ZIP_PATH="$OUTPUT_DIR/pro-seed-$STAMP.zip"
COPIED=0

cleanup() {
  rm -rf "$STAGE_DIR"
}
trap cleanup EXIT

while IFS= read -r raw || [[ -n "$raw" ]]; do
  line="${raw%%#*}"
  line="$(printf '%s' "$line" | tr -d '\r' | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')"
  if [[ -z "$line" ]]; then
    continue
  fi

  src="$ROOT_DIR/$line"
  if [[ ! -f "$src" ]]; then
    echo "skip missing file: $line" >&2
    continue
  fi

  dst="$STAGE_DIR/$line"
  mkdir -p "$(dirname "$dst")"
  cp "$src" "$dst"
  COPIED=$((COPIED + 1))
done < "$MANIFEST_PATH"

if [[ "$COPIED" -eq 0 ]]; then
  echo "no files copied, archive not created" >&2
  exit 1
fi

(cd "$STAGE_DIR" && zip -qr "$ZIP_PATH" .)
echo "created: $ZIP_PATH"
echo "files: $COPIED"
