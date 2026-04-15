package cmd

import "testing"

func TestNormalizeOutputPayloadReturnsDataRows(t *testing.T) {
	input := map[string]any{
		"success": true,
		"message": "ok",
		"data": map[string]any{
			"rows": []any{map[string]any{"id": 1}},
		},
	}

	got := normalizeOutputPayload(input)
	rows, ok := got.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", got)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
}

func TestNormalizeOutputPayloadReturnsDataMapWithoutEnvelope(t *testing.T) {
	input := map[string]any{
		"success": true,
		"message": "ok",
		"data":    map[string]any{"id": 5, "name": "Ada"},
	}

	got := normalizeOutputPayload(input)
	mapped, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", got)
	}
	if mapped["id"] != 5.0 && mapped["id"] != 5 {
		t.Fatalf("expected id 5, got %v", mapped["id"])
	}
}

func TestNormalizeOutputPayloadUsesColumnOrderForRows(t *testing.T) {
	input := map[string]any{
		"success": true,
		"data": map[string]any{
			"columns": []any{
				map[string]any{"key": "id"},
				map[string]any{"key": "name"},
			},
			"rows": []any{
				map[string]any{"name": "Acme", "id": 1, "extra": "x"},
			},
		},
	}

	got := normalizeOutputPayload(input)
	rows, ok := got.([]any)
	if !ok || len(rows) != 1 {
		t.Fatalf("expected rows list, got %T", got)
	}
	rowMap, ok := rows[0].(map[string]any)
	if !ok {
		t.Fatalf("expected row map, got %T", rows[0])
	}
	if rowMap["id"] != 1 && rowMap["id"] != float64(1) {
		t.Fatalf("expected id in row")
	}
	if rowMap["name"] != "Acme" {
		t.Fatalf("expected name in row")
	}
}
