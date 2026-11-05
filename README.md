# hase

Small fixed-point types for values with one or two decimal places, such as
temperatures, setpoints and sensor readings.

```go
go get github.com/kaatinga/hase
```

Requires Go 1.27 (generic methods).

## Why

Devices and wire formats rarely carry floats. A setpoint of 21.5 °C travels as
`215` (tenths) or `2150` (hundredths) in an `int16`. Code that handles these
values usually ends up as a scatter of bare integers and hand-written scaling:

```go
// Which unit is target? Tenths or hundredths?
func apply(target int16) {
	log.Printf("target %.2f", float32(target)/100)
	send(float32(target) / 100)
}

stored := int16(math.Round(float64(celsius) * 100))
```

With hase the scale is part of the type:

```go
func apply(target hase.CentiParameter) {
	log.Printf("target %s", target) // "target 21.50"
	send(target.Float[float32]())
}

stored := hase.NewCentiParameter(celsius)
```

What this buys you:

- **The unit is in the signature.** A function that takes a `CentiParameter`
  cannot be handed a `DeciParameter` or a raw count by accident; the compiler
  catches the mix-up that a bare `int16` lets through.
- **No more `/10`, `*100` and `math.Round` at every call site.** Conversion and
  rounding live in one place.
- **Printing is free.** Both types implement `fmt.Stringer`, so `%s`, `%v`,
  structured loggers and `html/template` (`{{ .Target }}`) show `21.50` or
  `21.5` without a format verb.
- **Same bytes on the wire.** The types are plain `int16` underneath: JSON,
  gob, binary encoding and database columns are unchanged, and `%d` still
  prints the raw integer.

## Types

| Type | Storage | Meaning | `String()` |
|---|---|---|---|
| `CentiParameter` | `int16` | value × 100, −327.68 … 327.67 | `"21.50"` |
| `DeciParameter` | `int16` | value × 10, −3276.8 … 3276.7 | `"21.5"` |
| `OffsetDeciParameter` *(experimental)* | `uint8` | (value − 5) × 10, 5.0 … 30.5 | `"21.5"` |

```go
c := hase.NewCentiParameter(21.456) // 2146, rounded to nearest
c.String()                          // "21.46"
c.Float[float64]()                  // 21.46

d := hase.NewDeciParameter(-3.25)   // -33
d.String()                          // "-3.3"
```

`New*` accepts any integer or float type.

### OffsetDeciParameter

A room setpoint between 5 and 30 °C with 0.1 °C steps fits in a single byte.
`OffsetDeciParameter` stores it that way, which halves the size of setpoint
tables, compact radio payloads and in-memory state. The encoding may still
change.

## Optional values

Protobuf and SQL often hand over `*int32` for an optional field. `PointerTo`
and `PointerFrom` convert between those and hase types, keeping `nil` as `nil`:

```go
var wire *int32 = hase.PointerTo[int32](reading.Temperature) // *CentiParameter → *int32
reading.Temperature = hase.PointerFrom[hase.CentiParameter](wire)
```

## Performance

`String()` avoids `strconv.FormatFloat` and the float round trip:

| Benchmark | ns/op | allocs |
|---|---|---|
| `DeciParameter.String` | 13 | 1 |
| `strconv.FormatFloat(…, 'f', 1, 64)` | 32 | 1 |
| `CentiParameter.String` | 13 | 1 |
| `strconv.FormatFloat(…, 'f', 2, 64)` | 39 | 1 |
| `OffsetDeciParameter.String` | 2 | 0 |

Apple M-series, Go 1.27. Run `go test -bench . -benchmem` to reproduce.

## Limits

- Values are assumed to be in range; `New*` does not clamp or report overflow.
- `New*` rounds to nearest. If an existing format truncates (`int16(x * 10)`),
  switching to hase changes the stored value — check before replacing.
