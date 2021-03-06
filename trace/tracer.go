package trace

import (
	"fmt"
	"io"
)

type Tracer interface {
	Trace(...interface{})
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	//t.out.Write([]byte(fmt.Sprint(a...)))
	//t.out.Write([]byte("\n"))

	//io.WriteString(t.out, fmt.Sprint(a...))
	//t.out.Write([]byte("\n"))

	fmt.Fprint(t.out, a...)
	fmt.Fprint(t.out, "\n")
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// Off はTraceメソッドの呼び出しを無視するTracerを返します
func Off() Tracer {
	return &nilTracer{}
}
