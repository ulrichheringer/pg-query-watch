package main

import (
	"strings"
	"testing"
)

func TestParseAndFilterStats(t *testing.T) {
	input := "{\"queryid\":\"12\",\"calls\":2,\"total_exec_time\":80,\"rows\":3,\"query\":\"select 1\"}\n{\"queryid\":\"13\",\"calls\":4,\"total_exec_time\":20,\"rows\":1,\"query\":\"select 2\"}\n"
	stats, err := ParseStats(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	slow := FilterSlow(stats, 30)
	if len(slow) != 1 || slow[0].QueryID != "12" || MeanMS(slow[0]) != 40 {
		t.Fatalf("unexpected results: %#v", slow)
	}
}
