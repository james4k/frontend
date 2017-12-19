#!/bin/bash
set -e

mkdir -p bin
GOOS=linux GOARCH=amd64 gb build -tags prod cmd/frontend
ls -t content/assets/script-*.js | awk 'NR>2' | xargs -L 1 rm
ls -t content/assets/style-*.css | awk 'NR>2' | xargs -L 1 rm
mend -f content/mend.json -o content/mend-versions.json content
cd admin
ansible-playbook -i prod deploy.yml

