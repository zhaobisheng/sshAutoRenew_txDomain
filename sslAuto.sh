#!/bin/bash

SecretId="API_ID"
secretKey="API_KEY"

action=$1 

if [[ "$action" != "clean" ]]; then
    action="add"
    ./txDomain -a=$action $SecretId $secretKey $CERTBOT_DOMAIN "_acme-challenge" TXT $CERTBOT_VALIDATION  >>"/var/log/certd.log"
else [[ "$action" == "clean" ]]; then  
    ./txDomain -a=$action $SecretId $secretKey $CERTBOT_DOMAIN
fi
