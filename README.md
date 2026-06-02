# constgen

Generate type-safe constants and helper functions for multiple programming languages from a single YAML definition.

**Quick Start**

- **Generate:** `go run ./cmd/constgen -cmd=generate` (reads YAML from [yaml](yaml) and writes to [generated](generated)).
- **Build:** `go build ./cmd/constgen` then run the binary with `-cmd=generate -input=yaml -out=generated`.

**YAML Format**

- **Name:** top-level `name` (used as filename base and class/object name). See [yaml/BooleanCodes.yaml](yaml/BooleanCodes.yaml).
- **Profiles:** `profiles` is a map of profile-name → constants. Profile names must match `[a-z][a-z0-9_]*`.
  - Constants are a map of UPPERCASE keys → scalar values (string, int, float, bool).
  - To add generated helper functions for a profile, include a `functions` map inside that profile. Example:

  ```yaml
  profiles:
  	string:
  		YES: "YES"
  		NO: "NO"
  		functions:
  			isYes: [YES]
  			isNo: [NO]
  ```

**Targets**

- Configure which languages to generate under `targets`. Each target must set `enabled: true` and list `profiles` to generate.
- Optional per-target function signature overrides live in `function_signatures`. They can be either:
  - A full function body/definition (emitted verbatim), or
  - A function signature which the generator will wrap with a default body.

Example (override `isNo` for PHP, Go, Java, JS):

```yaml
targets:
	php:
		enabled: true
		profiles: [string]
		function_signatures:
			isNo: "public static function isNo($value) { return in_array(strtoupper($value), [self::NO, 'N', 'n']); }"
	go:
		enabled: true
		profiles: [string]
		function_signatures:
			isNo: "func IsNo(value string) bool { v := strings.ToUpper(value); return v == NO || v == \"N\" || v == \"n\" }"
```

**How functions work**

- The `functions` map in a profile defines helper functions by name → list of constant keys they should compare against (or a single key string).
- By default for string-like constants the generator produces a helper that uppercases the input and compares to the constants.
- For numeric-like constants the helper performs equality checks.
- If a target provides a `function_signatures` entry for a function:
  - If the string contains a full function definition (e.g. contains `{` or `def `/`fun `), it is emitted verbatim and the generator skips the default body.
  - Otherwise the string is treated as a signature and the generator will place the default body inside it.

**Generated output layout**

- Files are written to: `generated/<language>/<profile>/<FileName>` where `FileName` is derived from `name` in the spec (preserves case).
- Generators attempt deterministic output by sorting constant keys.
