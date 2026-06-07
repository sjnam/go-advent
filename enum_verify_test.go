package main

import "testing"

func TestEnumValues(t *testing.T) {
	checks := []struct {
		name      string
		got, want int
	}{
		{"NOWHERE", int(NOWHERE), 74}, {"OFFICE", int(OFFICE), 73}, {"Y2", int(Y2), 69}, {"N", int(N), 0},
		{"CHAIN", int(CHAIN), 66}, {"GOLD", int(GOLD), 51}, {"WATER", int(WATER), 24}, {"BIRD", int(BIRD), 10}, {"KEYS", int(KEYS), 1},
		{"QUIT", int(QUIT), 29}, {"RELAX", int(RELAX), 10}, {"TAKE", int(TAKE), 1}, {"ABSTAIN", int(ABSTAIN), 0},
		{"troll", int(troll), 143}, {"didit", int(didit), 140}, {"inside", int(inside), 9},
		{"emist", int(emist), 15}, {"crack", int(crack), 129}, {"y2", int(y2), 54}, {"road", int(road), 1},
		{"inhand", int(inhand), -1}, {"limbo", int(limbo), 0},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %d, want %d", c.name, c.got, c.want)
		}
	}
}
