SRCDIR = src
SRCFILE = $(SRCDIR)/main.tex
OUTDIR = target
TARGET = $(OUTDIR)/main.pdf
EXPORT = ../iisc-test-ctf-slides.pdf

build: ${TARGET}

${TARGET}: ${SRCFILE}
	mkdir -p $(OUTDIR)
	lualatex --interaction=batchmode --output-directory=$(OUTDIR) $(SRCFILE)

export: ${TARGET}
	cp $(TARGET) $(EXPORT)

.PHONY: watch
watch:
	while inotifywait -e close_write -r ./src; do make build; done

.PHONE: clean
clean:
	rm -rf $(OUTDIR)
