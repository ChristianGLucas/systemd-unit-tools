// Package-internal helpers shared by every node: text parsing (wrapping
// coreos/go-systemd/v22/unit), the text-or-pre-parsed-unit resolution
// convenience, directive lookups, the Environment= quoted-word tokenizer,
// and unit-type detection. No node calls systemctl, D-Bus, the filesystem,
// the network, the wall clock, or randomness — everything here is a pure
// function of the caller-supplied text/UnitFile.
package nodes

import (
	"fmt"
	"strings"
	"unicode"

	sdunit "github.com/coreos/go-systemd/v22/unit"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// parseUnitText parses raw unit-file text into the canonical UnitFile
// envelope, delegating the actual lexing/grammar (comments, backslash line
// continuations, the systemd 2048-byte line-length ceiling, "garbage after
// section header", "option before section") to
// github.com/coreos/go-systemd/v22/unit.DeserializeSections. Repeated
// section headers of the same name are merged, in first-appearance order;
// directives are never deduplicated or overwritten.
func parseUnitText(text string) (*gen.UnitFile, error) {
	secs, err := sdunit.DeserializeSections(strings.NewReader(text))
	if err != nil {
		return nil, err
	}

	var sections []*gen.Section
	index := map[string]int{}
	for _, s := range secs {
		i, ok := index[s.Section]
		if !ok {
			sections = append(sections, &gen.Section{Name: s.Section})
			i = len(sections) - 1
			index[s.Section] = i
		}
		for _, e := range s.Entries {
			sections[i].Directives = append(sections[i].Directives, &gen.Directive{Key: e.Name, Value: e.Value})
		}
	}
	return &gen.UnitFile{Sections: sections}, nil
}

// resolveUnit implements the shared input-envelope convenience every node
// uses: a pre-parsed `unit` (non-empty Sections) wins when supplied —
// avoiding re-parsing text that was already parsed by ParseUnitFile earlier
// in a flow — otherwise `text` is parsed fresh. Returns a non-empty error
// string (never a Go error/panic) when neither input is usable.
func resolveUnit(text string, u *gen.UnitFile) (*gen.UnitFile, string) {
	if u != nil && len(u.Sections) > 0 {
		return u, ""
	}
	if strings.TrimSpace(text) == "" {
		return nil, "no input provided: supply either 'text' (raw unit-file content) or a pre-parsed 'unit'"
	}
	parsed, err := parseUnitText(text)
	if err != nil {
		return nil, err.Error()
	}
	return parsed, ""
}

// findSection returns the first (and, since parseUnitText merges repeated
// headers, only) section with the given name.
func findSection(u *gen.UnitFile, name string) (*gen.Section, bool) {
	for _, s := range u.Sections {
		if s.Name == name {
			return s, true
		}
	}
	return nil, false
}

func hasSection(u *gen.UnitFile, name string) bool {
	_, ok := findSection(u, name)
	return ok
}

// directiveValues returns every value of `key` within `section`, in source
// order — the repeated-key-safe lookup every typed extraction node builds
// on. Returns nil (not an error) when the section or key is absent.
func directiveValues(u *gen.UnitFile, section, key string) []string {
	s, ok := findSection(u, section)
	if !ok {
		return nil
	}
	var vals []string
	for _, d := range s.Directives {
		if d.Key == key {
			vals = append(vals, d.Value)
		}
	}
	return vals
}

// tokenizeAll whitespace-splits every value in `values` and concatenates
// the tokens in order — used for directives whose value is a
// space-separated list of unit/target names (After=, WantedBy=, ...), where
// the directive itself may also legitimately repeat.
func tokenizeAll(values []string) []string {
	var out []string
	for _, v := range values {
		out = append(out, strings.Fields(v)...)
	}
	return out
}

// splitQuotedWords tokenizes a systemd.syntax "quoted assignment list"
// value (as used by Environment=) on whitespace, treating a double-quoted
// span as one token even if it contains spaces, and interpreting the
// escape sequences systemd.syntax(7) documents for this context: \" \\ \'
// \s \t \n \r. Any other backslash-escaped character is left as a literal
// two-character "\x" sequence rather than guessed at.
func splitQuotedWords(s string) []string {
	var out []string
	var buf strings.Builder
	inQuote := false
	hasToken := false

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case r == '\\' && i+1 < len(runes):
			next := runes[i+1]
			switch next {
			case '"':
				buf.WriteRune('"')
			case '\\':
				buf.WriteRune('\\')
			case '\'':
				buf.WriteRune('\'')
			case 's':
				buf.WriteRune(' ')
			case 't':
				buf.WriteRune('\t')
			case 'n':
				buf.WriteRune('\n')
			case 'r':
				buf.WriteRune('\r')
			default:
				buf.WriteRune('\\')
				buf.WriteRune(next)
			}
			hasToken = true
			i++ // consumed `next` too
		case r == '"':
			inQuote = !inQuote
			hasToken = true
		case unicode.IsSpace(r) && !inQuote:
			if hasToken {
				out = append(out, buf.String())
				buf.Reset()
				hasToken = false
			}
		default:
			buf.WriteRune(r)
			hasToken = true
		}
	}
	if hasToken {
		out = append(out, buf.String())
	}
	return out
}

