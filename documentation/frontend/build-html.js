#!/usr/bin/env node

const fs = require("node:fs");
const HTMLParser = require("node-html-parser");
const { exit } = require("node:process");
const prettier = require("prettier");
const Path = require("path")

const outDir = "out";

fs.mkdirSync(outDir, {
  recursive: true,
});

if (process.argv.length < 3) {
  console.error("ERROR: Missing source files");
  console.error(`Usage: ./${Path.basename(process.argv[1])} ... FILES`);
  exit(1);
}

const sources = process.argv.slice(2)

async function build(sources, outDir) {
  for (let i = 0; i < sources.length; i++) {
    const source = sources[i]
    if (!fs.existsSync(source)) {
      console.error("ERROR: Source does not exist:", source);
      exit(1)
    }
    if (fs.lstatSync(source).isDirectory()) {
      console.error("ERROR: Source is a directory:", source);
      exit(1)
    }

    const workingDir = Path.dirname(source);
    const html = fs.readFileSync(source);
    const tree = HTMLParser.parse(html);
    let elem = tree.querySelector("layout");

    while (elem) {
      let layoutSrc = elem.getAttribute("src");
      if (!layoutSrc) {
        console.error("ERROR: Layout tag missing src attribute", layoutSrc);
        exit(1);
      }
      layoutSrc = Path.join(workingDir, layoutSrc)
      if (!fs.existsSync(layoutSrc)) {
        console.error("ERROR: Source does not exist:", layoutSrc);
        exit(1)
      }
      if (fs.lstatSync(layoutSrc).isDirectory()) {
        console.error("ERROR: Source is a directory:", layoutSrc);
        exit(1)
      }

      let layout = fs.readFileSync(layoutSrc);
      const layoutTree = HTMLParser.parse(layout);
      const layoutChildren = layoutTree.querySelector("children");
      if (!layoutChildren) {
        console.error("ERROR: Layout missing children tag", layoutSrc);
        exit(1);
      }
      layoutChildren.replaceWith(...elem.childNodes);

      elem.replaceWith(layoutTree);
      elem = tree.querySelector("layout");
    }

    await prettier
      .format(tree.toString(), {
        parser: "html",
      })
      .then((result) => {
        fs.writeFileSync(Path.join(outDir, Path.basename(source)), result);
      });
  }
}

build(sources, outDir)
