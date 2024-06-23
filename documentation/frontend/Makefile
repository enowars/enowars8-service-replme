OUTDIR = out

HTMLIN = $(wildcard src/html/*.html)
HTMLOUT = $(patsubst src/html/%.html,$(OUTDIR)/%.html,$(HTMLIN))

JSIN = $(filter-out src/ts/lib.ts, $(wildcard src/ts/*.ts))
JSOUT = $(patsubst src/ts/%.ts,$(OUTDIR)/static/js/%.js,$(JSIN))

CSSIN = $(wildcard src/css/*.css)
CSSOUT = $(subst src/,$(OUTDIR)/static/,$(CSSIN))

build: static $(HTMLOUT) $(JSOUT) $(CSSOUT)

$(OUTDIR)/%.html: src/html/%.html
	@mkdir -p $(OUTDIR)/static
	node build-html.js $^

$(OUTDIR)/static/js/%.js: src/ts/%.ts
	@mkdir -p $(OUTDIR)/static/js
	@npm run webpack -- --config webpack.config.js --entry ./$^ --output-filename ./$@

$(OUTDIR)/static/css/%.css: src/css/%.css
	@mkdir -p $(OUTDIR)/static/css
	@npm run tailwind -- -i $^ -o $@

.PHONY: static
static:
	@mkdir -p $(OUTDIR)/static
	cp -r static/* $(OUTDIR)/static

.PHONY: clean
clean:
	@rm -rf $(OUTDIR)
