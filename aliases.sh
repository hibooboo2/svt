# Source this file
alias gtags='git describe --tags | cut -f 1-2 -d "-" '
alias gntag='git fetch --all --tags;git checkout development;git pull; git tag `git describe --tags --match v*.*.* | svt -last-tag -mode dev | cut -f 1 -d "-"| xargs -n 1 svt -mode dev v0.0.1 ` && git push --tag'
alias gutag='git fetch --all --tags;git checkout staging;git pull; git tag `git describe --tags --match uat-* | svt -last-tag -mode uat | cut -f 1-2 -d "-"| xargs -n 1 svt -mode uat ` && git push --tag'
alias gptag='git fetch --all --tags;git checkout production;git pull; git tag `git describe --tags --match r* | svt -last-tag -mode prod | cut -f 1 -d "-"| xargs -n 1 svt -mode prod ` && git push --tag'
