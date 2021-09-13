package trace

import (
	"bytes"
	"testing"
)

// TestNew tests the tracing behaviour.
func TestNew(t *testing.T) {

	var buf bytes.Buffer
	tracer := New(&buf)

	if tracer == nil {
		t.Error("Newからの戻り値がnillです")
	}
	tracer.Trace("こんにちは、traceパッケージ")
	if buf.String() != "こんにちは、traceパッケージ\n" {
		t.Errorf("'%s'といった謝った文字列が検出されました", buf.String())
	}
}
// ログを無効化させる。
func TestOff (t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("データ")
}