// parseEnvironmentEntries expands every Environment= occurrence into
// individual KEY=VALUE pairs via splitQuotedWords. A token with no "="
// is not a valid assignment and is skipped (best-effort extraction, not a
// hard parse error).
func parseEnvironmentEntries(values []string) []*gen.EnvVar {
	var out []*gen.EnvVar
	for _, v := range values {
		for _, tok := range splitQuotedWords(v) {
			eq := strings.IndexByte(tok, '=')
			if eq < 0 {
				continue
			}
			out = append(out, &gen.EnvVar{Key: tok[:eq], Value: tok[eq+1:]})
		}
	}
	return out
}

// parseEnvironmentFiles expands every EnvironmentFile= occurrence, peeling
// off the leading "-" ("ignore if missing") into `optional`.
func parseEnvironmentFiles(values []string) []*gen.EnvironmentFileEntry {
	var out []*gen.EnvironmentFileEntry
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		optional := false
		if strings.HasPrefix(v, "-") {
			optional = true
			v = v[1:]
		}
		out = append(out, &gen.EnvironmentFileEntry{Path: v, Optional: optional})
	}
	return out
}

// sectionToType maps a type-DISTINGUISHING section name to its unit type.
// Sections not listed here (Unit, Install, and vendor "X-*" sections) carry
// no type signal on their own — that is precisely what makes a
// [Unit]+[Install]-only file ambiguous without a filename (it could be a
// .target, or a .slice/.scope with no body properties).
var sectionToType = map[string]string{
	"Service":   "service",
	"Socket":    "socket",
	"Mount":     "mount",
	"Automount": "automount",
	"Swap":      "swap",
	"Path":      "path",
	"Timer":     "timer",
	"Slice":     "slice",
	"Scope":     "scope",
}

// sectionTypePriority is the deterministic tie-break order when more than
// one type-distinguishing section is present (not a valid real-world unit,
// but the detector must still answer deterministically).
var sectionTypePriority = []string{"Service", "Socket", "Mount", "Automount", "Swap", "Path", "Timer", "Slice", "Scope"}

// extensionToType maps a recognized unit filename extension to its type.
var extensionToType = map[string]string{
	".service":   "service",
	".socket":    "socket",
	".mount":     "mount",
	".automount": "automount",
	".swap":      "swap",
	".path":      "path",
	".timer":     "timer",
	".slice":     "slice",
	".scope":     "scope",
	".target":    "target",
}

// detectUnitType implements the shared logic behind DetectUnitType and the
// unit_type field on ValidateUnit/SummarizeUnit. filename may be "".
func detectUnitType(u *gen.UnitFile, filename string) (unitType string, basis string) {
	if filename != "" {
		ext := lowerExt(filename)
		if t, ok := extensionToType[ext]; ok {
			return t, "filename_extension"
		}
	}
	present := map[string]bool{}
	for _, s := range u.Sections {
		present[s.Name] = true
	}
	for _, name := range sectionTypePriority {
		if present[name] {
			return sectionToType[name], "section_presence"
		}
	}
	return "", "ambiguous"
}

// lowerExt returns the lowercased extension (including the leading '.') of
// a filename, without importing path/filepath for one operation.
func lowerExt(filename string) string {
	i := strings.LastIndexByte(filename, '.')
	if i < 0 {
		return ""
	}
	return strings.ToLower(filename[i:])
}

// knownSections is the recognized systemd unit-file top-level section
// vocabulary. A section outside this set that also doesn't start with the
// conventional "X-" vendor-extension prefix is flagged by ValidateUnit.
var knownSections = map[string]bool{
	"Unit": true, "Install": true,
	"Service": true, "Socket": true, "Mount": true, "Automount": true,
	"Swap": true, "Path": true, "Timer": true, "Slice": true, "Scope": true,
}

// validServiceTypes is the documented vocabulary for [Service] Type=.
var validServiceTypes = map[string]bool{
	"simple": true, "forking": true, "oneshot": true, "notify": true,
	"notify-reload": true, "dbus": true, "exec": true, "idle": true,
}

