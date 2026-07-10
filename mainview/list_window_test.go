package mainview

import "testing"

func TestVisibleRange(t *testing.T) {
	tests := []struct {
		name          string
		selected      int
		total         int
		visible       int
		expectedStart int
		expectedEnd   int
	}{
		{
			name:          "empty list",
			selected:      0,
			total:         0,
			visible:       5,
			expectedStart: 0,
			expectedEnd:   0,
		},
		{
			name:          "all items fit",
			selected:      1,
			total:         3,
			visible:       5,
			expectedStart: 0,
			expectedEnd:   3,
		},
		{
			name:          "selected near top",
			selected:      1,
			total:         10,
			visible:       5,
			expectedStart: 0,
			expectedEnd:   5,
		},
		{
			name:          "selected near bottom",
			selected:      9,
			total:         10,
			visible:       5,
			expectedStart: 5,
			expectedEnd:   10,
		},
		{
			name:          "selected scrolls window",
			selected:      6,
			total:         10,
			visible:       5,
			expectedStart: 2,
			expectedEnd:   7,
		},
		{
			name:          "no visible rows",
			selected:      2,
			total:         10,
			visible:       0,
			expectedStart: 0,
			expectedEnd:   0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			start, end := visibleRange(test.selected, test.total, test.visible)
			if start != test.expectedStart || end != test.expectedEnd {
				t.Fatalf("expected %d..%d, got %d..%d", test.expectedStart, test.expectedEnd, start, end)
			}
		})
	}
}
