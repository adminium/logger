package extract

import (
	"bytes"
	"testing"
	"time"

	"github.com/adminium/logger"
	"go.uber.org/zap/zapcore"
)

func TestNewCoreFormat(t *testing.T) {
	entry := zapcore.Entry{
		LoggerName: "main",
		Level:      zapcore.InfoLevel,
		Message:    "scooby",
		Time:       time.Date(2010, 5, 23, 15, 14, 0, 0, time.UTC),
	}

	testCases := []struct {
		format logger.LogFormat
		want   string
	}{
		{
			format: logger.ColorizedOutput,
			want:   "2010-05-23T15:14:00.000Z\t\x1b[34mINFO\x1b[0m\tmain\tscooby\n",
		},
		{
			format: logger.JSONOutput,
			want:   `{"level":"info","ts":"2010-05-23T15:14:00.000Z","logger":"main","msg":"scooby"}` + "\n",
		},
		{
			format: logger.PlaintextOutput,
			want:   "2010-05-23T15:14:00.000Z\tINFO\tmain\tscooby\n",
		},
	}

	for _, tc := range testCases {
		buf := &bytes.Buffer{}
		ws := zapcore.AddSync(buf)

		core := logger.NewCore(tc.format, ws, logger.LevelDebug)
		if err := core.Write(entry, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		row, err := Extract(buf.String())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		t.Log("a:", buf.String())
		t.Log("b:", row.Output())

		//got := buf.String()
		//if got != tc.want {
		//	t.Errorf("got %q, want %q", got, tc.want)
		//}
	}

}
