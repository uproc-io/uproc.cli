package processes

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"text/tabwriter"
)

func isValidOverviewSection(section string) bool {
	switch strings.ToLower(strings.TrimSpace(section)) {
	case "all", "kpis", "charts", "tables":
		return true
	default:
		return false
	}
}

func formatModuleOverviewOutput(body []byte, moduleSlug, section string) (string, error) {
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil {
		return "", err
	}

	data, ok := root["data"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("missing overview data")
	}

	resolvedSection := strings.ToLower(strings.TrimSpace(section))
	if resolvedSection == "" {
		resolvedSection = "all"
	}

	parts := []string{fmt.Sprintf("Overview: %s", moduleSlug)}

	if resolvedSection == "all" || resolvedSection == "kpis" {
		parts = append(parts, renderOverviewKpis(data["kpis"]))
	}
	if resolvedSection == "all" || resolvedSection == "charts" {
		parts = append(parts, renderOverviewCharts(data))
	}
	if resolvedSection == "all" || resolvedSection == "tables" {
		parts = append(parts, renderOverviewTables(data["tables"]))
	}

	return strings.TrimSpace(strings.Join(parts, "\n\n")), nil
}

func renderOverviewKpis(kpisRaw any) string {
	header := "KPIs"
	kpisMap, ok := kpisRaw.(map[string]any)
	if !ok || len(kpisMap) == 0 {
		return header + "\n(empty)"
	}

	keys := make([]string, 0, len(kpisMap))
	for key := range kpisMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var b strings.Builder
	tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "key\tvalue")
	for _, key := range keys {
		_, _ = fmt.Fprintf(tw, "%s\t%v\n", key, kpisMap[key])
	}
	_ = tw.Flush()

	return header + "\n" + strings.TrimRight(b.String(), "\n")
}

func renderOverviewCharts(data map[string]any) string {
	header := "Charts"
	charts := normalizeCharts(data)
	if len(charts) == 0 {
		return header + "\n(empty)"
	}

	parts := []string{header}
	for _, chart := range charts {
		title := nonEmpty(chart["title"], chart["id"], "chart")
		chartType := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", chart["type"])))
		parts = append(parts, fmt.Sprintf("- %s (%s)", title, nonEmpty(chartType, "chart")))

		if isDonutChart(chartType, chart) {
			parts = append(parts, indentLines(renderDonutChart(chart), 1))
			continue
		}

		if isTimeSeriesChart(chartType) {
			if line := renderLineChart(chart); line != "" {
				parts = append(parts, indentLines(line, 1))
				continue
			}
		}

		if isBarLikeChart(chartType) {
			if bar := renderBarChart(chart); bar != "" {
				parts = append(parts, indentLines(bar, 1))
				continue
			}
		}

		if line := renderLineChart(chart); line != "" {
			parts = append(parts, indentLines(line, 1))
			continue
		}

		if bar := renderBarChart(chart); bar != "" {
			parts = append(parts, indentLines(bar, 1))
			continue
		}

		parts = append(parts, "  (empty)")
	}

	return strings.Join(parts, "\n")
}

func renderOverviewTables(tablesRaw any) string {
	header := "Tables"
	tables, ok := tablesRaw.([]any)
	if !ok || len(tables) == 0 {
		return header + "\n(empty)"
	}

	parts := []string{header}
	for _, tableRaw := range tables {
		table, ok := tableRaw.(map[string]any)
		if !ok {
			continue
		}

		tableName := nonEmpty(table["title"], table["id"], "table")
		parts = append(parts, fmt.Sprintf("- %s", tableName))

		rows, _ := table["rows"].([]any)
		if len(rows) == 0 {
			parts = append(parts, "  (empty)")
			continue
		}

		renderedRows := rows
		if len(renderedRows) > 20 {
			renderedRows = renderedRows[:20]
		}

		tableText := renderObjectTable(renderedRows, 1)
		parts = append(parts, tableText)
		if len(rows) > len(renderedRows) {
			parts = append(parts, fmt.Sprintf("  ... %d more rows", len(rows)-len(renderedRows)))
		}
	}

	return strings.Join(parts, "\n")
}

func normalizeCharts(data map[string]any) []map[string]any {
	chartsRaw, hasCharts := data["charts"]
	if hasCharts {
		if charts := toChartList(chartsRaw); len(charts) > 0 {
			return charts
		}
	}
	return nil
}

func toChartList(raw any) []map[string]any {
	list, ok := raw.([]any)
	if !ok {
		return nil
	}

	out := make([]map[string]any, 0, len(list))
	for _, item := range list {
		mapped, ok := item.(map[string]any)
		if ok {
			out = append(out, mapped)
		}
	}
	return out
}

func isDonutChart(chartType string, chart map[string]any) bool {
	if strings.Contains(chartType, "donut") || strings.Contains(chartType, "pie") {
		return true
	}
	_, labelsOk := extractChartLabels(chart)
	series := extractChartSeries(chart)
	if !labelsOk || len(series) == 0 {
		return false
	}
	if len(series) == 1 {
		name := strings.ToLower(strings.TrimSpace(fmt.Sprintf("%v", series[0]["name"])))
		return strings.Contains(name, "distribution") || strings.Contains(name, "share")
	}
	return false
}

