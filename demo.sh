#!/bin/bash



function NewTag() {
    #get highest tag number
    VERSION=`git describe --abbrev=0 --tags`

    #replace . with space so can split into an array
    VERSION_BITS=(${VERSION//./ })

    #get number parts and increase last one by 1
    VNUM1=${VERSION_BITS[0]}
    VNUM2=${VERSION_BITS[1]}
    VNUM3=${VERSION_BITS[2]}
    VNUM1=`echo $VNUM1 | sed 's/v//'`

    # Check for #major or #minor in commit message and increment the relevant version number
    MAJOR=`git log --format=%B -n 1 HEAD | grep '#major'`
    MINOR=`git log --format=%B -n 1 HEAD | grep '#minor'`

    if [ "$MAJOR" ]; then
        echo "Update major version"
        VNUM1=$((VNUM1+1))
        VNUM2=0
        VNUM3=0
    elif [ "$MINOR" ]; then
        echo "Update minor version"
        VNUM2=$((VNUM2+1))
        VNUM3=0
    else
        echo "Update patch version"
        VNUM3=$((VNUM3+1))
    fi
    NEW_TAG="v$VNUM1.$VNUM2.$VNUM3"
}

function AddTaggedCommit() {
    MESSAGE=$1
    NewTag
    echo "PREV:$VERSION, NEW: $NEW_TAG"
    sed "s/$VERSION/$NEW_TAG/g" .valint.yaml
    git commit -m $MESSAGE

    GIT_COMMIT=`git rev-parse HEAD`
    NEEDS_TAG=`git describe --contains $GIT_COMMIT`
    if [ -z "$NEEDS_TAG" ]; then
        git tag -a $NEW_TAG -m $MESSAGE
        echo "Tagging -a $NEW_TAG -m $MESSAGE "
    else
        echo "Already a tag on this commit"
    fi
}

function UpdateFile() {
    MESSAGE=$1
    FILE=$2
    echo -e "" >> $FILE
    git add $FILE
    git add .valint.yaml
    AddTaggedCommit $MESSAGE
    git push origin main
    git push --tags
}


UpdateFile "bump-version - demo release" "README.md"
# echo "Updating $VERSION to $NEW_TAG"