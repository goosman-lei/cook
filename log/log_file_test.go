package log

import (
	"testing"
	"time"
)

func TestCase_format_log_fname(t *testing.T) {
	var (
		r string
		d time.Time
	)

	r, d = parse_log_fname("cook.log-*{3}")
	t.Logf("%s %s", r, d)
	r, d = parse_log_fname("cook.log-*-*{3}")
	t.Logf("%s %s", r, d)
	r, d = parse_log_fname("cook.log-*-*-*{3}")
	t.Logf("%s %s", r, d)
	r, d = parse_log_fname("cook.log-*-*-*-*{3}")
	t.Logf("%s %s", r, d)
	r, d = parse_log_fname("cook.log-*-*-*-*-*{20}")
	t.Logf("%s %s", r, d)
	//t.Fail()
}

func Benchmark_format_log_fname(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parse_log_fname("cook.log-*-*-*-*{3}")
	}
}
