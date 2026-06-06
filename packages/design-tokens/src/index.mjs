// AGenNext Chat design tokens — JS API.
// Reads the shipped DTCG tokens; resolves semantic aliases to hex.
import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const distDir = resolve(here, "..", "dist");

/** Canonical design tokens (W3C DTCG format). */
export const tokens = JSON.parse(readFileSync(resolve(distDir, "tokens.json"), "utf8"));

/** Absolute paths to the shipped assets. */
export const paths = {
  tokensJson: resolve(distDir, "tokens.json"),
  tokensCss: resolve(distDir, "tokens.css"),
  logoSvg: resolve(distDir, "logo.svg"),
};

/**
 * Resolve a semantic color token (e.g. "primary", "accent", "danger") to its
 * hex value, following the DTCG alias to the primitive palette.
 * @param {string} name
 * @returns {string|undefined}
 */
export function color(name) {
  const sem = tokens.semantic?.[name];
  if (!sem) return undefined;
  const ref = String(sem.$value).match(/^\{color\.(.+)\}$/);
  if (!ref) return sem.$value;
  let node = tokens.color;
  for (const seg of ref[1].split(".")) node = node?.[seg];
  return node?.$value ?? sem.$value;
}

export default { tokens, paths, color };
