#!/bin/bash

set -e -x

cd $KURMA_DIR
/usr/local/bin/fpm \
	-C ./tmp-pkg \
	-f \
	-p ./resources \
	-v $VERSION \
	--after-install ./build/docker/rpm/startup.sh \
	-s dir \
	-t rpm \
	-n "kurmad-systemd" \
	.