func renderDonutChart(chart map[string]any) string {
	labels, ok := extractChartLabels(chart)
	series := extractChartSeries(chart)
	if !ok || len(series) == 0 {
		return "(empty)"
	}

	values := extractNumericSeriesData(series[0])
	if len(values) == 0 || len(labels) == 0 {
		return "(empty)"
	}

	count := len(values)
	if len(labels) < count {
		count = len(labels)
	}

	total := 0.0
	for i := 0; i < count; i++ {
		total += values[i]
	}
	if total <= 0 {
		return "(empty)"
	}

	indices := make([]int, count)
	for i := 0; i < count; i++ {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		return values[indices[i]] > values[indices[j]]
	})

	lines := make([]string, 0, count+1)
	for _, idx := range indices {
		pct := (values[idx] / total) * 100
		bar := scaledBar(pct, 20)
		lines = append(lines, fmt.Sprintf("%s\t%.2f\t(%.1f%%)\t%s", labels[idx], values[idx], pct, bar))
	}

	var b strings.Builder
	tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "segment\tvalue\tshare\tvisual")
	for _, line := range lines {
		_, _ = fmt.Fprintln(tw, line)
	}
	_, _ = fmt.Fprintf(tw, "Total\t%.2f\t\t\n", total)
	_ = tw.Flush()

	return strings.TrimRight(b.String(), "\n")
}

func renderBarChart(chart map[string]any) string {
	labels, ok := extractChartLabels(chart)
	if !ok || len(labels) == 0 {
		return ""
	}
	series := extractChartSeries(chart)
	if len(series) == 0 {
		return ""
	}

	values := extractNumericSeriesData(series[0])
	if len(values) == 0 {
		return ""
	}
	count := len(values)
	if len(labels) < count {
		count = len(labels)
	}
	if count == 0 {
		return ""
	}

	maxValue := 0.0
	for i := 0; i < count; i++ {
		if values[i] > maxValue {
			maxValue = values[i]
		}
	}
	if maxValue <= 0 {
		return ""
	}

	lines := make([]string, 0, count)
	for i := 0; i < count; i++ {
		share := (values[i] / maxValue) * 100
		lines = append(lines, fmt.Sprintf("%s\t%.2f\t%s", labels[i], values[i], scaledBar(share, 20)))
	}

	var b strings.Builder
	tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "label\tvalue\tvisual")
	for _, line := range lines {
		_, _ = fmt.Fprintln(tw, line)
	}
	_ = tw.Flush()
	return strings.TrimRight(b.String(), "\n")
}

