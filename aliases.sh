
pull(){
    git fetch --all --tags
    git checkout $1
    git pull
}

# Source this file
tagIfNoTag(){
    set +x
    CURRENT_BRANCH=$(git branch --show-current)

    if [[ "$CURRENT_BRANCH" != "$2" ]]; then
        echo "Error: Not on the '$2' branch."
        return
    fi

    TAG="$1"

    # Extract last numeric segment and decrement
    if [[ "$TAG" =~ ^(.*[^0-9])([0-9]+)$ ]]; then
        PREFIX="${BASH_REMATCH[1]}"
        NUM="${BASH_REMATCH[2]}"
        PREV_TAG="${PREFIX}$((NUM - 1))"
    else
        PREV_TAG=""
    fi

    TAGS_ON_COMMIT=$(git tag --points-at HEAD)

    if echo "$TAGS_ON_COMMIT" | grep -qx "$TAG" || \
       { [[ -n "$PREV_TAG" ]] && echo "$TAGS_ON_COMMIT" | grep -qx "$PREV_TAG"; }; then
        echo "Commit already tagged with: $TAG or $PREV_TAG"
    else
        echo "No matching tag found. Tagging with: $TAG"
        git tag "$TAG"
        git push --tag "$3" 
        # || git tag -d "$TAG"
    fi

    set -x
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