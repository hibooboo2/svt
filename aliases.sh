
pull(){
    git fetch --all --tags
    git checkout $1
    git pull
}

# Source this file
tagIfNoTag(){
    CURRENT_BRANCH=$(git branch --show-current)

    if [[ "$CURRENT_BRANCH" != "$2" ]]; then
        echo "Error: Not on the '$2' branch."
        return
    fi

    if [ "$(git tag --points-at HEAD | wc -l)" -eq 0 ]; then
        echo "No tags found on the current commit tagging with: $1"
        git tag $1
    else
        echo "Commit already tagged: $(git tag --points-at HEAD)"
    fi
    git push --tag $3 || git tag -d $1
}

alias gtags='git describe --tags | cut -f 1-2 -d "-" '

ntag(){
    git fetch --all --tags
    TAGNAME=$(git describe --tags --match v*.*.* | cut -f 1 -d "-"| xargs -n 1 svt -mode dev v0.0.1)
    echo $TAGNAME
}

gntag(){
    pull development
    ntag
    tagIfNoTag $TAGNAME development $1
}

ttag() {
    git fetch --all --tags

    # Get the latest tag matching v*.*.*
    BASE_TAG=$(git describe --tags --match "v*.*.*" --abbrev=0)

    # Check if the latest tag already has "-test" suffix
    if [[ "$BASE_TAG" == *-test* ]]; then
        # Extract base version (before "-test") and current test number
        MAIN_TAG="${BASE_TAG%-test*}"
        CURRENT_NUM=$(echo "$BASE_TAG" | sed -n 's/.*-test\([0-9]*\)$/\1/p')
        if [[ -z "$CURRENT_NUM" ]]; then
            CURRENT_NUM=0
        fi
        NEXT_NUM=$((CURRENT_NUM + 1))
        TAGNAME="${MAIN_TAG}-test${NEXT_NUM}"
    else
        TAGNAME="${BASE_TAG}-test1"
    fi

    echo "$TAGNAME"
}

gttag(){
    git fetch --all
    ttag
    echo $TAGNAME
    tagIfNoTag $TAGNAME $(git branch --show-current) $1
}


utag(){
    git fetch --all --tags
    TAGNAME=$(git describe --tags --match uat-* | cut -f 1-2 -d "-"| xargs -n 1 svt -mode uat)
    echo $TAGNAME
}

gutag(){
    pull staging
    utag
    tagIfNoTag $TAGNAME staging $1
}

ptag(){
    git fetch --all --tags
    TAGNAME=$(git describe --tags --match r*.* | cut -f 1 -d "-"| xargs -n 1 svt -mode prod)
    echo $TAGNAME
}

gptag(){
    pull production
    ptag
    tagIfNoTag $TAGNAME production $1
}

gmtag(){
    pull main
    ntag
    tagIfNoTag $TAGNAME main $1
}