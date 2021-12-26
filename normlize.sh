#!/bin/bash

echo "Enter next tag"
read tag
echo -e "" >> go.sum
git add go.sum
git commit -m "normlize"
git tag -a $tag -m "normlize"
git push origin main
git push --tags