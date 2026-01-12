# PREFIX is the prefix on the targetsystem
# DESTDIR can be used to prefix ALL paths, e.g., to do a dummy-install in a
# fake root dir, e.g., for building packages. Users mainly want PREFIX

PREFIX = /usr/local
DOCDIR = ${PREFIX}/share/doc/skogul

GIT_DESCRIBE := $(shell git describe --always --tag --dirty)
VERSION_NO = $(shell echo ${GIT_DESCRIBE} | sed s/[v-]//g)
OS := $(shell uname -s | tr A-Z a-z)
ARCH := $(shell uname -m)

skogul: $(wildcard *.go */*.go */*/*.go)
	@echo 🤸 go build !
	@CGO_ENABLED=0 go build -ldflags "-X main.versionNo=${VERSION_NO}" -o skogul ./cmd/skogul

generate:
	@echo 🔧 Generating protocol buffer code
	@./gen/generate.sh

docs/skogul.rst: skogul
	@echo 😽 Generating documentation$@
	@./skogul -make-man > $@

skogul.1: docs/skogul.rst
	@echo 🎢 Generating man-file$@
	@rst2man < $< > $@

notes: docs/NEWS
	@echo ⛲ Extracting release notes.
	@./build/release-notes.sh $$(echo ${GIT_DESCRIBE} | sed s/-dirty//) > notes

all: skogul skogul.1 docs/skogul.rst

install: skogul skogul.1 docs/skogul.rst
	@echo 🙅 Installing
	@install -D -m 0755 skogul ${DESTDIR}${PREFIX}/bin/skogul
	@install -D -m 0644 skogul.1 ${DESTDIR}${PREFIX}/share/man/man1/skogul.1
	@install -D -m 0644 docs/examples/basics/default.json ${DESTDIR}/etc/skogul/conf.d/default.json
	@cd docs && find . -type f -exec install -D -m 0644 {} ${DESTDIR}${DOCDIR}/{} \;
	@install -D -m 0644 README.rst LICENSE -t ${DESTDIR}${DOCDIR}/

# Any complaints on this macro-substitution without patches and I introduce m4.
build/redhat-skogul.spec: build/redhat-skogul.spec.in FORCE
	@echo ❕Building spec-file
	@cat $< | sed "s/xxVxx/${GIT_DESCRIBE}/g; s/xxARCHxx/${ARCH}/g; s/xxVERSION_NOxx/${VERSION_NO}/g" > $@
	@if [ ! -f /etc/redhat-release ]; then \
		echo 🆒 Adding debian-workaround for rpm build; \
		sed -i 's/^BuildReq/\#Debian hack, auto-commented out: BuildReq/g'$@; \
	fi

# Build RPM. The spec has a blank %prep, so it assumes sources are already
# available. This isn't perfect, since it creates a tight coupling between
# Makefile and specfile, but it isn't all that bad either, since it allows
# building with minimal redundant effort, and without having to commit to
# git.
rpm: build/redhat-skogul.spec
	@echo 🎇 Triggering huge-as-heck rpm build
	@mkdir -p rpm-prep/BUILDROOOT
	@DEFAULT_UNIT_DIR=/usr/lib/systemd/system ; \
	RPM_UNIT_DIR=$$(rpm --eval $%{_unitdir}) ; \
	if [ "$${RPM_UNIT_DIR}" = "$%{_unitdir}" ]; then \
		echo "😭 _unitdir not set, setting _unitdir to $$DEFAULT_UNIT_DIR"; \
		rpmbuild --quiet --bb \
		--nodebuginfo \
		--build-in-place \
		--define "_rpmdir $$(pwd)" \
		--define "_topdir $$(pwd)" \
		--define "_unitdir $$DEFAULT_UNIT_DIR" \
		--buildroot "$$(pwd)/rpm-prep/BUILDROOT" \
		build/redhat-skogul.spec; \
	else \
		rpmbuild --quiet --bb \
		--nodebuginfo \
		--build-in-place \
		--define "_rpmdir $$(pwd)" \
		--define "_topdir $$(pwd)" \
		--buildroot "$$(pwd)/rpm-prep/BUILDROOT" \
		build/redhat-skogul.spec; \
	fi
	@cp x86_64/skogul-${VERSION_NO}-1.x86_64.rpm .
	@echo ⭐ RPM built: ./skogul-${VERSION_NO}-1.x86_64.rpm

check: test fmtcheck vet exampletest exampletestdep checkconfigs printfcheck

vet:
	@echo 🔬 Vetting code
	@go vet ./...

fmtcheck:
	@echo 🦉 Checking format with gofmt -d -s
	@if [ "x$$(find . -name '*.go' -not -wholename './gen/*' -and -not -wholename './vendor/*' -exec gofmt -d -s {} +)" != "x" ]; then \
	find . -name '*.go' -not -wholename './gen/*' -and -not -wholename './vendor/*' -exec gofmt -d -s {} +; \
		exit 1; \
	fi

fmtfix:
	@echo 🎨 Fixing formating
	@find . -name '*.go' -not -wholename './gen/*' -and -not -wholename './vendor/*' -exec gofmt -d -s -w {} +

printfcheck:
	@echo 📖 Looking for printf-debugging left over
	@! find -not -wholename './cmd/*' -and -not -wholename '*_test.go' -and -not -wholename './config/parse.go' -and -not -wholename './sender/debug.go' -and -name '*.go' -exec egrep fmt.Printf {} +

exampletest: skogul
	@echo 📖 Verifying examples
	@failed=0; for a in $$(find docs/examples/ -name '*json'  | grep -v payloads | grep -v client-certificates | grep -v sql-tls | grep -v juniper); do \
		./skogul -show -f $$a >/dev/null 2>&1 ; \
		if [ $$? -ne 0 ]; then \
			echo 🚩 Example $$a is not valid; \
			failed=$$(( failed + 1 )); \
		fi; \
	done; \
	exit $${failed}
	@echo 📖 Verifying junos example
	@./skogul -show -d docs/examples/juniper >/dev/null 2>&1; \
	if [ $$? -ne 0 ]; then \
		echo 🚩 Junos-example is not valid; \
		exit 1; \
	fi

checkbadconfigs: skogul
	@echo 📖 Verifying that invalid configuration files are caught
	@failed=0; for a in $$(find testdata/invalid_configs/ -name '*json'); do \
		./skogul -show -f $$a >/dev/null 2>&1 ; \
		if [ $$? -eq 0 ]; then \
			echo 🚩 Invalid config $$a was accepted, but should fail; \
			failed=$$(( failed + 1 )); \
		fi; \
	done; \
	exit $${failed}

checkokconfigs: skogul
	@echo 📖 Verifying that valid configuration files are accepted
	@failed=0; for a in $$(find testdata/valid_configs/ -name '*json'); do \
		./skogul -show -f $$a >/dev/null 2>&1 ; \
		if [ $$? -ne 0 ]; then \
			echo 🚩 Valid config $$a was rejected; \
			failed=$$(( failed + 1 )); \
		fi; \
	done; \
	exit $${failed}

checkconfigs: checkbadconfigs checkokconfigs

exampletestdep: exampletest
	@echo 📖 Checking examples for deprecation warnings
	@failed=0; for a in $$(find docs/examples/ -name '*json'  | grep -v payloads | grep -v client-certificates | grep -v sql-tls | grep -v juniper); do \
		./skogul -show -f $$a 2>&1 | egrep -q "deprecation warning for" ; \
		if [ $$? -eq 0 ]; then \
			echo 🚩 Example $$a has deprecation warnings; \
			failed=$$(( failed + 1 )); \
		fi; \
	done; \
	exit $${failed}
	@echo 📖 Checking junos example for deprecation warnings
	@./skogul -show -d docs/examples/juniper 2>&1 | egrep -q 'deprecation warning for'; \
	if [ $$? -eq 0 ]; then \
		echo 🚩 Junos-example has deprecation warnings; \
		exit 1; \
	fi

test:
	@echo 🧐 Testing, without SQL-tests
	@go test -short ./...

test-sql-tls: skogul
	@echo 🔐 Running SQL TLS integration tests
	@cd testdata/sql-tls && ./run-tests.sh

bench:
	@echo 🏋 Benchmarking
	@go test -run ^Bench -benchtime 1s -bench Bench ./... | grep --line-buffered Benchmark | awk -v term_width=$$(tput cols 2>/dev/null || echo 120) 'BEGIN { \
		num_width = 13; \
		spacing = 2; \
		fixed_width = (num_width * 4) + (spacing * 4); \
		name_width = term_width - fixed_width; \
		if (name_width < 30) name_width = 30; \
		if (name_width > 60) name_width = 60; \
		fmt_name = "%-" name_width "s"; \
		fmt_num = "%" num_width "s"; \
		printf fmt_name " " fmt_num " " fmt_num " " fmt_num " " fmt_num "\n", "Test", "Iterations", "ns/op", "B/op", "allocs/op"; \
		printf fmt_name " " fmt_num " " fmt_num " " fmt_num " " fmt_num "\n", "----", "----------", "-----", "----", "---------"; \
		fflush(); \
	} { \
		name = $$1; \
		sub(/^Benchmark/, "", name); \
		sub(/ProtoBuf/, "PB", name); \
		sub(/InfluxDB/, "Influx", name); \
		sub(/MemoryFootprint/, "Mem", name); \
		sub(/WithoutTimestamp/, "NoTS", name); \
		sub(/SmallMessage/, "Small", name); \
		sub(/LargeMessage/, "Large", name); \
		iters = $$2; \
		ns = $$3; \
		bytes = "-"; \
		allocs = "-"; \
		for (i = 4; i <= NF; i++) { \
			if ($$i == "B/op" && i > 4) { \
				bytes = $$(i-1); \
			} \
			if ($$i == "allocs/op" && i > 4) { \
				allocs = $$(i-1); \
			} \
		} \
		printf fmt_name " " fmt_num " " fmt_num " " fmt_num " " fmt_num "\n", name, iters, ns, bytes, allocs; \
		fflush(); \
	}'

covergui:
	@echo 🧠 Testing, with coverage analysis
	@go test -short -coverpkg ./... -covermode=atomic -coverprofile=coverage.out ./...
	@echo 💡 Generating HTML coverage report and opening browser
	@go tool cover -html coverage.out

release:
	@if ! git diff-index --quiet HEAD; then echo "git working directory is not clean, it would be unfair to not release everything!" && exit 1; fi
	@git tag -a $$(head -n 1 docs/NEWS) -m "Skogul $$(head -n 1 docs/NEWS)"

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
	@echo " - check - run all tests we know and care about - run automatically on every push/tag/pr"
	@echo ""
	@echo " - rpm - build RPM"
	@echo " - clean - remove known build crap - use git clean -fdx for more thorough cleaning"
	@echo " - generate - regenerate protocol buffer code"
	@echo " - test / bench - run go test, with and without benchmarks "
	@echo "                  note that this uses "-short" to avoid mysql/postgres dependencies. "
	@echo " - test-sql-tls - Run SQL TLS integration tests using Docker (requires docker-compose)"
	@echo " - fmtcheck - Runs gofmt -d -s, excluding generated code"
	@echo " - fmtfix - Runs gofmt -d -s -w, excluding generated code (e.g.: fix formating)"
	@echo " - covergui - Run tests, track test coverage and open coverage analysis in browser"

.PHONY: all clean check checkconfigs test test-sql-tls bench help install rpm release FORCE
