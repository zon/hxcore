#!/bin/bash

VERSION=$(cat ./version)

git tag -a $VERSION -m $VERSION
git push origin $VERSION