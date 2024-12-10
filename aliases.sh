
pull(){
    git fetch --all --tags
    git checkout $1
    git pull
}

# Source this file
tagIfNoTag(){
    if [ "$(git tag --points-at HEAD | wc -l)" -eq 0 ]; then
        echo "No tags found on the current commit tagging with: $1"
        git tag $1
    else
        echo "Commit already tagged: $(git tag --points-at HEAD)"
    fi
    git push --tag
}

alias gtags='git describe --tags | cut -f 1-2 -d "-" '

gntag(){
    pull development
    TAGNAME=$(git describe --tags --match v*.*.* | cut -f 1 -d "-"| xargs -n 1 svt -mode dev v0.0.1)
    echo $TAGNAME
    tagIfNoTag $TAGNAME
}
gutag(){
    pull staging
    TAGNAME=$(git describe --tags --match uat-* | cut -f 1-2 -d "-"| xargs -n 1 svt -mode uat)
    echo $TAGNAME
    tagIfNoTag $TAGNAME
}
gptag(){
    pull production
    TAGNAME=$(git describe --tags --match r* | cut -f 1 -d "-"| xargs -n 1 svt -mode prod)
    echo $TAGNAME
    tagIfNoTag $TAGNAME
}