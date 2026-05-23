---
name: pp-prediction-goat
description: "Every Polymarket and Kalshi market in one slim agent-native CLI, with cross-venue topic search and screens no other... Trigger phrases: `what are the odds on`, `polymarket odds for`, `kalshi odds for`, `find prediction markets for`, `what's trending on polymarket`, `what's resolving this week`, `compare polymarket and kalshi on`, `use prediction-goat`, `run prediction-goat`."
author: "Matt Van Horn"
license: "Apache-2.0"
argument-hint: "<command> [args] | install cli|mcp"
allowed-tools: "Read Bash"
metadata:
  openclaw:
    requires:
      bins:
        - prediction-goat-pp-cli
---

# Polymarket + Kalshi — Printing Press CLI

## Prerequisites: Install the CLI

This skill drives the `prediction-goat-pp-cli` binary. **You must verify the CLI is installed before invoking any command from this skill.** If it is missing, install it first:

1. Install via the Printing Press installer:
   ```bash
   npx -y @mvanhorn/printing-press install prediction-goat --cli-only
   ```
2. Verify: `prediction-goat-pp-cli --version`
3. Ensure `$GOPATH/bin` (or `$HOME/go/bin`) is on `$PATH`.

If the `npx` install fails (no Node, offline, etc.), fall back to a direct Go install (requires Go 1.26.3 or newer):

```bash
go install github.com/mvanhorn/printing-press-library/library/payments/prediction-goat/cmd/prediction-goat-pp-cli@latest
```

If `--version` reports "command not found" after install, the install step did not put the binary on `$PATH`. Do not proceed with skill commands until verification succeeds.

