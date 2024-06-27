#!/bin/bash
set -e -o pipefail
cwd=$(pwd)
pushd "$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)" 1>/dev/null || exit
pushd gen_tfmoduleindex 1>/dev/null || exit

go run "./gen_tfmoduleindex.go" -dir="$cwd"
