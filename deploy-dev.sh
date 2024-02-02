#!/bin/bash

source vars.public
source vars.secret

mkdir -p deploy-tmp
rsync -a --delete frontend/ deploy-tmp/
sed -i "s/TURNSTILE-SITEKEY/$sitekey_dev/" deploy-tmp/index.html
sed -i "s#BASE-URL#$baseurl_dev#" deploy-tmp/index.html
rsync -a --delete deploy-tmp/ /srv/http/tapas/

cd backend
./make.sh
./waitlist-linux ../dev-testdata localhost:6000 $secretkey_dev
