# Source this file
tagIfNoTag(){
    git fetch --all --tags
    git checkout development
    git pull
    if [ "$(git tag --points-at HEAD | wc -l)" -eq 0 ]; then
        echo "No tags found on the current commit."
        git tag $1
        git push --tag
    else
        echo "Commit already tagged"
    fi
}

alias gtags='git describe --tags | cut -f 1-2 -d "-" '

gntag(){
    TAGNAME=$(git describe --tags --match v*.*.* | cut -f 1 -d "-"| xargs -n 1 svt -mode dev v0.0.1)
    echo $TAGNAME
    tagIfNoTag $TAGNAME
}
gutag(){
    TAGNAME=$(git describe --tags --match uat-* | cut -f 1-2 -d "-"| xargs -n 1 svt -mode uat)
    echo $TAGNAME
    tagIfNoTag $TAGNAME
}
gptag(){
    TAGNAME=$(git describe --tags --match r* | cut -f 1 -d "-"| xargs -n 1 svt -mode prod)
    echo $TAGNAME
    tagIfNoTag $TAGNAME
}