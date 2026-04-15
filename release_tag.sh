#!/usr/bin/env bash

set -euo pipefail

BUMP_TYPE="patch"
ALLOW_NON_MAIN="false"
REMOTE="origin"
START_TAG="v0.1.0"

usage() {
  cat <<'EOF'
Usage: ./release_tag.sh [--major|--minor|--patch] [--allow-non-main] [--remote <name>]

Automatically creates and pushes the next semantic version tag.

Options:
  --major            Increment major version
  --minor            Increment minor version
  --patch            Increment patch version (default)
  --allow-non-main   Allow running from branches other than main
  --remote <name>    Git remote to push (default: origin)
  --help             Show this help message
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --major)
      BUMP_TYPE="major"
      shift
      ;;
    --minor)
      BUMP_TYPE="minor"
      shift
      ;;
    --patch)
      BUMP_TYPE="patch"
      shift
      ;;
    --allow-non-main)
      ALLOW_NON_MAIN="true"
      shift
      ;;
    --remote)
      if [[ $# -lt 2 ]]; then
        echo "Missing value for --remote" >&2
        exit 1
      fi
      REMOTE="$2"
      shift 2
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage
      exit 1
      ;;
  esac
done

current_branch="$(git rev-parse --abbrev-ref HEAD)"
if [[ "$ALLOW_NON_MAIN" != "true" && "$current_branch" != "main" ]]; then
  echo "Release tags must be created from main (current: $current_branch)." >&2
  echo "Use --allow-non-main to override." >&2
  exit 1
fi

if [[ -n "$(git status --porcelain)" ]]; then
  echo "Working tree is not clean. Commit/stash changes before tagging." >&2
  exit 1
fi

git fetch --tags "$REMOTE"

last_tag="$(git tag -l 'v[0-9]*.[0-9]*.[0-9]*' --sort=-v:refname | head -n 1)"
if [[ -z "$last_tag" ]]; then
  next_tag="$START_TAG"
else
  version="${last_tag#v}"
  IFS='.' read -r major minor patch <<< "$version"

  case "$BUMP_TYPE" in
    major)
      major=$((major + 1))
      minor=0
      patch=0
      ;;
    minor)
      minor=$((minor + 1))
      patch=0
      ;;
    patch)
      patch=$((patch + 1))
      ;;
  esac

  next_tag="v${major}.${minor}.${patch}"
fi

if git rev-parse "$next_tag" >/dev/null 2>&1; then
  echo "Tag already exists: $next_tag" >&2
  exit 1
fi

echo "Last tag: ${last_tag:-<none>}"
echo "Next tag: $next_tag"

git tag "$next_tag"
git push "$REMOTE" "$next_tag"

echo "Tag pushed successfully: $next_tag"
