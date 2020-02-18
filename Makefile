
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
# No, I'm not entirely sure this is 105% correct, but it's in the right
# neighbourhood.

PREFIX=/usr/local

GIT_DESCRIBE:=$(shell git describe --always --tag --dirty)
VERSION_NO=$(shell echo ${GIT_DESCRIBE} | sed s/[v-]//g)
OS:=$(shell uname -s | tr A-Z a-z)
ARCH:=$(shell uname -m)

skogul: $(wildcard *.go */*.go */*/*.go)
	go build -ldflags "-X main.versionNo=$V" -o skogul ./cmd/skogul

docs/skogul.rst: skogul
	./skogul -make-man > $@

skogul.1: docs/skogul.rst
	rst2man < $< > $@

# Extract release notes - used by drone
notes:
	./build/release-notes.sh > notes

# MAGIC - for creating directories and not littering stdout with redundant
# mkdir -p's
%/:
	mkdir -p $@

all: skogul skogul.1 docs/skogul.rst

install: skogul skogul.1 docs/skogul.rst
	install -D -m 0755 skogul ${DESTDIR}${PREFIX}/bin/skogul
	install -D -m 0644 skogul.1 ${DESTDIR}${PREFIX}/share/man/man1/skogul.1
	install -D -m 0644 docs/examples/default.json ${DESTDIR}/etc/skogul/default.json
	cd docs; \
	find -type f -exec install -D -m 0644 {} ${DESTDIR}${PREFIX}/share/doc/skogul/{} \;
	install -D -m 0644 README.rst LICENSE -t ${DESTDIR}${PREFIX}/share/doc/skogul/


rpm-prep/SPECS/skogul.spec: build/redhat-skogul.spec.in | rpm-prep/SPECS/
	cat $< | sed "s/xxVxx/${GIT_DESCRIBE}/g; s/xxARCHxx/${ARCH}/g; s/xxVERSION_NOxx/${VERSION_NO}/g" > $@

rpm: rpm-prep/SPECS/skogul.spec | rpm-prep/BUILDROOT/ rpm-prep/RPMS/ rpm-prep/SPECS/ rpm-prep/SRPMS/
	# Hacky as heck, and creates a tight coupling between makefile and
	# spec. But I just can't be bothered to fix this right now.
	test -h rpm-prep/BUILD || ln -s ./ rpm-prep/BUILD
	
	# Taken from CentOS Linux release 7.6.1810 (Core)
	cd rpm-prep; \
	DEFAULT_UNIT_DIR=/usr/lib/systemd/system ;\
	RPM_UNIT_DIR=$$(rpm --eval $%{_unitdir}) ;\
	if [ "$${RPM_UNIT_DIR}" = "$%{_unitdir}" ]; then \
	    echo "_unitdir not set, setting _unitdir to $$DEFAULT_UNIT_DIR"; \
	    rpmbuild --bb \
		--define "_rpmdir $$(pwd)" \
		--define "_sourcedir $$(pwd)/SOURCES" \
		--define "_topdir $$(pwd)" \
		--define "_unitdir $$DEFAULT_UNIT_DIR" \
		--buildroot "$$(pwd)/BUILDROOT" \
		SPECS/skogul.spec; \
	else \
	    rpmbuild --bb \
		    --define "_rpmdir $$(pwd)" \
		    --define "_sourcedir $$(pwd)/SOURCES" \
		    --define "_topdir $$(pwd)" \
		    --buildroot "$$(pwd)/BUILDROOT" \
		    SPECS/skogul.spec ;\
	fi
	cp rpm-prep/x86_64/* .

test:
	go test -short ./...

bench:
	go test -run ^Bench -benchtime 1s -bench Bench ./... | grep Benchmark

clean:
	-rm -fr dist
	-rm -fr rpm-prep
	-rm -f skogul
	-rm -f docs/skogul.rst
	-rm -f skogul.1

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

