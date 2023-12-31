<!-- markdownlint-disable MD033 -->

# envexpander

Expand the unordered referenced variable in the environment variables.

## How does it work?

Expansion is a process of replacing a variable reference with its value.

For example:

<table>
<tr>
    <th>Input</th>
    <th>Output</th>
</tr><tr>
<td>

```env
A=B
B=${A}
```

</td>
<td>

```env
A=B
B=B
```

</td>
</tr>
</table>

This tool can expand the variable reference `${A}` to `B` in the above example.

It is NOT MATTER for this tool whether the definition is ordered. So the following example is also expandable.

<table>
<tr>
    <th>Input</th>
    <th>Output</th>
</tr><tr>
<td>

```env
A=${B}
B=A123
```

</td>
<td>

```env
A=A123
B=A123
```

</td>
</tr>
</table>

For complex examples, it also handles well.

<table>
<tr>
    <th>Input</th>
    <th>Output</th>
</tr><tr>
<td>

```env
CONTACT_MAIL=foo@bar.tld
LISTEN_HOST=http://0.0.0.0:${PORT}
DATABASE_URI=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable
PORT=8080
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_HOST=internal.postgres.host
POSTGRES_PORT=1145
POSTGRES_DB=testdb
```

</td>
<td>

```env
CONTACT_MAIL=foo@bar.tld
LISTEN_HOST=http://0.0.0.0:8080
DATABASE_URI=postgres://postgres:postgres@internal.postgres.host:1145/testdb?sslmode=disable
PORT=8080
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_HOST=internal.postgres.host
POSTGRES_PORT=1145
POSTGRES_DB=testdb
```

</td>
</tr>
</table>

For circular references, we break it by replacing it with an empty string.

<table>
<tr>
    <th>Input</th>
    <th>Output</th>
</tr><tr>
<td>

```env
A=${B}C
B=C${A}
```

</td>
<td>

```env
A=C
B=C
```

</td>
</tr>
</table>

For undefined variables, we leave it as it is.

<table>
<tr>
    <th>Input</th>
    <th>Output</th>
</tr><tr>
<td>

```env
A=${CCC}
```

</td>
<td>

```env
A=${CCC}
```

</td>
</tr>
</table>

## Benchmark

```plain
goos: darwin
goarch: arm64
pkg: github.com/pan93412/envexpander/v2
BenchmarkFindVariableReferenceMap-8              9632380               111.0 ns/op           192 B/op           2 allocs/op
BenchmarkResolveEnvVariable_Basic-8               261889              4229 ns/op            4298 B/op          46 allocs/op
BenchmarkResolveEnvVariable_Complex-8             452300              2588 ns/op            3248 B/op          37 allocs/op
BenchmarkReplacerIntegrate-8                     6312372               190.6 ns/op           168 B/op           6 allocs/op
BenchmarkReplacerIntegrateWithCache-8           10643762               107.8 ns/op            56 B/op           3 allocs/op
```

## Test Coverage

```plain
ok      github.com/pan93412/envexpander/v2      0.587s  coverage: 100.0% of statements
```

## Used in

- [zeabur/zbpack](https://github.com/zeabur/zbpack)

## Examples

- [expander.go](./examples/expander.go): Receive the dotenv-like format from stdin and output the expanded result to stdout. Note that `envexpander` does not include the dotenv parser, so you need to parse it by yourself (or reference `expander.go`'s example :D).

## License

MIT License. See [LICENSE](./LICENSE) file.
