package sum

import "testing"

func TestSum(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{1, 2, 3}, 6},
		{"empty", []int{}, 0},
		{"mixed", []int{-1, 1, 2}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.nums...); got != tt.want {
				t.Fatalf("Sum(%v) = %d; want %d", tt.nums, got, tt.want)
			}
		})
	}
}
