#!/bin/sh

set -e

ver=v$(gobump show -r)
make crossbuild
ghr -username Asuforce -replace ${ver} dist/${ver}
