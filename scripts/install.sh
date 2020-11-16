#!/bin/sh
# shellcheck shell=dash

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
OS=$(uname)
ARCH=$(uname -m)
TMP_DIR="${REPO}_tmp"
DEFAULT_INSTALL_DIR=/usr/local/bin
INSTALL_DIR=${INSTALL_DIR:-$DEFAULT_INSTALL_DIR}
DEFAULT_VERSION=latest
VERSION=${VERSION:-$DEFAULT_VERSION}

# By default, list latest release
if [ "$VERSION" = latest ]; then
    VERSION=$(curl -s "https://api.github.com/repos/${USER}/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo "Failed to determine latest release version of ${USER}/${REPO}"
        exit 1
    fi
fi
FILENAME="${REPO}_${VERSION}_${OS}_${ARCH}.tar.gz"

# Download archive from GitHub Releases
if ! curl -sLO -w '' -f "https://github.com/${USER}/${REPO}/releases/download/v${VERSION}/${FILENAME}"; then
    echo "Failed to download ${USER}/${REPO} at version \"${VERSION}\""
    exit 1
fi

# Unpack into tmp directory
mkdir -p ${TMP_DIR}
if ! tar xzf "${FILENAME}" -C "${TMP_DIR}"; then
    echo "Failed to unpack ${FILENAME}"
    exit 1
fi

# Prompt for sudo if directory is not writable
MV_CMD="mv"
if [ ! -w "${INSTALL_DIR}" ]; then
    MV_CMD="sudo mv"
fi

# Install any executables
for f in "${TMP_DIR}"/*; do
    if [ -x "$f" ]; then
        if ! $MV_CMD "$f" "$INSTALL_DIR"; then
            echo "Failed to install ${INSTALL_DIR}/${f}"
            exit 1
        fi
    fi
done

# Cleanup
rm -f ${TMP_DIR}/*
rmdir ${TMP_DIR}
echo "${REPO} has been installed!"
