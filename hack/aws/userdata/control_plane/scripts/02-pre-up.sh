#!/bin/bash


# wait until the dev bin exists, then copy into right place, automatically fails after 20 tries, 5s each
aws s3api wait object-exists --bucket ${s3_bucket_name} --key bin/crit
aws s3 cp s3://${s3_bucket_name}/bin/crit /usr/bin/crit
chmod +x /usr/bin/crit

#####################################################################
# Retreive cluster secrets
#####################################################################

aws s3api wait object-exists --bucket ${s3_bucket_name} --key pki/ca.crt
aws s3api wait object-exists --bucket ${s3_bucket_name} --key pki/ca.key
aws s3api wait object-exists --bucket ${s3_bucket_name} --key pki/etcd/ca.crt
aws s3api wait object-exists --bucket ${s3_bucket_name} --key pki/etcd/ca.key
aws s3 cp --recursive s3://${s3_bucket_name}/pki/ /etc/kubernetes/pki/


#####################################################################
# Setup e2d
#####################################################################

export AWS_LOCAL_IPV4=$(curl -s http://169.254.169.254/latest/meta-data/local-ipv4)

# create host-specific certs
e2d pki gencerts \
  --ca-cert=/etc/kubernetes/pki/etcd/ca.crt \
  --ca-key=/etc/kubernetes/pki/etcd/ca.key \
  --hosts=$AWS_LOCAL_IPV4,${control_plane_endpoint} \
  --output-dir=/etc/kubernetes/pki/etcd

systemctl enable --now e2d.service
