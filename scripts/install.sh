#!/bin/sh

# Copyright 2020 Critical Stack, LLC
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

USER="criticalstack"
REPO="crit"
VERSION=$(curl -s "https://api.github.com/repos/${USER}/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')
OS=$(uname)
ARCH=$(uname -m)
FILENAME="${REPO}_${VERSION}_${OS}_${ARCH}.tar.gz"
TMP_DIR="${REPO}_tmp"
INSTALL_DIR=/usr/local/bin

# Download archive from GitHub Releases
curl -sLO -w '' "https://github.com/${USER}/${REPO}/releases/download/v${VERSION}/${FILENAME}"

if [ $? -ne 0 ]; then
    echo "Failed to download ${USER}/${REPO}"
    exit 1
fi

# Unpack into tmp directory
mkdir -p ${TMP_DIR}
tar xzf ${FILENAME} -C ${TMP_DIR}

if [ $? -ne 0 ]; then
    echo "Failed to unpack ${FILENAME}"
    exit 1
fi

# Install any executables
for f in ${TMP_DIR}/*; do
    if [ -x $f ]; then
        sudo mv $f $INSTALL_DIR
    fi
done

# Cleanup
rm -f ${TMP_DIR}/*
rmdir ${TMP_DIR}
echo "${REPO} has been installed!"
