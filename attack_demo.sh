#!/bin/bash

echo "Enter next tag"
read tag
echo -e "" >> go.sum
git add go.sum
git commit -m "attack,sign,verify"
git tag -a $tag -m "attack,sign,verify"
git push origin main
git push --tags