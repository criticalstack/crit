#!/bin/bash


# wait until the dev bin exists, then copy into right place, automatically fails after 20 tries, 5s each
aws s3api wait object-exists --bucket ${s3_bucket_name} --key bin/crit
aws s3 cp s3://${s3_bucket_name}/bin/crit /usr/bin/crit
chmod +x /usr/bin/crit

#####################################################################
# Retreive cluster certificate
#####################################################################

aws s3api wait object-exists --bucket ${s3_bucket_name} --key pki/ca.crt
aws s3 cp s3://${s3_bucket_name}/pki/ca.crt /etc/kubernetes/pki/ca.crt