func renderLineChart(chart map[string]any) string {
	series := extractChartSeries(chart)
	if len(series) == 0 {
		return ""
	}

	labels, labelsOk := extractChartLabels(chart)
	rangeSummary := ""
	if labelsOk && len(labels) > 1 {
		rangeSummary = fmt.Sprintf("range: %s -> %s (%d points)", labels[0], labels[len(labels)-1], len(labels))
	} else if labelsOk && len(labels) == 1 {
		rangeSummary = fmt.Sprintf("range: %s (1 point)", labels[0])
	}

	lines := make([]string, 0, len(series)*7)
	for _, serie := range series {
		name := nonEmpty(serie["name"], "series")
		values := extractNumericSeriesData(serie)
		if len(values) == 0 {
			continue
		}
		if len(values) > 40 {
			values = downsampleSeries(values, 40)
		}

		lines = append(
			lines,
			fmt.Sprintf("series: %s  min=%.2f  max=%.2f  last=%.2f", name, minFloat(values), maxFloat(values), values[len(values)-1]),
		)

		areaRows := renderSeriesArea(values, 4)
		for _, row := range areaRows {
			lines = append(lines, "  "+row)
		}
		lines = append(lines, "  "+sparkline(values))
	}

	if len(lines) == 0 {
		return ""
	}

	var b strings.Builder
	if rangeSummary != "" {
		_, _ = fmt.Fprintln(&b, rangeSummary)
	}
	for _, line := range lines {
		_, _ = fmt.Fprintln(&b, line)
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderSeriesArea(values []float64, height int) []string {
	if len(values) == 0 || height <= 0 {
		return nil
	}

	minV := minFloat(values)
	maxV := maxFloat(values)
	levels := make([]int, len(values))
	if maxV-minV == 0 {
		for i := range levels {
			levels[i] = height / 2
		}
	} else {
		for i, value := range values {
			normalized := (value - minV) / (maxV - minV)
			levels[i] = int(math.Round(normalized * float64(height)))
		}
	}

	rows := make([]string, 0, height)
	for row := height; row >= 1; row-- {
		var builder strings.Builder
		for _, level := range levels {
			if level >= row {
				builder.WriteRune('█')
			} else {
				builder.WriteRune(' ')
			}
		}
		rows = append(rows, builder.String())
	}

	return rows
}

func downsampleSeries(values []float64, maxPoints int) []float64 {
	if maxPoints <= 0 || len(values) <= maxPoints {
		return values
	}

	out := make([]float64, 0, maxPoints)
	step := float64(len(values)) / float64(maxPoints)
	for index := 0; index < maxPoints; index++ {
		start := int(math.Floor(float64(index) * step))
		end := int(math.Floor(float64(index+1) * step))
		if end <= start {
			end = start + 1
		}
		if end > len(values) {
			end = len(values)
		}

		sum := 0.0
		count := 0
		for i := start; i < end; i++ {
			sum += values[i]
			count++
		}
		if count == 0 {
			continue
		}
		out = append(out, sum/float64(count))
	}

	if len(out) == 0 {
		return values
	}

	return out
}

func extractChartLabels(chart map[string]any) ([]string, bool) {
	if labelsRaw, ok := chart["labels"].([]any); ok {
		return toStringSlice(labelsRaw), true
	}
	if value, ok := chart["value"].(map[string]any); ok {
		if labelsRaw, ok := value["labels"].([]any); ok {
			return toStringSlice(labelsRaw), true
		}
	}
	return nil, false
}

func extractChartSeries(chart map[string]any) []map[string]any {
	if seriesRaw, ok := chart["series"].([]any); ok {
		return normalizeSeries(seriesRaw)
	}
	if value, ok := chart["value"].(map[string]any); ok {
		if seriesRaw, ok := value["series"].([]any); ok {
			return normalizeSeries(seriesRaw)
		}
	}
	return nil
}

func normalizeSeries(seriesRaw []any) []map[string]any {
	series := make([]map[string]any, 0, len(seriesRaw))
	for idx, entry := range seriesRaw {
		switch typed := entry.(type) {
		case map[string]any:
			series = append(series, typed)
		case []any:
			series = append(series, map[string]any{"name": fmt.Sprintf("series_%d", idx+1), "data": typed})
		default:
			series = append(series, map[string]any{"name": fmt.Sprintf("series_%d", idx+1), "data": []any{typed}})
		}
	}
	return series
}

func extractNumericSeriesData(series map[string]any) []float64 {
	rawData, ok := series["data"].([]any)
	if !ok {
		if value, ok := series["value"].([]any); ok {
			rawData = value
		} else {
			return nil
		}
	}

	values := make([]float64, 0, len(rawData))
	for _, raw := range rawData {
		if parsed, ok := toFloat(raw); ok {
			values = append(values, parsed)
		}
	}
	return values
}

func toStringSlice(input []any) []string {
	out := make([]string, 0, len(input))
	for _, item := range input {
		out = append(out, fmt.Sprintf("%v", item))
	}
	return out
}

func toFloat(v any) (float64, bool) {
	switch value := v.(type) {
	case float64:
		return value, true
	case float32:
		return float64(value), true
	case int:
		return float64(value), true
	case int64:
		return float64(value), true
	case int32:
		return float64(value), true
	case int16:
		return float64(value), true
	case int8:
		return float64(value), true
	case uint:
		return float64(value), true
	case uint64:
		return float64(value), true
	case uint32:
		return float64(value), true
	case uint16:
		return float64(value), true
	case uint8:
		return float64(value), true
	case json.Number:
		parsed, err := value.Float64()
		if err == nil {
			return parsed, true
		}
	}
	return 0, false
}

func sparkline(values []float64) string {
	if len(values) == 0 {
		return ""
	}
	minV := minFloat(values)
	maxV := maxFloat(values)
	if maxV-minV == 0 {
		return strings.Repeat("▅", len(values))
	}

	chars := []rune("▁▂▃▄▅▆▇█")
	out := make([]rune, 0, len(values))
	for _, v := range values {
		normalized := (v - minV) / (maxV - minV)
		idx := int(math.Round(normalized * float64(len(chars)-1)))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(chars) {
			idx = len(chars) - 1
		}
		out = append(out, chars[idx])
	}
	return string(out)
}

func scaledBar(percent float64, width int) string {
	if width <= 0 {
		width = 10
	}
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := int(math.Round((percent / 100.0) * float64(width)))
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

func minFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	minV := values[0]
	for _, v := range values[1:] {
		if v < minV {
			minV = v
		}
	}
	return minV
}

func maxFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	maxV := values[0]
	for _, v := range values[1:] {
		if v > maxV {
			maxV = v
		}
	}
	return maxV
}

func indentLines(text string, indent int) string {
	pad := strings.Repeat("  ", indent)
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	for i, line := range lines {
		lines[i] = pad + line
	}
	return strings.Join(lines, "\n")
}

func nonEmpty(values ...any) string {
	for _, value := range values {
		text := strings.TrimSpace(fmt.Sprintf("%v", value))
		if text != "" && text != "<nil>" {
			return text
		}
	}
	return ""
}

func isTimeSeriesChart(chartType string) bool {
	return strings.Contains(chartType, "time-series") || strings.Contains(chartType, "line")
}

func isBarLikeChart(chartType string) bool {
	return strings.Contains(chartType, "bar") || strings.Contains(chartType, "column")
}
