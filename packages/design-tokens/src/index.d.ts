export interface Tokens {
  $schema?: string;
  [group: string]: unknown;
}

/** Canonical design tokens (W3C DTCG format). */
export const tokens: Tokens;

/** Absolute paths to the shipped assets. */
export const paths: {
  tokensJson: string;
  tokensCss: string;
  logoSvg: string;
};

/** Resolve a semantic color token (e.g. "primary") to its hex value. */
export function color(name: string): string | undefined;

declare const _default: { tokens: Tokens; paths: typeof paths; color: typeof color };
export default _default;
