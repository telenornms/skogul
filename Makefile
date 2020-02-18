
# Makefile-lol 101:
#  PREFIX is typically overridden locally. Should default to /usr/local,
#  and be set to /usr explicitly. It is used for most items, but not all.
#  It indicates where the program will be installed on the target system.
#
# DESTDIR is used to do installations during packaging etc. We will install
# all files under DESTDIR$PREFIX, but any internal references will ignore
# DESTDIR. E.g.: rpmbuild intends for the target systems to install to
# /usr, so PREFIX should be /usr, but during the RPM build process, DESTDIR
# is set to rpm-prep/BUILDROOT, so you get
# rpm-prep/BUILDROOT/usr/bin/skogul, etc.
#
# DESTDIR is empty/undefined by default.
#
# No, I'm not entirely sure this is 105% correct, but it's in the right
# neighbourhood.

PREFIX=/usr/local
DOCDIR=${PREFIX}/share/doc/skogul

GIT_DESCRIBE:=$(shell git describe --always --tag --dirty)
VERSION_NO=$(shell echo ${GIT_DESCRIBE} | sed s/[v-]//g)
OS:=$(shell uname -s | tr A-Z a-z)
ARCH:=$(shell uname -m)

skogul: $(wildcard *.go */*.go */*/*.go)
	@echo 🤸 go build !
	@go build -ldflags "-X main.versionNo=$V" -o skogul ./cmd/skogul

docs/skogul.rst: skogul
	@echo 😽 Generating documentation $@
	@./skogul -make-man > $@

skogul.1: docs/skogul.rst
	@echo 🎢 Generating man-file $@
	@rst2man < $< > $@

notes: docs/NEWS
	@echo ⛲ Extracting release notes.
	@./build/release-notes.sh > notes

# MAGIC (I hate noisy Make-runs)
%/:
	@mkdir -p $@

all: skogul skogul.1 docs/skogul.rst

install: skogul skogul.1 docs/skogul.rst
	@echo 🙅 Installing
	@install -D -m 0755 skogul ${DESTDIR}${PREFIX}/bin/skogul
	@install -D -m 0644 skogul.1 ${DESTDIR}${PREFIX}/share/man/man1/skogul.1
	@install -D -m 0644 docs/examples/default.json ${DESTDIR}/etc/skogul/default.json
	@cd docs; \
	find -type f -exec install -D -m 0644 {} ${DESTDIR}${DOCDIR}/{} \;
	@install -D -m 0644 README.rst LICENSE -t ${DESTDIR}${DOCDIR}/


# Any complaints on this macro-substitution without patches and I introduce m4.
rpm-prep/SPECS/skogul.spec: build/redhat-skogul.spec.in | rpm-prep/SPECS/
	@echo  ❕Building spec-file
	@cat $< | sed "s/xxVxx/${GIT_DESCRIBE}/g; s/xxARCHxx/${ARCH}/g; s/xxVERSION_NOxx/${VERSION_NO}/g" > $@
	@which dpkg >/dev/null && { echo 🆒 Adding debian-workaround for rpm build; sed -i 's/^BuildReq/\#Debian hack, auto-commented out: BuildReq/g' $@; }

# Build RPM. The spec has a blank %prep, so it assumes sources are already
# available. This isn't perfect, since it creates a tight coupling between
# Makefile and specfile, but it isn't all that bad either, since it allows
# building with minimal redundant effort, and without having to commit to
# git.
rpm: rpm-prep/SPECS/skogul.spec | rpm-prep/BUILDROOT/
	@echo 🎇 Triggering huge-as-heck rpm build
	@DEFAULT_UNIT_DIR=/usr/lib/systemd/system ;\
	RPM_UNIT_DIR=$$(rpm --eval $%{_unitdir}) ;\
	if [ "$${RPM_UNIT_DIR}" = "$%{_unitdir}" ]; then \
	    echo "😭 _unitdir not set, setting _unitdir to $$DEFAULT_UNIT_DIR"; \
	    rpmbuild --quiet --bb \
	    	--build-in-place \
		--define "_rpmdir $$(pwd)" \
		--define "_topdir $$(pwd)" \
		--define "_unitdir $$DEFAULT_UNIT_DIR" \
		--buildroot "$$(pwd)/rpm-prep/BUILDROOT" \
		rpm-prep/SPECS/skogul.spec; \
	else \
	    rpmbuild --quiet --bb \
	    	--build-in-place \
		--define "_rpmdir $$(pwd)" \
		--define "_topdir $$(pwd)" \
		--buildroot "$$(pwd)/rpm-prep/BUILDROOT" \
		rpm-prep/SPECS/skogul.spec; \
	fi
	@cp x86_64/skogul-${VERSION_NO}-1.x86_64.rpm .
	@echo ⭐ RPM built: ./skogul-${VERSION_NO}-1.x86_64.rpm


test:
	go test -short ./...

bench:
	go test -run ^Bench -benchtime 1s -bench Bench ./... | grep Benchmark

clean:
	@echo 💩Cleaning up
	@-rm -fr dist
	@-rm -fr rpm-prep
	@-rm -f skogul
	@-rm -f docs/skogul.rst
	@-rm -f skogul.1
	@-rm -f *.rpm

help:
	@echo "Several targets exist:"
	@echo 
	@echo " - skogul - build the binary"
	@echo " - all - build binary and documentation"
	@echo " - install - install binary and docs. Honors PREFIX, default prefix: ${PREFIX}"
	@echo " - rpm - build RPM"
	@echo " - clean - remove build crap"
	@echo " - test / bench - run go test, with and without benchmarks "
	@echo "                  note that this uses "-short" to avoid mysql/postgres dependencies. "

.PHONY: clean test bench help install rpm
