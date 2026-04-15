package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"bizzmod-cli/internal/api"
	"github.com/spf13/cobra"
)

type adminContractField struct {
	Name     string   `json:"name"`
	Label    string   `json:"label"`
	Type     string   `json:"type"`
	Required bool     `json:"required"`
	Default  any      `json:"default"`
	Enum     []string `json:"enum"`
	Help     string   `json:"help"`
}

type adminContract struct {
	Resource             string                    `json:"resource"`
	Action               string                    `json:"action"`
	Method               string                    `json:"method"`
	PathTemplate         string                    `json:"path_template"`
	PathParams           []adminContractField      `json:"path_params"`
	Fields               []adminContractField      `json:"fields"`
	PayloadTemplate      map[string]any            `json:"payload_template"`
	Taxonomy             map[string]any            `json:"taxonomy"`
	ConfigDefaultsByType map[string]map[string]any `json:"config_defaults_by_type"`
}

func runAdminInteractiveContractAction(cmd *cobra.Command, resource, action string) error {
	client, err := mustClient()
	if err != nil {
		return err
	}

	contract, err := fetchAdminContract(client, resource, action)
	if err != nil {
		return err
	}

	pathParamValues, payload, err := promptPayloadFromContract(cmd, contract)
	if err != nil {
		return err
	}

	path := applyPathTemplate(contract.PathTemplate, pathParamValues)
	method := strings.ToUpper(strings.TrimSpace(contract.Method))
	if method == "" {
		if action == "create" {
			method = "POST"
		} else {
			method = "PUT"
		}
	}

	body := []byte{}
	if method != "GET" && method != "DELETE" {
		encoded, marshalErr := json.Marshal(payload)
		if marshalErr != nil {
			return fmt.Errorf("failed to serialize payload: %w", marshalErr)
		}
		body = encoded
	}

	respBody, status, reqErr := client.Do(method, path, body)
	return printResponse(cmd, respBody, status, reqErr)
}

func fetchAdminContract(client *api.Client, resource, action string) (adminContract, error) {
	path := fmt.Sprintf("/api/v1/external/admin/contracts/%s/%s", resource, action)
	body, status, err := client.Do("GET", path, nil)
	if err != nil {
		return adminContract{}, fmt.Errorf("failed to fetch contract: %w", err)
	}

	var parsed map[string]any
	if unmarshalErr := json.Unmarshal(body, &parsed); unmarshalErr != nil {
		return adminContract{}, fmt.Errorf("invalid contract response: %w", unmarshalErr)
	}

	data := parsed["data"]
	if data == nil {
		data = parsed
	}

	encoded, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		return adminContract{}, fmt.Errorf("invalid contract payload: %w", marshalErr)
	}

	var contract adminContract
	if unmarshalErr := json.Unmarshal(encoded, &contract); unmarshalErr != nil {
		return adminContract{}, fmt.Errorf("invalid contract schema: %w", unmarshalErr)
	}

	if strings.TrimSpace(contract.PathTemplate) == "" {
		return adminContract{}, fmt.Errorf("invalid contract schema: missing path_template")
	}

	_ = status
	return contract, nil
}

func promptPayloadFromContract(cmd *cobra.Command, contract adminContract) (map[string]string, map[string]any, error) {
	reader := bufio.NewReader(cmd.InOrStdin())

	pathValues := make(map[string]string)
	for _, field := range contract.PathParams {
		value, _, err := promptContractFieldValue(reader, cmd, field, nil)
		if err != nil {
			return nil, nil, err
		}
		pathValues[field.Name] = fmt.Sprintf("%v", value)
	}

	payload := deepCopyMap(contract.PayloadTemplate)
	for _, field := range contract.Fields {
		defaultValue := field.Default
		if existing, ok := payload[field.Name]; ok {
			defaultValue = existing
		}

		value, provided, err := promptContractFieldValue(reader, cmd, field, defaultValue)
		if err != nil {
			return nil, nil, err
		}

		if provided {
			payload[field.Name] = value
		}
	}

	applyContractDerivedValues(contract, payload)
	return pathValues, payload, nil
}

