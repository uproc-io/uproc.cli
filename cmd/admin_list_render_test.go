package cmd

import (
	"encoding/json"
	"testing"
)

func TestApplyVisibleFieldsFromContractFiltersRows(t *testing.T) {
	body := []byte(`{
		"success": true,
		"message": "ok",
		"data": {
			"rows": [
				{"id": 1, "name": "Acme", "country": "ES", "active": true, "extra": "x"}
			],
			"columns": [
				{"key": "id"}, {"key": "name"}, {"key": "country"}, {"key": "active"}, {"key": "extra"}
			]
		}
	}`)

	filtered, err := applyVisibleFieldsFromContract(body, []string{"id", "name", "country", "active"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(filtered, &parsed); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	data := parsed["data"].(map[string]any)
	rows := data["rows"].([]any)
	row := rows[0].(map[string]any)

	if _, ok := row["extra"]; ok {
		t.Fatalf("expected extra field to be removed")
	}
	if row["id"] != float64(1) {
		t.Fatalf("unexpected id value: %v", row["id"])
	}
	if row["name"] != "Acme" {
		t.Fatalf("unexpected name value: %v", row["name"])
	}
}
