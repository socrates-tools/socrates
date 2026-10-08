package e2e_test

import (
	"os"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestMain(m *testing.M) {
	os.Exit(source.RunSharingBuilds(m))
}
