#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Bump the named shared libraries to master.

LIB="${LIB:-}"
if [[ "${LIB}" == "" ]]; then
  >&2 echo "error: LIB not set. Please supply one of {api,build-machinery-go,apiserver-library-go,client-go,library-go}."
  exit 1
fi

REPO="${REPO:-}"
if [[ "${REPO}" == "" ]]; then
  >&2 echo "error: REPO not set."
  exit 1
fi

BRANCH="${BRANCH:-}"
if [[ "${BRANCH}" == "" ]]; then
  >&2 echo "error: BRANCH not set"
  exit 1
fi

MESSAGE="${MESSAGE:-}"
if [[ "${MESSAGE}" == "" ]]; then
  >&2 echo "error: MESSAGE not set"
  exit 1
fi

TITLE="${TITLE:-}"
if [[ "${TITLE}" == "" ]]; then
  >&2 echo "error: TITLE not set"
  exit 1
fi

GH_USER_ID=${GH_USER_ID}
if [[ "${GH_USER_ID}" == "" ]]; then
  >&2 echo "error: GH_USER_ID not set"
  exit 1
fi

ROOT_PATH="$( dirname "${BASH_SOURCE[0]}" )/.."
REPO_PATH="${ROOT_PATH}/repos/${REPO}"

echo "Updating ${REPO_PATH}"
pushd "${REPO_PATH}" > /dev/null
  if [[ ! -f "go.mod" ]]; then
    >&2 echo "warning: this repo is not using go modules"
    # Skip, don't error out.
    exit 0
  fi

  # Add ssh-based remote if not already present
  if [[ "$( git remote | grep upstream)" == "" ]]; then
    git remote rename origin upstream
    git remote add origin "https://github.com/${GH_USER_ID}/${REPO}"
    git remote set-url --push origin "git@github.com:${GH_USER_ID}/${REPO}.git"
  fi

  # Ensure a fork exists
  gh repo fork --remote=false

  # Bump the library
  git co -b "${BRANCH}"
  go get "github.com/openshift/${LIB}@master"
  go mod tidy
  go mod vendor
  git add .
  git ci -m "${MESSAGE}"

  # Create a PR for the bump
  gh pr create --title="${TITLE}" --body="${BODY}"
popd > /dev/null
