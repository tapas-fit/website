#!/bin/bash

source vars.public
source vars.secret

mkdir -p deploy-tmp
rsync -au --delete frontend/ deploy-tmp/
sed -i "s/TURNSTILE-SITEKEY/$sitekey_dev/" deploy-tmp/index.html
sed -i "s#BASE-URL#$baseurl_dev#" deploy-tmp/index.html
rsync -au --delete deploy-tmp/ /srv/http/tapas/

cd backend
./make.sh
./waitlist-linux ../dev-testdata localhost:6000 $secretkey_dev
