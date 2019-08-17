# execx
[![GoDoc](https://godoc.org/github.com/izumin5210/execx?status.svg)](https://godoc.org/github.com/izumin5210/execx)
[![License](https://img.shields.io/github/license/izumin5210/execx.svg)](./LICENSE)

Wrapper of `os/exec` to stop commands correctly.

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

cmd := execx.CommandContext(ctx, "sh", "-c", "sleep 5; echo done")
out, err := cmd.Output()

st := err.(*execx.ExitStatus)

fmt.Println(out, err, st.Signaled, st.Killed)

// Output: [] context deadline exceeded true false
```

## Reference

- https://github.com/Songmu/timeout
