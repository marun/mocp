#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

RELEASE_IMAGE="${1:-}"

if [[ "${RELEASE_IMAGE}" == "" ]]; then
  >&2 echo "$0: usage <image>"
  exit 1
fi

# Synchronize the repos referenced by a release with submodules.

ROOT_PATH="$( dirname "${BASH_SOURCE[0]}" )/.."
REPOS_PATH="${ROOT_PATH}/repos"

RAW_REPO_URLS="$( oc adm release info "${RELEASE_IMAGE}" --commit-urls |\
 grep 'https://github.com' |\
 awk '{print $2}' |\
 sed -e 's+/commit.*++' |\
 uniq |\
 sort
)"
REPO_URLS=(${RAW_REPO_URLS})

pushd "${REPOS_PATH}"> /dev/null

for repo_url in "${REPO_URLS[@]}"; do
  repo_name="$( basename "${repo_url}" )"
  if [[ ! -d "${repo_name}" ]]; then
    git submodule add ${repo_url}
  fi
done

popd > /dev/null
