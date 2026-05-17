package main

import (
	"fmt"
	"sort"
	"strings"
)

func rootReadme() string {
	return `# @chill-institute/cli

npm distribution for the ` + "`chilly`" + ` command-line client.

` + "```sh" + `
npm install -g @chill-institute/cli
chilly version
` + "```" + `
`
}

func platformReadme(t target) string {
	return "# " + platformPackageName(t) + "\n\n" +
		"Platform binary package for `@chill-institute/cli` on `" + t.npmOS + "-" + t.npmArch + "`.\n"
}

func wrapperScript() string {
	entries := make([]string, 0, len(targets))
	for _, t := range targets {
		entries = append(entries, fmt.Sprintf(
			"  %q: { packageName: %q, binaryName: %q }",
			t.npmOS+" "+t.npmArch,
			platformPackageName(t),
			t.binaryFile,
		))
	}
	sort.Strings(entries)

	return `#!/usr/bin/env node
"use strict";

const { spawn } = require("node:child_process");
const { dirname, join } = require("node:path");

const targets = {
` + strings.Join(entries, ",\n") + `
};

const targetKey = process.platform + " " + process.arch;
const target = targets[targetKey];
if (!target) {
  console.error("chilly is not available for " + process.platform + "-" + process.arch);
  process.exit(1);
}

let packageRoot;
try {
  packageRoot = dirname(require.resolve(target.packageName + "/package.json"));
} catch (error) {
  console.error("Missing optional dependency " + target.packageName + "; reinstall @chill-institute/cli with optional dependencies enabled.");
  process.exit(1);
}

const child = spawn(join(packageRoot, "bin", target.binaryName), process.argv.slice(2), {
  stdio: "inherit",
  windowsHide: false
});

child.on("error", (error) => {
  console.error("Failed to start chilly: " + error.message);
  process.exit(1);
});

child.on("exit", (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }
  process.exit(code ?? 1);
});
`
}
