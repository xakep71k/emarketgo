#!/bin/bash -xe
(cd emarket ; go build)
docker build -t emarket .
rm -f emarket/emarket
