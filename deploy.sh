#!/bin/bash
set -e

mkdir -p bin
GOOS=linux GOARCH=amd64 go build -o ./bin/frontend -tags prod .
rm -f content/assets/script-*.js content/assets/style-*.css
mend -f content/mend.json -o content/mend-versions.json content
cd admin
ansible-playbook -i prod deploy.yml

