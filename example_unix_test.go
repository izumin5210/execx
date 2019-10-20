// +build !windows

package execx_test

import (
	"context"
	"fmt"
	"time"

	"github.com/izumin5210/execx"
)

func ExampleCommandContext() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := execx.CommandContext(ctx, "sh", "-c", "sleep 5; echo done")
	out, err := cmd.Output()

	st := err.(*execx.ExitStatus)

	fmt.Println(out, err, st.Signaled, st.Killed)

	// Output: [] signal: terminated true false
}
