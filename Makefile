APP := secretd
AUR_PKGNAME = secretd-git
VERSION ?= $(shell git describe --tags --dirty --always 2>/dev/null || echo dev)

PREFIX ?= /usr
BINDIR := $(PREFIX)/bin
MANDIR := $(PREFIX)/share/man/man1
BASHDIR := $(PREFIX)/share/bash-completion/completions
ZSHDIR  := $(PREFIX)/share/zsh/site-functions
FISHDIR := $(PREFIX)/share/fish/vendor_completions.d
LICENSEDIR := $(PREFIX)/share/licenses
SYSCONFDIR := $(PREFIX)/share/secretd


AUR_DIR := packaging/aur
BUILD_DIR := build

.PHONY: all build install uninstall clean aur release check

all: build

build:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build \
		-trimpath \
		-ldflags "-s -w -X main.version=$(VERSION)" \
		-o $(BUILD_DIR)/$(APP) ./cmd/secretd

check:
	go vet ./...
	go test ./...
	go build ./cmd/secretd

install: build
	install -Dm755 $(BUILD_DIR)/$(APP) $(DESTDIR)$(BINDIR)/$(APP)
	install -Dm644 LICENSE $(DESTDIR)$(LICENSEDIR)/$(APP)/LICENSE
	install -Dm644 man/$(APP).1 $(DESTDIR)$(MANDIR)/$(APP).1

	install -Dm644 scripts/completions/$(APP).bash $(DESTDIR)$(BASHDIR)/$(APP)
	install -Dm644 scripts/completions/$(APP).zsh  $(DESTDIR)$(ZSHDIR)/_$(APP)
	install -Dm644 scripts/completions/$(APP).fish $(DESTDIR)$(FISHDIR)/$(APP).fish

	install -Dm644 internal/config/defaults.toml $(DESTDIR)$(SYSCONFDIR)/config.toml

uninstall:
	rm -f $(DESTDIR)$(BINDIR)/$(APP)
	rm -rf $(DESTDIR)$(LICENSEDIR)/$(APP)
	rm -f $(DESTDIR)$(MANDIR)/$(APP).1
	rm -f $(DESTDIR)$(BASHDIR)/$(APP)
	rm -f $(DESTDIR)$(ZSHDIR)/_$(APP)
	rm -f $(DESTDIR)$(FISHDIR)/$(APP).fish

clean:
	rm -rf $(BUILD_DIR)

aur: build
	@echo "Preparing AUR directory..."
	mkdir -p $(AUR_DIR)
	rsync -av --exclude='.git' --exclude='build' ./packaging/aur/ $(AUR_DIR)/
	@echo "Updating PKGBUILD version..."
	sed -i "s/^pkgver=.*/pkgver=$(VERSION)/" $(AUR_DIR)/PKGBUILD
	@echo "Generating .SRCINFO..."
	cd $(AUR_DIR) && makepkg --printsrcinfo > .SRCINFO
	@echo "Adding to git..."
	cd $(AUR_DIR) && git add PKGBUILD .SRCINFO
	cd $(AUR_DIR) && git commit -m "Update $(AUR_PKGNAME) version $(VERSION)" || echo "Nothing to commit"
	cd $(AUR_DIR) && git push origin master

release: clean
	git archive \
		--prefix=$(APP)-$(VERSION)/ \
		-o $(APP)-$(VERSION).tar.gz \
		HEAD

