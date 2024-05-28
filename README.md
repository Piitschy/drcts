# Directus Data Model CLI
 Directus Data Model CLI (drcts) to migrate schemas from one instance to an other.


![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Piitschy/drcts)
![GitHub Release](https://img.shields.io/github/v/release/Piitschy/drcts)
![GitHub Release Date](https://img.shields.io/github/release-date/Piitschy/drcts)
![GitHub branch check runs](https://img.shields.io/github/check-runs/Piitschy/drcts/main)

## Installation

```bash
go install github.com/Piitschy/drcts@latest
```

## Usage

### Migrate

```bash
drcts --bu <base-url> --bt <base-token> --tu <target-url> --tt <target-token> migrate
```
 or just run `drcts migrate` and follow the instructions.

### Export

To export the schema of a Directus instance to a file, run:

```bash
drcts --bu <base-url> --bt <base-token> save -o <output-file>
```

Formats supported: `json`, `yaml`, `csv`, `xml`. But only `json` is appliable.

### Apply

To apply a schema from a file to a Directus instance, run:

```bash
drcts --tu <target-url> --tt <target-token> apply -i <input-file>
```

Its only possible to apply a schema in `json` format.

### Diff 

To compare the schema of two Directus instances, run:

```bash
drcts --bu <base-url> --bt <base-token> --tu <target-url> --tt <target-token> save-diff -o <diff-output-file>
```
or
```bash
drcts  --tu <target-url> --tt <target-token> save-diff -i <base-schema-file> -o <diff-output-file>
```

or just run `drcts save-diff` and follow the instructions.

Formats supported: `json`, `yaml`, `csv`, `xml`. But only `json` is appliable.

### Apply Diff

To apply a schema diff from a file to a Directus instance, run:

```bash
drcts --tu <target-url> --tt <target-token> apply-diff -i <diff-file>
```

or just run `drcts apply-diff` and follow the instructions.

Its only possible to apply a schema diff in `json` format.

## Environment Variables

You can also use environment variables to set the Directus instance URL and token:

- `DRCTS_BASE_URL`
- `DRCTS_BASE_TOKEN`
- `DRCTS_TARGET_URL`
- `DRCTS_TARGET_TOKEN`


- `DRCTS_SCHEMA_FILE` (for `save` and `apply`)
- `DRCTS_DIFF_FILE` (for `save-diff` and `apply-diff`)

