package playlist

import (
	"testing"
)

func TestPickDistribution(t *testing.T) {
	seen := map[string]int{}
	const runs = 200
	for i := 0; i < runs; i++ {
		e := Pick()
		seen[e.Name]++
	}

	t.Logf("Track distribution over %d picks:", runs)
	for _, track := range Tracks {
		count := seen[track.Name]
		t.Logf("  %-38s %3d picks (%.0f%%)", track.Name, count, float64(count)/runs*100)
	}

	// Every track with assets should appear at least once in 200 picks
	for name, count := range seen {
		if count == 0 {
			t.Errorf("track %q never picked in %d runs", name, runs)
		}
	}

	if len(seen) < 2 {
		t.Errorf("only %d unique track(s) selected — randomisation not working", len(seen))
	}
}
