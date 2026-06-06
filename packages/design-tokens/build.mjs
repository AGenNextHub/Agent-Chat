// Sync the canonical design assets from docs/brand/ into the package's dist/.
// docs/brand/tokens.json is the single source of truth; this package is a
// generated artifact (no hand-edited copies — keeps the brand canonical).
import { mkdirSync, copyFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const root = resolve(here, "..", ".."); // repo root
const dist = resolve(here, "dist");

mkdirSync(dist, { recursive: true });
for (const f of ["tokens.json", "tokens.css", "logo.svg"]) {
  copyFileSync(resolve(root, "docs", "brand", f), resolve(dist, f));
}
copyFileSync(resolve(root, "LICENSE"), resolve(here, "LICENSE"));
console.log("design-tokens: synced canonical assets from docs/brand/ → dist/");
