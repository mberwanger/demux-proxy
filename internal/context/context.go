package context

import (
	ctx "context"
	"os"
	"strings"
	"time"
)

type Env map[string]string

func (e Env) Copy() Env {
	var out = Env{}
	for k, v := range e {
		out[k] = v
	}
	return out
}

func (e Env) Strings() []string {
	var result = make([]string, 0, len(e))
	for k, v := range e {
		result = append(result, k+"="+v)
	}
	return result
}

type Context struct {
	ctx.Context
	BindAddress string
	Env         Env
	Date        time.Time
}

func New() *Context {
	return Wrap(ctx.Background())
}

func Wrap(ctx ctx.Context) *Context {
	return &Context{
		Context: ctx,
		Env:     splitEnv(os.Environ()),
		Date:    time.Now(),
	}
}

func splitEnv(env []string) map[string]string {
	r := map[string]string{}
	for _, e := range env {
		p := strings.SplitN(e, "=", 2)
		r[p[0]] = p[1]
	}
	return r
}
