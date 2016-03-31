#!/usr/bin/env bash
set -e
# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that dir because we expect that
cd $DIR
source app.env

: ${APP_NAME:?Not set. Ex. APP_NAME=findami}
: ${VERSION:?Not set. Ex. VERSION=0.0.1}

# Tag, unless told not to
if [ -z $NOTAG ]; then
  echo "==> Tagging..."
  git commit --allow-empty -a -m "Cut version $VERSION"
  git tag -a -m "Version $VERSION" "v${VERSION}" $RELBRANCH
  git push origin "v${VERSION}"
fi

# Zip all the files
rm -rf ./pkg/dist
mkdir -p ./pkg/dist
for PLATFORM in $(find ./pkg -mindepth 1 -maxdepth 1 -type d); do
  OSARCH=$(basename ${PLATFORM})

  if [ $OSARCH = "dist" ]; then
    continue
  fi

  echo "--> ${OSARCH}"
  pushd $PLATFORM >/dev/null 2>&1
  zip ../dist/${APP_NAME}_${VERSION}_${OSARCH}.zip ./*
  popd >/dev/null 2>&1
done
