# @agennext/design-tokens

The canonical **AGenNext Chat** design tokens (W3C DTCG) + logo, as a lightweight,
OCI-compatible package. It is a **tool, not an app**: it carries no Penpot/Figma
application and **no database** — design tools consume the tokens via their SDK/plugin.

> Single source of truth: `docs/brand/tokens.json` in the monorepo. This package's `dist/`
> is generated from it (`npm run build`); never hand-edit derived files.

## Install (npm)

```bash
npm install @agennext/design-tokens
```

```js
import { tokens, color, paths } from "@agennext/design-tokens";
import "@agennext/design-tokens/tokens.css"; // CSS custom properties (--agx-*)

color("primary"); // "#4338CA"
```

## CLI

```bash
npx agennext-tokens dtcg   # W3C DTCG JSON — import into Penpot/Figma token plugins
npx agennext-tokens css    # CSS custom properties
npx agennext-tokens list   # semantic tokens → hex
```

## Use with Penpot / Figma (SDK as tool, not the app)

The `dtcg` output is the standard W3C Design Tokens format. Import it through the design
tool's **tokens plugin/SDK** — no Penpot app or database is bundled or required here:

- **Penpot:** Tokens → Import → select the `dtcg` JSON.
- **Figma:** Tokens Studio (or similar) → Import → the `dtcg` JSON.

## Distribute as an OCI artifact

The same files publish as a single, lightweight OCI artifact (no registry lock-in):

```bash
oci/push.sh ghcr.io/agennext/design-tokens:0.1.0
# pull anywhere:
oras pull ghcr.io/agennext/design-tokens:0.1.0
```

## Publishing (run with YOUR credentials)

This package is not auto-published. To release:

```bash
npm run build
npm publish --access public      # requires your npm auth
# and/or
oci/push.sh <your-registry>/design-tokens:0.1.0   # requires your registry auth
```

## License

Apache-2.0 — see [LICENSE](LICENSE).
