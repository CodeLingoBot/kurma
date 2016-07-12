#!/bin/bash

set -e -x

PACKAGE_DIR="$(pwd)/tmp-pkg"
trap "rm -rf $PACKAGE_DIR" EXIT

mkdir -p $PACKAGE_DIR/usr/bin \
	$PACKAGE_DIR/usr/share/kurmad \
	$PACKAGE_DIR/etc/kurmad \
	$PACKAGE_DIR/var/cache/kurmad/images \
	$PACKAGE_DIR/var/cache/kurmad/volumes \
	$PACKAGE_DIR/var/cache/kurmad/pods \
	$PACKAGE_DIR/var/run \
	$PACKAGE_DIR/lib/systemd/system
cp ./bin/kurma-cli ./bin/kurmad $PACKAGE_DIR/usr/bin
cp ./build/release/base-config.yml $PACKAGE_DIR/etc/kurmad/config.yml
cp ./bin/kurma-api.aci \
	./bin/console.aci \
	./bin/stager-container.aci \
	./bin/busybox.aci \
	./bin/cni-netplugin.aci \
	$PACKAGE_DIR/usr/share/kurmad
cp ./build/release/kurmad.service $PACKAGE_DIR/lib/systemd/system

docker pull kurma/debian-fpm:latest
docker pull kurma/centos-fpm:latest

if [ -n "$IN_DOCKER" ]; then
	docker run --rm \
		   --volumes-from $HOSTNAME \
		   -e KURMA_DIR=$(pwd) \
		   -e VERSION=$VERSION \
		   kurma/centos-fpm:latest

else
	docker run --rm \
		   -v $(pwd):"/kurma" \
		   -e KURMA_DIR="/kurma" \
		   -e VERSION=$VERSION \
		   kurma/centos-fpm:latest
fi
