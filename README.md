# Koda Academy - Weekly 3

A Go program for ordering and paying for restaurant menu items.
A port of [Weekly 1](../weekly1) using a Bubble Tea + Huh TUI.

## Running

```sh
go run .
```

## Test

```sh
go test ./...
```

## Known issues

`huh` v2.0.3 hides the last option of a multi-select that has a title: it
derives the viewport height from the options alone, then subtracts the title
height anyway. `Select` guards this case, `MultiSelect` does not.

`optionsForm` in [model/main.go](model/main.go) works around it by passing an
explicit height. The upstream fix is in
[patches/](patches/huh-v2.0.3-multiselect-height.patch) — Go does not apply
patch files, so it is kept only to send upstream, and applies to `huh` itself:

```sh
git am < patches/huh-v2.0.3-multiselect-height.patch
```

Drop the workaround once a `huh` release includes the fix.
`TestOptionsFormShowsEveryOption` will fail if it is removed too early.

## License

[MIT](LICENSE)
