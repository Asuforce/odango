#!/bin/sh

set -e

if [[ $(git status -s) != "" ]]; then
    echo 2>&1 "git is currently in a dirty state"
    exit 1
fi

current_version=$(gobump show -r)

echo "current version: $current_version"
read -p "input next version: " next_version

echo "--> Bumping version $next_version"
gobump set "$next_version" -w
echo "--> Generating CHANGELOG"
ghch -w -N "v$next_version"

git commit -am "Bump version $next_version"
git tag "v$next_version"
git push && git push --tags
