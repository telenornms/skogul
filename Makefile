PREFIX=/usr/local
TOPDOCS=README.rst LICENSE

GIT_DESCRIBE:=$(shell git describe --always --tag --dirty)
VERSION_NO=$(shell echo ${GIT_DESCRIBE} | sed s/[v-]//g)
OS:=$(shell uname -s | tr A-Z a-z)
ARCH:=$(shell uname -m)
TARBALL=skogul-${GIT_DESCRIBE}.${OS}-${ARCH}.tar.bz2

skogul: $(wildcard *.go */*.go */*/*.go)
	go build -ldflags "-X main.versionNo=$V" -o skogul ./cmd/skogul

docs/skogul.rst: skogul
	./skogul -make-man > $@

skogul.1: docs/skogul.rst
	rst2man < $< > $@

notes:
	./build/release-notes.sh > notes

all: skogul skogul.1 docs/skogul.rst

install: skogul skogul.1 docs/skogul.rst
	install -D -m 0755 skogul ${DESTDIR}${PREFIX}/bin/skogul
	install -D -m 0644 skogul.1 ${DESTDIR}${PREFIX}/share/man/man1/skogul.1
	install -D -m 0644 docs/examples/default.json ${DESTDIR}/etc/skogul/default.json
	cd docs; \
	find -type f -exec install -D -m 0644 {} ${DESTDIR}${PREFIX}/share/doc/skogul/{} \;
	install -D -m 0644 ${TOPDOCS} -t ${DESTDIR}${PREFIX}/share/doc/skogul/

%/:
	mkdir -p $@

rpm-prep/SPECS/skogul.spec: build/redhat-skogul.spec.in | rpm-prep/SPECS/
	cat $< | sed "s/xxVxx/${GIT_DESCRIBE}/g; s/xxARCHxx/${ARCH}/g; s/xxVERSION_NOxx/${VERSION_NO}/g" > $@

rpm: rpm-prep/SPECS/skogul.spec | rpm-prep/BUILDROOT/ rpm-prep/RPMS/ rpm-prep/SPECS/ rpm-prep/SRPMS/
	test -h rpm-prep/BUILD || ln -s ./ rpm-prep/BUILD
	build/trigger-rpm.sh
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
	@echo " - skogul - build the binary "
	@echo " - install - install binary and docs. Honors PREFIX, default prefix: ${PREFIX}"
	@echo " - rpm - build RPM"
	@echo " - test / bench - run go test, with and without benchmarks "
	@echo "                  note that this uses "-short" to avoid mysql/postgres dependencies. "

.PHONY: clean test bench help install

