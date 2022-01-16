#!/bin/bash

echo "Enter next tag"
read tag
echo -e "" >> go.sum
git add go.sum
# MESSAGE="attack,sign,verify"
MESSAGE="attack"
git commit -m $MESSAGE
git tag -a $tag -m $MESSAGE
git push origin main
git push --tags