Read-only by design and by CI lint — the binary structurally cannot trade. `topic <name>` returns every related Polymarket + Kalshi market/event/tag in one ~3KB ranked bundle (vs the official Polymarket CLI's ~250KB firehose). Local SQLite + FTS5 keeps queries instant and free after one sync. Six screens (`trending`, `resolving`, `liquid`, `mispriced`, `movers`, `new`) cover the workflows agents and odds researchers run every week.

## When to Use This CLI

Reach for prediction-goat-pp-cli when an agent needs current prediction-market odds across both Polymarket and Kalshi without trading. The killer commands are `topic`, `compare`, and the six screens — every other command exists so power users can drill into one venue or one market. The CLI is read-only by structural CI guarantee: it cannot place orders, hold a wallet, or sign trades, which makes it safe to embed in agent toolchains.

## Unique Capabilities

These capabilities aren't available in any other tool for this API.

### Cross-venue intelligence
- **`topic`** — Get every related Polymarket and Kalshi market for a topic in one slim ranked ~3KB bundle — kanye-west, argentina, chatgpt-5.

  _When an agent needs current odds on a topic, this is the one-call answer across both venues without fanning out to two platform tools and re-ranking by hand._

  ```bash
  prediction-goat-pp-cli topic kanye-west --json
  ```
- **`mispriced`** — Find same-outcome markets where Polymarket and Kalshi disagree on implied probability by more than a threshold.

  _The clearest signal that one venue is wrong or one side is mispricing — useful for calibration research, not trading._

  ```bash
  prediction-goat-pp-cli mispriced --threshold 0.05 --json
  ```
- **`compare`** — Side-by-side YES/NO and implied probability for the same topic across Polymarket and Kalshi.

  _Tells an agent or analyst 'which venue has the better/different number on this outcome' in one read-only call._

  ```bash
  prediction-goat-pp-cli compare 'arizona basketball' --json
  ```
- **`markets diff`** — Field-by-field structural diff between a specific Polymarket market and a specific Kalshi market.

  _When you already know the two slugs/tickers (e.g. from `topic <theme>`), diff shows you exactly where the venues disagree._

  ```bash
  prediction-goat-pp-cli markets diff <pm-slug> <kalshi-ticker> --json
  ```

### Screens
- **`trending`** — Top movers by 24h volume across both venues, ranked.

  _One command answers 'what should I be watching today' without scraping two homepages._

  ```bash
  prediction-goat-pp-cli trending --json --limit 20
  ```
- **`resolving`** — Markets resolving in the next week/month/days, sorted by liquidity.

  _Tells an agent 'what's about to settle' without re-paging two cursors._

  ```bash
  prediction-goat-pp-cli resolving --week --json
  ```
- **`liquid`** — Markets above a normalized volume/liquidity floor across both venues.

  _Filters out thin markets that will move on a single 100-dollar bet._

  ```bash
  prediction-goat-pp-cli liquid --min-volume 100000 --json
  ```
- **`movers`** — Biggest implied-probability deltas over a 24h or 7d window across both venues.

  _Surfaces narrative shifts (price-driven) vs hype shifts (volume-driven from )._

  ```bash
  prediction-goat-pp-cli movers --window 7d --json
  ```
- **`new`** — Markets created in the last N days across both venues.

  _Newly listed markets are where the alpha and mispricings live._

  ```bash
  prediction-goat-pp-cli new --days 7 --json
  ```

## Command Reference

**comments** — Comment system and user interactions

- `prediction-goat-pp-cli comments get-by-id` — Get comments by comment id
- `prediction-goat-pp-cli comments get-by-user-address` — Get comments by user address
- `prediction-goat-pp-cli comments list` — List comments

**events** — Event management and event-related operations

- `prediction-goat-pp-cli events get` — Get event by id
- `prediction-goat-pp-cli events get-by-slug` — Get event by slug
- `prediction-goat-pp-cli events get-creator` — Get event creator by id
- `prediction-goat-pp-cli events list` — List events
- `prediction-goat-pp-cli events list-creators` — List event creators
- `prediction-goat-pp-cli events list-keyset` — Returns events using cursor-based (keyset) pagination for stable, efficient paging through large result sets. Use...
- `prediction-goat-pp-cli events list-pagination` — List events (paginated)
- `prediction-goat-pp-cli events list-sport-results` — List sport events results

**markets** — Market data and market-related operations

- `prediction-goat-pp-cli markets get` — Get market by id
- `prediction-goat-pp-cli markets get-abridged` — Query abridged markets by information filters
- `prediction-goat-pp-cli markets get-by-slug` — Get market by slug
- `prediction-goat-pp-cli markets get-information` — Query markets by information filters
- `prediction-goat-pp-cli markets list` — List markets
- `prediction-goat-pp-cli markets list-keyset` — Returns markets using cursor-based (keyset) pagination for stable, efficient paging through large result sets. Use...

**profiles** — User profile management

- `prediction-goat-pp-cli profiles <user_address>` — Get public profile by user address

**public-profile** — Manage public profile

- `prediction-goat-pp-cli public-profile` — Get public profile by wallet address

**public-search** — Manage public search

- `prediction-goat-pp-cli public-search` — Search markets, events, and profiles

**series** — Series management and related operations

- `prediction-goat-pp-cli series get` — Get series by id
- `prediction-goat-pp-cli series list` — List series

**series-summary** — Manage series summary

- `prediction-goat-pp-cli series-summary get-by-id` — Get series summary by id
- `prediction-goat-pp-cli series-summary get-by-slug` — Get series summary by slug

**sports** — Sports-related endpoints including teams and game data

- `prediction-goat-pp-cli sports get-market-types` — Get valid sports market types
- `prediction-goat-pp-cli sports get-metadata` — Get sports metadata information

**status** — Manage status

- `prediction-goat-pp-cli status` — Gamma API Health check

**tags** — Tag management and related tag operations

- `prediction-goat-pp-cli tags get` — Get tag by id
- `prediction-goat-pp-cli tags get-by-slug` — Get tag by slug
- `prediction-goat-pp-cli tags get-related-by-slug` — Get related tags (relationships) by tag slug
- `prediction-goat-pp-cli tags get-related-to-atag-by-slug` — Get tags related to a tag slug
- `prediction-goat-pp-cli tags list` — List tags

**teams** — Manage teams

- `prediction-goat-pp-cli teams get` — Get team by id
- `prediction-goat-pp-cli teams list` — List teams


### Finding the right command

When you know what you want to do but not which command does it, ask the CLI directly:

```bash
prediction-goat-pp-cli which "<capability in your own words>"
```

`which` resolves a natural-language capability query to the best matching command from this CLI's curated feature index. Exit code `0` means at least one match; exit code `2` means no confident match — fall back to `--help` or use a narrower query.

## Recipes


### Find every market for a topic across both venues

```bash
prediction-goat-pp-cli topic kanye-west --json --select markets.title,markets.venue,markets.yesProbability,markets.endDate
```

Slim ranked bundle of every related PM + Kalshi market. `--select` reduces to four fields per row; an agent sees ~1KB instead of the firehose.

### What's settling this week?

```bash
prediction-goat-pp-cli resolving --week --json --select title,venue,endDate,liquidity --limit 10
```

Local SQL filter on end_date < now+7d across both venues, sorted by liquidity descending. `--select` keeps the response tiny.

### Side-by-side odds comparison

```bash
prediction-goat-pp-cli compare 'arizona basketball' --json
```

Resolves the topic to paired markets and renders YES/NO + implied prob for each venue side-by-side. Tells you exactly where PM and Kalshi disagree.

### Enumerate every child market under a Kalshi event

When the event ticker is known (often discovered via `kalshi-series-search` then `kalshi events list --series`), one call returns every child market with live prices:

```bash
prediction-goat-pp-cli kalshi events get KXMENWORLDCUP-26 --with-markets --agent
```

Passes `with_nested_markets=true` to the upstream `/events/{ticker}` endpoint. The response includes a `markets` array with ticker, title, yes_sub_title, status, yes_ask_dollars, no_ask_dollars, volume_24h_fp, and expiration_time for each child. This is the lightweight alternative to a full sync walk when you only need one event's children.

### List every Kalshi event under a series

```bash
prediction-goat-pp-cli kalshi events list --series KXMENWORLDCUP --agent
```

Filters `/events` by `series_ticker`. The `--series` flag forces `--data-source live` since the local store doesn't index by series ticker. Use this before `kalshi events get --with-markets` when you don't know the exact event ticker — series → event → markets in three calls.

### Catch mispricings across venues

```bash
prediction-goat-pp-cli mispriced --threshold 0.05 --json --select pair.pm.title,pair.kalshi.title,delta
```

Returns same-outcome market pairs where implied probabilities diverge by 5+ percentage points. Slim output via `--select`.

## Auth Setup

No authentication required.

Run `prediction-goat-pp-cli doctor` to verify setup.

## Agent Mode

Add `--agent` to any command. Expands to: `--json --compact --no-input --no-color --yes`.

- **Pipeable** — JSON on stdout, errors on stderr
- **Filterable** — `--select` keeps a subset of fields. Dotted paths descend into nested structures; arrays traverse element-wise. Critical for keeping context small on verbose APIs:

  ```bash
  prediction-goat-pp-cli comments list --agent --select id,name,status
  ```
- **Previewable** — `--dry-run` shows the request without sending
- **Offline-friendly** — sync/search commands can use the local SQLite store when available
- **Non-interactive** — never prompts, every input is a flag
- **Explicit retries** — use `--idempotent` only when an already-existing create should count as success

### Response envelope

Commands that read from the local store or the API wrap output in a provenance envelope:

```json
{
  "meta": {"source": "live" | "local", "synced_at": "...", "reason": "..."},
  "results": <data>
}
```

Parse `.results` for data and `.meta.source` to know whether it's live or local. A human-readable `N results (live)` summary is printed to stderr only when stdout is a terminal AND no machine-format flag (`--json`, `--csv`, `--compact`, `--quiet`, `--plain`, `--select`) is set — piped/agent consumers and explicit-format runs get pure JSON on stdout.

## Agent Feedback

When you (or the agent) notice something off about this CLI, record it:

```
prediction-goat-pp-cli feedback "the --since flag is inclusive but docs say exclusive"
prediction-goat-pp-cli feedback --stdin < notes.txt
prediction-goat-pp-cli feedback list --json --limit 10
```

Entries are stored locally at `~/.prediction-goat-pp-cli/feedback.jsonl`. They are never POSTed unless `PREDICTION_GOAT_FEEDBACK_ENDPOINT` is set AND either `--send` is passed or `PREDICTION_GOAT_FEEDBACK_AUTO_SEND=true`. Default behavior is local-only.

Write what *surprised* you, not a bug report. Short, specific, one line: that is the part that compounds.

## Output Delivery

Every command accepts `--deliver <sink>`. The output goes to the named sink in addition to (or instead of) stdout, so agents can route command results without hand-piping. Three sinks are supported:

| Sink | Effect |
|------|--------|
| `stdout` | Default; write to stdout only |
| `file:<path>` | Atomically write output to `<path>` (tmp + rename) |
| `webhook:<url>` | POST the output body to the URL (`application/json` or `application/x-ndjson` when `--compact`) |

Unknown schemes are refused with a structured error naming the supported set. Webhook failures return non-zero and log the URL + HTTP status on stderr.

## Named Profiles

A profile is a saved set of flag values, reused across invocations. Use it when a scheduled agent calls the same command every run with the same configuration - HeyGen's "Beacon" pattern.

```
prediction-goat-pp-cli profile save briefing --json
prediction-goat-pp-cli --profile briefing comments list
prediction-goat-pp-cli profile list --json
prediction-goat-pp-cli profile show briefing
prediction-goat-pp-cli profile delete briefing --yes
```

Explicit flags always win over profile values; profile values win over defaults. `agent-context` lists all available profiles under `available_profiles` so introspecting agents discover them at runtime.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 2 | Usage error (wrong arguments) |
| 3 | Resource not found |
| 5 | API error (upstream issue) |
| 7 | Rate limited (wait and retry) |
| 10 | Config error |

## Argument Parsing

Parse `$ARGUMENTS`:

1. **Empty, `help`, or `--help`** → show `prediction-goat-pp-cli --help` output
2. **Starts with `install`** → ends with `mcp` → MCP installation; otherwise → see Prerequisites above
3. **Anything else** → Direct Use (execute as CLI command with `--agent`)

## MCP Server Installation

Install the MCP binary from this CLI's published public-library entry or pre-built release, then register it:

```bash
claude mcp add prediction-goat-pp-mcp -- prediction-goat-pp-mcp
```

Verify: `claude mcp list`

## Direct Use

1. Check if installed: `which prediction-goat-pp-cli`
   If not found, offer to install (see Prerequisites at the top of this skill).
2. Match the user query to the best command from the Unique Capabilities and Command Reference above.
3. Execute with the `--agent` flag:
   ```bash
   prediction-goat-pp-cli <command> [subcommand] [args] --agent
   ```
4. If ambiguous, drill into subcommand help: `prediction-goat-pp-cli <command> --help`.
