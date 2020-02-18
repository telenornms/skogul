#!/bin/bash

echo "Building"

cd rpm-prep

# Taken from CentOS Linux release 7.6.1810 (Core)
DEFAULT_UNIT_DIR=/usr/lib/systemd/system
RPM_UNIT_DIR=$(rpm --eval %{_unitdir})
if [ "${RPM_UNIT_DIR}" = "%{_unitdir}" ]; then
    echo "_unitdir not set, setting _unitdir to $DEFAULT_UNIT_DIR"

    rpmbuild --bb \
        --define "_rpmdir $(pwd)" \
        --define "_sourcedir $(pwd)/SOURCES" \
        --define "_topdir $(pwd)" \
        --define "_unitdir $DEFAULT_UNIT_DIR" \
        --buildroot "$(pwd)/BUILDROOT" \
        SPECS/skogul.spec
else
    rpmbuild --bb \
            --define "_rpmdir $(pwd)" \
            --define "_sourcedir $(pwd)/SOURCES" \
            --define "_topdir $(pwd)" \
            --buildroot "$(pwd)/BUILDROOT" \
            SPECS/skogul.spec
fi
