#!/bin/bash

set -x

PACKAGE_DIR="$(pwd)/tmp-pkg"
trap "rm -rf $PACKAGE_DIR" EXIT

mkdir -p $PACKAGE_DIR/usr/bin \
	$PACKAGE_DIR/usr/share/kurmad \
	$PACKAGE_DIR/etc/kurmad \
	$PACKAGE_DIR/var/cache/kurmad/images \
	$PACKAGE_DIR/var/cache/kurmad/volumes \
	$PACKAGE_DIR/var/cache/kurmad/pods \
	$PACKAGE_DIR/var/run
cp ./bin/kurma-cli ./bin/kurmad $PACKAGE_DIR/usr/bin
cp ./build/release/base-config.yml $PACKAGE_DIR/etc/kurmad/config.yml
cp ./bin/kurma-api.aci \
	./bin/console.aci \
	./bin/kurma-upgrader.aci \
	./bin/stager-container.aci \
	./bin/busybox.aci \
	./bin/cni-netplugin.aci \
	$PACKAGE_DIR/usr/share/kurmad


docker pull kurma/debian-fpm:latest
docker pull kurma/centos-fpm:latest

if [ -n "$IN_DOCKER" ]; then
	docker run \
		   --rm \
		   --volumes-from $HOSTNAME \
		   -e KURMA_DIR=$(pwd) \
		   -e VERSION=$VERSION \
		   -e TARGET_INIT="systemd" \
		   kurma/debian-fpm:latest

	docker run \
		   --rm \
		   --volumes-from $HOSTNAME \
		   -e KURMA_DIR=$(pwd) \
		   -e VERSION=$VERSION \
		   -e TARGET_INIT="upstart" \
		   kurma/debian-fpm:latest

else
	docker run --rm \
		   -v $(pwd):"/kurma" \
		   -e KURMA_DIR="/kurma" \
		   -e VERSION=$VERSION \
		   -e TARGET_INIT="systemd" \
		   kurma/debian-fpm:latest

	docker run --rm \
		   -v $(pwd):"/kurma" \
		   -e KURMA_DIR="/kurma" \
		   -e VERSION=$VERSION \
		   -e TARGET_INIT="upstart" \
		   kurma/debian-fpm:latest
fi
