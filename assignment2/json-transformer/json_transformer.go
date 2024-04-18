package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func sanitizeValue(value string) string {
	// Sanitize the value of trailing and leading whitespace.
	return strings.TrimSpace(value)
}

func transformString(value string) interface{} {
	sanitizedValue := sanitizeValue(value)
	// Transform RFC3339 formatted Strings to Unix Epoch in Numeric data type.
	layout := "2006-01-02T15:04:05Z"
	if _, err := time.Parse(layout, sanitizedValue); err == nil {
		t, _ := time.Parse(layout, sanitizedValue)
		return t.Unix()
	} else if sanitizedValue == "" {
		// Omit fields with empty values.
		return nil
	}

	return sanitizedValue
}

func transformNumber(value string) interface{} {
	sanitizedValue := sanitizeValue(value)
	// Omit fields with invalid Numeric values.
	if val, err := strconv.ParseFloat(sanitizedValue, 64); err == nil {
		// Strip the leading zeros.
		if strings.Contains(sanitizedValue, ".") {
			return val
		}
		return int(val)
	}
	return nil
}

func transformBoolean(value string) interface{} {
	sanitizedValue := sanitizeValue(value)
	// Omit fields with invalid Boolean values.
	switch sanitizedValue {
	case "1", "t", "T", "TRUE", "true", "True":
		return true
	case "0", "f", "F", "FALSE", "false", "False":
		return false
	}
	return nil
}

func transformNull(value string) interface{} {
	sanitizedValue := sanitizeValue(value)
	// Omit fields with invalid Boolean values.
	if sanitizedValue == "1" || strings.ToLower(sanitizedValue) == "true" {
		return nil
	}
	return sanitizedValue
}

func transformList(value map[string]interface{}) []interface{} {
	// Transform value to the List data type.
	if val, ok := value["L"]; ok {
		if list, ok := val.([]interface{}); ok {
			return list
		}
	}
	return nil
}

func transformMap(value map[string]interface{}) map[string]interface{} {
	// Transform value to the Map data type.
	result := make(map[string]interface{})
	for key, val := range value {
		if key == "M" {
			if m, ok := val.(map[string]interface{}); ok {
				for k, v := range m {
					result[k] = v
				}
			}
		} else {
			result[key] = val
		}
	}
	return result
}

func jsonTransformer(inputJSON map[string]map[string]interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for key, item := range inputJSON {
		sanitizedKey := sanitizeValue(key)
		if sanitizedKey != "" { // Omit fields with empty keys.
			transformedValue := make(map[string]interface{})
			for dataType, val := range item {
				switch dataType {
				case "S":
					transformedValue[sanitizedKey] = transformString(val.(string))
				case "N":
					transformedValue[sanitizedKey] = transformNumber(val.(string))
				case "BOOL":
					transformedValue[sanitizedKey] = transformBoolean(val.(string))
				case "NULL":
					transformedValue[sanitizedKey] = transformNull(val.(string))
				case "L":
					transformedValue[sanitizedKey] = transformList(val.(map[string]interface{}))
				case "M":
					transformedValue[sanitizedKey] = transformMap(val.(map[string]interface{}))
				}
			}
			if len(transformedValue) > 0 {
				result = append(result, transformedValue)
			}
		}
	}
	return result
}

func main() {
	// Specify the file path
	filePath := "input.json"

	// Loading the contents of json file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Unmarshaling JSON content into a map
	var jsonDict map[string]map[string]interface{}
	if err := json.Unmarshal(content, &jsonDict); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Transforming the JSON input
	transformedJSON := jsonTransformer(jsonDict)

	// Marshaling the transformed JSON into a byte slice
	transformedBytes, err := json.MarshalIndent(transformedJSON, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Print the transformed JSON
	fmt.Println(string(transformedBytes))
}
