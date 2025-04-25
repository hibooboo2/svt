
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
    git push --tag $2 || git tag -d $1
}

alias gtags='git describe --tags | cut -f 1-2 -d "-" '

ntag(){
    pull development
    TAGNAME=$(git describe --tags --match v*.*.* | cut -f 1 -d "-"| xargs -n 1 svt -mode dev v0.0.1)
    echo $TAGNAME
}

gntag(){
    ntag
    tagIfNoTag $TAGNAME development $1
}

utag(){
    pull staging
    TAGNAME=$(git describe --tags --match uat-* | cut -f 1-2 -d "-"| xargs -n 1 svt -mode uat)
    echo $TAGNAME
}

gutag(){
    utag
    tagIfNoTag $TAGNAME staging $1
}

ptag(){
    pull production
    TAGNAME=$(git describe --tags --match r*.* | cut -f 1 -d "-"| xargs -n 1 svt -mode prod)
    echo $TAGNAME
}

gptag(){
    ptag
    tagIfNoTag $TAGNAME production $1
}