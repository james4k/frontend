#!/bin/bash
set -e

mkdir -p bin
GOOS=linux GOARCH=amd64 go build -o ./bin/frontend -tags prod .
mend -f content/mend.json -o content/mend-versions.json content
cd admin
ansible-playbook -i prod deploy.yml

