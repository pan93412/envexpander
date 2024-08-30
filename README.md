<!-- markdownlint-disable MD033 -->

# envexpander

Expand the unordered referenced variable in the environment variables efficiently.

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
pkg: github.com/pan93412/envexpander/v3
BenchmarkEnvExpand_V3/simple-16                  2946846               407.5 ns/op
BenchmarkEnvExpand_V3/simple-2-16                2984750               398.1 ns/op
BenchmarkEnvExpand_V3/complex-16                  863078              1422 ns/op
BenchmarkEnvExpand_V3/very-complex-16             482428              2456 ns/op
BenchmarkEnvExpand_V3/unknown-reference-16               4462561               271.5 ns/op
BenchmarkEnvExpand_V3/unknown-reference-2-16             2693028               446.9 ns/op
BenchmarkEnvExpand_V3/realcase-1-16                       404270              2791 ns/op
BenchmarkEnvExpand_V3/simple-3-16                        1742944               698.6 ns/op
BenchmarkEnvExpand_V3/selfreference-16                   4691149               253.2 ns/op
BenchmarkRefV3/simple-16                                57546478                20.14 ns/op
BenchmarkRefV3/simple-2-16                              56987947                20.82 ns/op
BenchmarkRefV3/simple-3-16                              29251585                40.99 ns/op
BenchmarkRefV3/simple-4-16                              28097368                41.83 ns/op
BenchmarkRefV3/invalid-syntax-1-16                      227145477                5.252 ns/op
BenchmarkRefV3/invalid-syntax-2-16                      180704725                6.537 ns/op
BenchmarkRefV3/invalid-syntax-3-16                      47285127                23.72 ns/op
BenchmarkRefV3/invalid-syntax-4-16                      53172633                22.21 ns/op
BenchmarkRefV3/invalid-syntax-5-16                      45558663                26.01 ns/op
BenchmarkRefV3/escape-16                                46021900                23.86 ns/op
```

## Test Coverage

```plain
ok      github.com/pan93412/envexpander/v3      0.183s  coverage: 98.6% of statements
```

## Used in

- [zeabur/zbpack](https://github.com/zeabur/zbpack)

## Examples

- [expander.go](./examples/expander.go): Receive the dotenv-like format from stdin and output the expanded result to stdout.

## License

MIT License. See [LICENSE](./LICENSE) file.
