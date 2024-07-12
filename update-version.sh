#!/bin/bash

# Hole den letzten Git-Tag
LATEST_TAG=$(git describe --tags --abbrev=0 | sed s/v//g)

# Überprüfe, ob ein Git-Tag gefunden wurde
if [ -z "$LATEST_TAG" ]; then
  echo "no git version found"
  exit 1
fi

# Aktualisiere die Version in der package.json
jq ".version = \"$LATEST_TAG\"" package.json > package.tmp.json && mv package.tmp.json package.json

# Überprüfe, ob jq erfolgreich war
if [ $? -ne 0 ]; then
  echo "Error while updating version in package.json."
  exit 1
fi

echo "Version in package.json updated $LATEST_TAG successfully."
