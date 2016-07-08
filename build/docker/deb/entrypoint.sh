#!/bin/bash

set -e -x

if [ "$TARGET_INIT" = "systemd" ]; then
	cd $KURMA_DIR
	fpm -f \
		-C ./tmp-pkg \
		-p ./resources \
		-v $VERSION \
		--deb-systemd ./build/release/kurmad.service \
		-s dir \
		-t deb \
		-n "kurmad-systemd" \
		.
else
	cd $KURMA_DIR
	fpm -f \
		-C ./tmp-pkg \
		-p ./resources \
		-v $VERSION \
		--deb-upstart ./build/release/kurmad.conf \
		--after-install ./build/docker/deb/upstart-postinst.sh \
		-s dir \
		-t deb \
		-n "kurmad-upstart" \
		.
fi
