
TOPDOCS=README.rst LICENSE
DIST=dist/
DOCDIR=${DIST}share/doc/skogul
MANDIR=${DIST}share/man/man1
MANFILE=${MANDIR}/skogul.1
DOCFILE=${DOCDIR}/skogul.rst
GIT_DESCRIBE:=$(shell git describe --always --tag --dirty)
OS:=$(shell uname -s | tr A-Z a-z)
ARCH:=$(shell uname -m)
DOCS=${TOPDOCS} ${MANFILE} ${DOCFILE} 
SKOGUL=${DIST}bin/skogul
TARBALL=skogul-${GIT_DESCRIBE}.${OS}-${ARCH}.tar.bz2

${SKOGUL}: $(wildcard *.go */*.go */*/*.go) | $(dir ${SKOGUL})
	go build -ldflags "-X main.versionNo=$V" -o ${SKOGUL} ./cmd/skogul

$(dir ${SKOGUL} ${MANFILE} ${DOCFILE}):
	mkdir -p $@

$(addprefix ${DOCDIR}/,${TOPDOCS}): ${TOPDOCS}
	cp $$(basename $@ ) $@

${MANFILE}: ${DOCFILE} | $(dir ${MANFILE})
	rst2man < $< > $@

${DOCFILE}: ${SKOGUL} | $(dir ${DOCFILE})
	${SKOGUL} -make-man > $@

${DOCDIR}/% : docs/%
	cp -a $< $@

notes:
	./build/release-notes.sh > notes

${TARBALL}: dist
	tar -C dist/ -cjf ${TARBALL} .

tar: ${TARBALL}

dist: ${SKOGUL} ${DOCS} $(addprefix ${DOCDIR}/,$(notdir $(wildcard docs/*)))
	touch dist

test:
	go test -short ./...

bench:
	go test -run ^Bench -benchtime 1s -bench Bench ./... | grep Benchmark

clean:
	rm -r dist

help:
	@echo "Several targets exist:"
	@echo 
	@echo " - dist/bin/skogul - build the binary "
	@echo " - dist - build the entire distribution in dist/ "
	@echo " - tar - build a tar ball (currently ${TARBALL}) "
	@echo " - test / bench - run go test, with and without benchmarks "
	@echo "                  note that this uses "-short" to avoid mysql/postgres dependencies. "

.PHONY: clean test bench tar help

