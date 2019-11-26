#!/bin/bash
set -e

GIT_DESCRIBE="$(git describe --always --tag --dirty)"
OS=$(uname -s | tr A-Z a-z)
ARCH=$(uname -m)
V=${1:-unknown}
if [ "x${V}" = "xunknown" ] && [ "x${GIT_DESCRIBE}" != "x" ]; then
	V=${GIT_DESCRIBE}
fi

rm -rf rpm-prep
mkdir -p rpm-prep/BUILD
mkdir -p rpm-prep/BUILDROOT
mkdir -p rpm-prep/RPMS
mkdir -p rpm-prep/SOURCES
mkdir -p rpm-prep/SPECS
mkdir -p rpm-prep/SRPMS


# Copy required files to build dir
cp LICENSE rpm-prep
cp make-docs.sh rpm-prep/BUILDROOT
cp skogul.service rpm-prep/SOURCES


VERSION_NO=$(echo $V | sed s/v//)

cat <<EOF > rpm-prep/SPECS/skogul.spec
Name:           skogul
Version:        $VERSION_NO
Release:        1%{?dist}
Summary:        Skogul metric engine

Group:          telenornms
License:        LGPL-2.1
URL:            https://github.com/telenornms/skogul
Source0:        https://github.com/telenornms/skogul/archive/v%{version}.tar.gz


BuildArch:      $ARCH
# The build requirements. However, because we build
# in a docker container with ubuntu
# rpm doesn't know that the requirements are available.
#BuildRequires:  go >= 1.13, python-docutils


%description
Skogul metric engine

# Executable files require a build id; let's stop that
# https://github.com/rpm-software-management/rpm/issues/367
%undefine _missing_build_ids_terminate_build

%prep
%setup -q


%build
go build -o dist/%{name} ./cmd/%{name}
bash %{buildroot}/make-docs.sh
rm -f %{buildroot}/make-docs.sh


%install
mkdir -p %{buildroot}%{_bindir}
mkdir -p %{buildroot}%{_mandir}/man1
mkdir -p %{buildroot}%{_defaultdocdir}/%{name}-%{version}
mkdir -p %{buildroot}%{_datadir}/licenses/%{name}-%{version}
install -m 0755 dist/%{name} %{buildroot}%{_bindir}/%{name}
cp dist/share/man/man1/%{name}.1 %{buildroot}%{_mandir}/man1/%{name}.1
cp -r docs/* %{buildroot}%{_defaultdocdir}/%{name}-%{version}
mkdir -p %{buildroot}%{_unitdir}
install -m 0644 %{_sourcedir}/%{name}.service %{buildroot}%{_unitdir}/%{name}.service

%pre
getent group skogul >/dev/null || groupadd -r skogul
getent passwd skogul >/dev/null || \
       useradd -r -g skogul -d /var/lib/skogul -s /sbin/nologin \
               -c "Skogul metric collector" skogul
exit 0

%post
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service


%files
%license LICENSE
%{_bindir}/%{name}
%{_mandir}/man1/%{name}.1*
%docdir %{_defaultdocdir}/%{name}-%{version}
%{_defaultdocdir}/%{name}-%{version}
%{_unitdir}/%{name}.service



%changelog
EOF

wget -O rpm-prep/SOURCES/$V.tar.gz https://github.com/telenornms/skogul/archive/$V.tar.gz

echo "Building"

cd rpm-prep
rpmbuild --bb \
    --define "_rpmdir $(pwd)" \
    --define "_sourcedir $(pwd)/SOURCES" \
    --define "_topdir $(pwd)" \
    --buildroot "$(pwd)/BUILDROOT" \
    SPECS/skogul.spec
