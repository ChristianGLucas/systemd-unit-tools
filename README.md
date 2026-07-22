# systemd-unit-tools

Composable Axiom nodes for deterministic parsing and structural inspection of
systemd unit files, published as `christiangeorgelucas/systemd-unit-tools`.

Eighteen stateless nodes cover the systemd unit-file surface broadly: parsing
`.service`, `.timer`, `.socket`, `.mount`, `.automount`, `.swap`, `.path`,
`.target`, `.slice` and `.scope` files into a structured, repeated-directive-
preserving envelope, then extracting every semantically meaningful piece of
it — dependencies/ordering, the execution lifecycle, restart policy,
environment, install targets, timer schedules, socket listeners, process
identity, resource-limit/sandboxing hardening, metadata, and structural
validation.

## Why this is distinct from a generic INI/config parser

Systemd unit files are INI-*like* — `[Section]` headers, `Key=Value` lines,
`#`/`;` comments, `\`-continued lines — but they are not INI. Directives
legitimately **repeat** (`ExecStartPre=`, `After=`, `OnCalendar=`, ...), and a
plain INI parser that treats duplicate keys as "last one wins" silently
drops data a systemd-aware caller needs. This package's canonical `UnitFile`
envelope preserves every repeated directive, in source order, as a list —
never collapsed. See `christiangeorgelucas/config-tools` for the separate,
genuinely different domain of `.env`/INI/`.properties` files, which have no
such repeated-directive semantics.

## The `UnitFile` envelope

`ParseUnitFile` turns raw unit-file text into:

```
UnitFile { sections: [ Section { name, directives: [ Directive { key, value } ] } ] }
```

Every other node accepts either raw `text` (parsed internally — the common
case calling a node standalone) or a pre-parsed `unit: UnitFile` from
`ParseUnitFile`, to avoid re-parsing when a caller (an agent, a script)
already holds a parsed structure and wants to run several extraction nodes
against it. Note: as of Axiom's current flow compiler, a nested/repeated
MESSAGE-kind field (like `unit`) cannot cross a plain `flow.yaml` edge
adapter yet — pass `text` on each node when composing this package inside a
visual/flow.yaml graph; the `unit` shortcut is a direct-invoke convenience.

## What wraps what

Parsing (the genuinely hard, error-prone part — comment/continuation
handling, section-header grammar, the systemd 2048-byte line-length ceiling)
is delegated entirely to
[`coreos/go-systemd/v22/unit`](https://pkg.go.dev/github.com/coreos/go-systemd/v22/unit)
(Apache-2.0) — CoreOS/Red Hat's own unit-file lexer, the same one used
throughout the container/cloud ecosystem. Its `DeserializeSections` groups
directives by section while genuinely preserving repeats; this package adds
the systemd-directive semantic knowledge on top (which directives mean what,
how they group, which values are lists vs. single values vs. lists-of-
whitespace-tokens).

## Safety

The unit file is always supplied as text by the caller: no `systemctl`, no
D-Bus, no filesystem, no network access, no wall-clock, no randomness. Input
text is capped at 1 MB, checked before any parsing work. A malformed unit
file (unterminated section header, garbage after a section header, a
directive line missing `=`) always returns a structured error, never a
crash.

## License

MIT. Built for the Axiom marketplace.
