# Source this file
alias gtags='git describe --tags | cut -f 1-2 -d "-" '
alias gntag='git fetch --all --tags;git checkout development;git pull; git tag `git describe --tags --match v*.*.* | cut -f 1 -d "-"| xargs -n 1 svt -mode dev v0.0.1 ` && git push --tag'
alias gutag='git fetch --all --tags;git checkout staging;git pull; git tag `git describe --tags --match uat | xargs -n 1 svt -mode uat ` && git push --tag'