// hardeningAllowlist is the resource-limit / sandboxing directive
// vocabulary GetHardeningDirectives filters [Service] against.
var hardeningAllowlist = map[string]bool{
	"MemoryLimit": true, "MemoryMax": true, "MemoryHigh": true, "MemoryLow": true,
	"MemoryMin": true, "MemorySwapMax": true,
	"CPUQuota": true, "CPUWeight": true, "CPUShares": true, "CPUAffinity": true,
	"TasksMax":    true,
	"LimitNOFILE": true, "LimitNPROC": true, "LimitCPU": true, "LimitAS": true,
	"LimitCORE": true, "LimitFSIZE": true, "LimitDATA": true, "LimitSTACK": true,
	"LimitRSS": true, "LimitMEMLOCK": true,
	"IOWeight": true, "IOReadBandwidthMax": true, "IOWriteBandwidthMax": true,
	"ProtectSystem": true, "ProtectHome": true, "ProtectKernelTunables": true,
	"ProtectKernelModules": true, "ProtectKernelLogs": true, "ProtectControlGroups": true,
	"ProtectClock": true, "ProtectHostname": true, "ProtectProc": true, "ProcSubset": true,
	"PrivateTmp": true, "PrivateNetwork": true, "PrivateDevices": true,
	"PrivateUsers": true, "PrivateMounts": true, "PrivateIPC": true,
	"NoNewPrivileges": true,
	"ReadOnlyPaths":   true, "ReadWritePaths": true, "InaccessiblePaths": true,
	"TemporaryFileSystem": true, "BindPaths": true, "BindReadOnlyPaths": true,
	"CapabilityBoundingSet": true, "AmbientCapabilities": true,
	"SystemCallFilter": true, "SystemCallArchitectures": true,
	"SystemCallErrorNumber": true, "SystemCallLog": true,
	"RestrictNamespaces": true, "RestrictRealtime": true, "RestrictSUIDSGID": true,
	"RestrictAddressFamilies": true,
	"LockPersonality":         true, "MemoryDenyWriteExecute": true,
	"KeyringMode": true, "UMask": true,
	"DeviceAllow": true, "DevicePolicy": true,
	"IPAddressAllow": true, "IPAddressDeny": true,
}

// validateUnit implements ValidateUnit's structural checks against an
// already-parsed UnitFile. Static analysis only: it never invokes systemd,
// so it cannot and does not check whether a referenced unit actually
// exists on any system.
func validateUnit(u *gen.UnitFile) (issues []*gen.ValidationIssue, valid bool, unitType string) {
	hasError := false
	addErr := func(section, msg string) {
		issues = append(issues, &gen.ValidationIssue{Severity: "error", Section: section, Message: msg})
		hasError = true
	}
	addWarn := func(section, msg string) {
		issues = append(issues, &gen.ValidationIssue{Severity: "warning", Section: section, Message: msg})
	}

	if len(u.Sections) == 0 {
		addErr("", "unit file has no sections")
	}
	for _, s := range u.Sections {
		if !knownSections[s.Name] && !strings.HasPrefix(s.Name, "X-") {
			addWarn(s.Name, fmt.Sprintf("unrecognized section [%s]", s.Name))
		}
	}
	if !hasSection(u, "Unit") {
		addWarn("", "no [Unit] section present")
	}

	unitType, _ = detectUnitType(u, "")

	if svc, ok := findSection(u, "Service"); ok {
		hasExec := false
		declaredType := ""
		for _, d := range svc.Directives {
			if d.Key == "ExecStart" {
				hasExec = true
			}
			if d.Key == "Type" && declaredType == "" {
				declaredType = d.Value
			}
		}
		if !hasExec {
			addWarn("Service", "no ExecStart directive found (valid only for an ordering-only oneshot service with RemainAfterExit=yes)")
		}
		if declaredType != "" && !validServiceTypes[declaredType] {
			addErr("Service", fmt.Sprintf("unrecognized Type value %q", declaredType))
		}
	}
	if sock, ok := findSection(u, "Socket"); ok {
		hasListen := false
		for _, d := range sock.Directives {
			if strings.HasPrefix(d.Key, "Listen") {
				hasListen = true
				break
			}
		}
		if !hasListen {
			addWarn("Socket", "no Listen* directive found")
		}
	}
	if tmr, ok := findSection(u, "Timer"); ok {
		hasSchedule := false
		for _, d := range tmr.Directives {
			if strings.HasPrefix(d.Key, "On") {
				hasSchedule = true
				break
			}
		}
		if !hasSchedule {
			addWarn("Timer", "no On*= schedule directive found")
		}
	}

	return issues, !hasError, unitType
}
