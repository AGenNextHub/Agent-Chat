# Channels

The core is **headless** and has one contract — **CloudEvents**. A **channel** is a
composable **adapter** that connects a touch point to the core: it maps an inbound message
to a CloudEvent and an outbound result back to its medium.

**One interface, every platform, any tool.**
- **One interface** — `pkg/channel.Adapter`. The same shape is the touchpoint for every
  platform; nothing in the core is special-cased per channel.
- **Every platform** — web, Slack, Mattermost, Matrix, WhatsApp … add one by implementing
  `Adapter`; the core is untouched.
- **Any tool** — a tool is a `Capability` (a contract). It is exposed through the same
  interface by **composition**, not by bespoke wiring.
- **Multimodal** — payloads are content-type agnostic via the event envelope
  (`Data` + `DataContentType`).

| Channel | Status | Adapter |
|---|---|---|
| web (`channels/web`) | static client | calls `POST /v1/chat` directly |
| slack | planned | Bolt / Events API |
| mattermost | planned | bot API |
| matrix | planned | client SDK / appservice |
| whatsapp | planned | WhatsApp Business API |

## web channel

`channels/web/index.html` — a single static page, no build step, styled with the design
tokens (`packages/design-tokens`), calling `POST /v1/chat` and rendering the inspectable
trace. Serve it from any static host, or from the daemon (same origin), pointed at the API.

It is also the surface for the **agent builder**: constructing an agent = defining a
`Capability` contract (name · scope · provides) — the same contract the kernel admits. Build
an agent, expose a tool, talk to it — through one interface, on any platform.