func promptContractFieldValue(reader *bufio.Reader, cmd *cobra.Command, field adminContractField, defaultValue any) (any, bool, error) {
	label := strings.TrimSpace(field.Label)
	if label == "" {
		label = field.Name
	}

	typeName := strings.TrimSpace(field.Type)
	if typeName == "" {
		typeName = "string"
	}

	for {
		hints := []string{typeName}
		if field.Required {
			hints = append(hints, "required")
		} else {
			hints = append(hints, "optional")
		}
		if len(field.Enum) > 0 {
			hints = append(hints, "options: "+strings.Join(field.Enum, ", "))
		}
		if field.Help != "" {
			hints = append(hints, field.Help)
		}

		defaultText := ""
		if defaultValue != nil {
			defaultText = fmt.Sprintf(" [default: %s]", toCompactJSON(defaultValue))
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s (%s)%s: ", label, strings.Join(hints, " | "), defaultText)
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, false, err
		}
		line = strings.TrimSpace(line)

		if line == "" {
			if defaultValue != nil {
				return defaultValue, true, nil
			}
			if field.Required {
				fmt.Fprintln(cmd.OutOrStdout(), "value is required")
				continue
			}
			return nil, false, nil
		}

		if len(field.Enum) > 0 && !containsString(field.Enum, line) {
			fmt.Fprintln(cmd.OutOrStdout(), "invalid option")
			continue
		}

		switch strings.ToLower(typeName) {
		case "boolean", "bool":
			parsed, ok := parseBool(line)
			if !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "expected boolean value (true/false, yes/no, 1/0)")
				continue
			}
			return parsed, true, nil
		case "number", "int", "integer":
			if strings.Contains(line, ".") {
				parsed, parseErr := strconv.ParseFloat(line, 64)
				if parseErr != nil {
					fmt.Fprintln(cmd.OutOrStdout(), "expected numeric value")
					continue
				}
				return parsed, true, nil
			}
			parsed, parseErr := strconv.Atoi(line)
			if parseErr != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "expected numeric value")
				continue
			}
			return parsed, true, nil
		case "json", "object", "array":
			var parsed any
			if parseErr := json.Unmarshal([]byte(line), &parsed); parseErr != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "expected valid JSON")
				continue
			}
			return parsed, true, nil
		default:
			return line, true, nil
		}
	}
}

func applyContractDerivedValues(contract adminContract, payload map[string]any) {
	if contract.Resource != "credentials" {
		return
	}

	taxonomyRaw, hasTaxonomy := contract.Taxonomy["backend_category_by_subcategory"]
	if hasTaxonomy {
		if categoryMap, ok := taxonomyRaw.(map[string]any); ok {
			subcategory := fmt.Sprintf("%v", payload["subcategory"])
			if mapped, found := categoryMap[subcategory]; found {
				payload["category"] = mapped
			}
		}
	}

	credentialType := fmt.Sprintf("%v", payload["type"])
	currentConfig := payload["config"]
	if defaults, ok := contract.ConfigDefaultsByType[credentialType]; ok {
		if currentConfigMap, isMap := currentConfig.(map[string]any); !isMap || len(currentConfigMap) == 0 {
			payload["config"] = deepCopyMap(defaults)
		}
	}

	delete(payload, "main_category")
	delete(payload, "subcategory")
}

func applyPathTemplate(pathTemplate string, pathParams map[string]string) string {
	path := pathTemplate
	for key, value := range pathParams {
		path = strings.ReplaceAll(path, "{"+key+"}", value)
	}
	return path
}

func deepCopyMap(input map[string]any) map[string]any {
	if input == nil {
		return map[string]any{}
	}

	encoded, err := json.Marshal(input)
	if err != nil {
		copy := make(map[string]any, len(input))
		for key, value := range input {
			copy[key] = value
		}
		return copy
	}

	copy := map[string]any{}
	if err := json.Unmarshal(encoded, &copy); err != nil {
		return map[string]any{}
	}
	return copy
}

func toCompactJSON(value any) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	return string(encoded)
}

func containsString(options []string, value string) bool {
	for _, option := range options {
		if option == value {
			return true
		}
	}
	return false
}

func parseBool(input string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "true", "1", "yes", "y", "si", "s":
		return true, true
	case "false", "0", "no", "n":
		return false, true
	default:
		return false, false
	}
}
