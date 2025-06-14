#!/bin/bash
# get script's work dir: https://www.ostricher.com/2014/10/the-right-way-to-get-the-directory-of-a-bash-script/
cd "$(dirname "${BASH_SOURCE[0]}")"
./fieldday fd2025.db
