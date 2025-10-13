prefix = /usr/local
exec_prefix = $(prefix)
bindir = $(exec_prefix)/bin
BASHCOMPLETIONSDIR = $(exec_prefix)/share/bash-completion/completions


RM = rm -f
INSTALL = install -D
MKDIRP = mkdir -p

.PHONY: install uninstall build clean default
default: build
build:
	go build
clean:
	go clean
reinstall: uninstall install
install:
	$(INSTALL) km-probe $(DESTDIR)$(bindir)/km-probe
	$(MKDIRP) $(DESTDIR)$(BASHCOMPLETIONSDIR)
	$(DESTDIR)$(bindir)/km-probe completion bash > $(DESTDIR)$(BASHCOMPLETIONSDIR)/km-probe
	@echo "================================="
	@echo ">> Now run the following command:"
	@echo "\tsource $(DESTDIR)$(BASHCOMPLETIONSDIR)/km-probe"
	@echo "================================="
uninstall:
	$(RM) $(DESTDIR)$(bindir)/km-probe
	$(RM) $(DESTDIR)$(BASHCOMPLETIONSDIR)/km-probe
