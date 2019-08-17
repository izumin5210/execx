# exec
[![GoDoc](https://godoc.org/github.com/izumin5210/exec?status.svg)](https://godoc.org/github.com/izumin5210/exec)
[![License](https://img.shields.io/github/license/izumin5210/exec.svg)](./LICENSE)

Wrapper of `os/exec` to stop commands correctly.

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

cmd := exec.CommandContext(ctx, "sh", "-c", "sleep 5; echo done")
out, err := cmd.Output()

st := err.(*exec.ExitStatus)

fmt.Println(out, err, st.Signaled, st.Killed)

// Output: [] context deadline exceeded true false
```

## Reference

- https://github.com/Songmu/timeout
