package exec_test

import (
	"context"
	"fmt"
	"time"

	"github.com/izumin5210/exec"
)

func ExampleCommandTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", "sleep 5; echo done")
	out, err := cmd.Output()

	st := err.(*exec.ExitStatus)

	fmt.Println(out, err, st.Signaled, st.Killed)

	// Output: [] context deadline exceeded true false
}
