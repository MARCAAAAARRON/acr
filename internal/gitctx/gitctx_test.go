package gitctx

import "testing"

func TestParsePorcelainSingleFile(t *testing.T) {
	input := " M src/ratelimiter/ratelimiter.guard.ts\n"
	files := parsePorcelain(input)

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d: %v", len(files), files)
	}
	if files[0] != "src/ratelimiter/ratelimiter.guard.ts" {
		t.Errorf("expected full correct path, got %q", files[0])
	}
}