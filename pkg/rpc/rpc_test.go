package rpc

import (
	"os"
	"testing"

	"github.com/asymmetric-research/solana-exporter/pkg/slog"
)

func TestMain(m *testing.M) {
	slog.Init()
	code := m.Run()
	os.Exit(code)
}
