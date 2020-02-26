# PREFIX is the prefix on the targetsystem
# DESTDIR can be used to prefix ALL paths, e.g., to do a dummy-install in a
# fake root dir, e.g., for building packages. Users mainly want PREFIX

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
	@./build/release-notes.sh ${GIT_DESCRIBE} > notes

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
	@which dpkg >/dev/null 2>&1 && { echo 🆒 Adding debian-workaround for rpm build; sed -i 's/^BuildReq/\#Debian hack, auto-commented out: BuildReq/g' $@; } || true

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
	        --nodebuginfo \
	    	--build-in-place \
		--define "_rpmdir $$(pwd)" \
		--define "_topdir $$(pwd)" \
		--define "_unitdir $$DEFAULT_UNIT_DIR" \
		--buildroot "$$(pwd)/rpm-prep/BUILDROOT" \
		rpm-prep/SPECS/skogul.spec; \
	else \
	    rpmbuild --quiet --bb \
	        --nodebuginfo \
	    	--build-in-place \
		--define "_rpmdir $$(pwd)" \
		--define "_topdir $$(pwd)" \
		--buildroot "$$(pwd)/rpm-prep/BUILDROOT" \
		rpm-prep/SPECS/skogul.spec; \
	fi
	@cp x86_64/skogul-${VERSION_NO}-1.x86_64.rpm .
	@echo ⭐ RPM built: ./skogul-${VERSION_NO}-1.x86_64.rpm

check: test fmtcheck vet

vet:
	@echo 🔬 Vetting code
	@go vet ./...

fmtcheck:
	@echo 🦉 Checking format with gofmt -d -s
	@if [ "x$$(find -name '*.go' -not -wholename './gen/*' -exec gofmt -d -s {} +)" != "x" ]; then find -name '*.go' -not -wholename './gen/*' -exec gofmt -d -s {} +; exit 1; fi

fmtfix:
	@echo 🎨 Fixing formating
	@find -name '*.go' -not -wholename './gen/*' -exec gofmt -d -s -w {} +

test:
	@echo 🧐 Testing, without SQL-tests
	@go test -short ./...

bench:
	@echo 🏋 Benchmarking
	@go test -run ^Bench -benchtime 1s -bench Bench ./... | grep Benchmark

covergui:
	@echo 🧠 Testing, with coverage analysis
	@go test -short -coverpkg ./... -covermode=atomic -coverprofile=coverage.out ./...
	@echo 💡 Generating HTML coverage report and opening browser
	@go tool cover -html coverage.out

clean:
	@echo 💩Cleaning up
	@-rm -fr dist
	@-rm -fr rpm-prep
	@-rm -f skogul
	@-rm -f docs/skogul.rst
	@-rm -f skogul.1
	@-rm -f *.rpm
	@-rm -f coverage.out

help:
	@echo "Several targets(🎯) exist:"
	@echo 
	@echo " - skogul - build the binary (the default)"
	@echo " - all - build binary and documentation"
	@echo " - install - install binary and docs. Honors PREFIX, default prefix: ${PREFIX}"
	@echo " - rpm - build RPM"
	@echo " - clean - remove known build crap - use git clean -fdx for more thorough cleaning"
	@echo " - test / bench - run go test, with and without benchmarks "
	@echo "                  note that this uses "-short" to avoid mysql/postgres dependencies. "
	@echo " - fmtcheck - Runs gofmt -d -s, excluding generated code"
	@echo " - fmtfix - Runs gofmt -d -s -w, excluding generated code (e.g.: fix formating)"
	@echo " - covergui - Run tests, track test coverage and open coverage analysis in browser"

.PHONY: clean test bench help install rpm
