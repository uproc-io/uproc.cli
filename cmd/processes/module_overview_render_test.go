package processes

import (
	"strings"
	"testing"
)

func TestIsValidOverviewSection(t *testing.T) {
	valid := []string{"all", "kpis", "charts", "tables", "CHARTS"}
	for _, section := range valid {
		if !isValidOverviewSection(section) {
			t.Fatalf("expected %s to be valid", section)
		}
	}

	invalid := []string{"", "chart", "overview", "table"}
	for _, section := range invalid {
		if isValidOverviewSection(section) {
			t.Fatalf("expected %s to be invalid", section)
		}
	}
}

func TestFormatModuleOverviewOutputRendersSections(t *testing.T) {
	body := []byte(`{
		"success": true,
		"data": {
			"kpis": {"processed": 10, "failed": 2},
			"charts": [
				{
					"id": "chart1",
					"type": "line",
					"title": "Trend",
					"value": {
						"labels": ["D1", "D2", "D3"],
						"series": [{"name": "Processed", "data": [1, 3, 2]}]
					}
				}
			],
			"tables": [
				{"id": "recent", "rows": [{"id": 1, "name": "A"}]}
			]
		}
	}`)

	rendered, err := formatModuleOverviewOutput(body, "invoice-ingest", "all")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, rendered, "Overview: invoice-ingest")
	assertContains(t, rendered, "KPIs")
	assertContains(t, rendered, "Charts")
	assertContains(t, rendered, "series")
	assertContains(t, rendered, "series:")
	assertContains(t, rendered, "Tables")
}

func TestFormatModuleOverviewOutputWithoutChartsAndTablesShowsEmptySections(t *testing.T) {
	body := []byte(`{
		"success": true,
		"data": {
			"kpis": {"processed": 5}
		}
	}`)

	rendered, err := formatModuleOverviewOutput(body, "invoice-ingest", "all")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, rendered, "KPIs")
	assertContains(t, rendered, "Charts")
	assertContains(t, rendered, "(empty)")
	assertContains(t, rendered, "Tables")
}

func TestTimeSeriesChartRendersHorizontalSparkline(t *testing.T) {
	body := []byte(`{
		"success": true,
		"data": {
			"charts": [
				{
					"id": "chart1",
					"type": "time-series",
					"title": "Timeline",
					"value": {
						"labels": ["2026-04-01", "2026-04-02", "2026-04-03"],
						"series": [{"name": "processed", "data": [0, 3, 1]}]
					}
				}
			]
		}
	}`)

	rendered, err := formatModuleOverviewOutput(body, "invoice-ingest", "charts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, rendered, "range: 2026-04-01 -> 2026-04-03")
	assertContains(t, rendered, "processed")
	assertContains(t, rendered, "series: processed")
	assertContains(t, rendered, "min=0.00")
	assertContains(t, rendered, "max=3.00")
	assertContains(t, rendered, "█")
}

func TestFormatModuleOverviewOutputRendersDonut(t *testing.T) {
	body := []byte(`{
		"success": true,
		"data": {
			"kpis": {},
			"charts": [
				{
					"id": "by_status",
					"type": "donut",
					"title": "By status",
					"value": {
						"labels": ["open", "closed"],
						"series": [{"name": "distribution", "data": [30, 70]}]
					}
				}
			]
		}
	}`)

	rendered, err := formatModuleOverviewOutput(body, "support", "charts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, rendered, "By status")
	assertContains(t, rendered, "segment")
	assertContains(t, rendered, "open")
	assertContains(t, rendered, "closed")
	assertContains(t, rendered, "Total")
}

func assertContains(t *testing.T, text string, token string) {
	t.Helper()
	if !strings.Contains(text, token) {
		t.Fatalf("expected output to contain %q, got:\n%s", token, text)
	}
}
