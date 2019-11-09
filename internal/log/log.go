package log

import (
	"context"
	"fmt"

	"github.com/mikemackintosh/twonicorn/config"
)

func Infof(ctx context.Context, msg string, params ...interface{}) {
	m := fmt.Sprintf(msg, params...)
	fmt.Printf("[INFO] Req:%s "+m+"\n", ctx.Value("reqid"))
}

func Debugf(ctx context.Context, msg string, params ...interface{}) {
	if config.IsDebug() {
		m := fmt.Sprintf(msg, params...)
		fmt.Printf("[DEBG] Req:%s "+m+"\n", ctx.Value("reqid"))
	}
}

func Printf(ctx context.Context, msg string, params ...interface{}) {
	m := fmt.Sprintf(msg, params...)
	fmt.Printf("[    ] Req:%s "+m+"\n", ctx.Value("reqid"))
}
