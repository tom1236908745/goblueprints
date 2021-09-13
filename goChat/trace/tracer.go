package trace

import (
	"fmt"
	"io"
)

// Tracerはコード内での出来事を記録できるオブジェクトを表すインターフェース
type Tracer interface {
	Trace(...interface{})
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

// 小文字から始まるので、外部には公開されない。
type tracer struct {
	out io.Writer // 情報が出力される
}

func (t *tracer) Trace(a ...interface{}) {
	// fmt.Sprintを使用する事で、引数aを文字列に変換している。
	// そして、その結果がoutフィールド(io.Writer型)のメソッドに渡されている。
	// 文字列はbyte型に型変換されている。
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}
type nilTracer struct {}
func (t *nilTracer) Trace(a ...interface{}) {}

// OffはTraceメソッドの呼び出しを無視するTracerを返す。
func Off() Tracer {
	return &nilTracer{}
}
