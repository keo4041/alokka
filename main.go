package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	input := `
{
  "number_1": {
    "N": "1.50"
  },
  "string_1": {
    "S": "784498 "
  },
  "string_2": {
    "S": "2014-07-16T20:55:46Z"
  },
  "map_1": {
    "M": {
      "bool_1": {
        "BOOL": "truthy"
      },
      "null_1": {
        "NULL ": "true"
      },
      "list_1": {
        "L": [
          {
            "S": ""
          },
          {
            "N": "011"
          },
          {
            "N": "5215s"
          },
          {
            "BOOL": "f"
          },
          {
            "NULL": "0"
          }
        ]
      }
    }
  },
  "list_2": {
    "L": "noop"
  },
  "list_3": {
    "L": [
      "noop"
    ]
  },
  "": {
    "S": "noop"
  }
}
`

	var inputData map[string]map[string]interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		fmt.Println("Error parsing input JSON:", err)
		return
	}

	output := transform(inputData)
	outputJSON, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling output JSON:", err)
		return
	}

	fmt.Println(string(outputJSON))
}

func transform(input map[string]map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	transformed := make(map[string]interface{})

	for key, value := range input {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		for dataType, val := range value {
			var transformedValue interface{}
			switch dataType {
			case "S":
				transformedValue = transformString(val.(string))
			case "N":
				transformedValue = transformNumber(val.(string))
			case "BOOL":
				transformedValue = transformBool(val.(string))
			case "NULL":
				transformedValue = transformNull(val.(string))
			case "L":
				transformedValue = transformList(val)
			case "M":
				transformedValue = transformMap(val.(map[string]interface{}))
			}

			if transformedValue != nil {
				if dataType == "NULL" {
					transformed[key] = nil

				} else {
					transformed[key] = transformedValue

				}
			}
		}
	}

	result = append(result, transformed)
	return result
}

func transformString(val string) interface{} {
	val = strings.TrimSpace(val)
	if val == "" {
		return nil
	}

	if t, err := time.Parse(time.RFC3339, val); err == nil {
		return t.Unix()
	}

	return val
}

func transformNumber(val string) interface{} {
	val = strings.TrimSpace(val)
	if val == "" {
		return nil
	}

	if num, err := strconv.ParseFloat(val, 64); err == nil {
		return num
	}

	return nil
}

func transformBool(val string) interface{} {
	val = strings.TrimSpace(strings.ToLower(val))
	switch val {
	case "1", "t", "true":
		return true
	case "0", "f", "false":
		return false
	default:
		return nil
	}
}

func transformNull(val string) interface{} {
	val = strings.TrimSpace(strings.ToLower(val))
	switch val {
	case "1", "t", "true":
		return "true"
	default:
		return nil
	}
}

func transformList(vals interface{}) interface{} {
	valSlice, ok := vals.([]interface{})
	if !ok {
		return nil // Omit if the value is not a slice
	}
	var result []interface{}
	for _, v := range valSlice {
		if m, ok := v.(map[string]interface{}); ok {
			for dataType, val := range m {
				var transformed interface{}
				switch dataType {
				case "S":
					transformed = transformString(val.(string))
				case "N":
					transformed = transformNumber(val.(string))
				case "BOOL":
					transformed = transformBool(val.(string))
				case "NULL":
					transformed = transformNull(val.(string))
				}
				if transformed != nil {
					if dataType == "NULL" {
						result = append(result, nil)

					} else {
						result = append(result, transformed)

					}
				}
			}
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func transformMap(input map[string]interface{}) interface{} {
	transformed := make(map[string]interface{})
	for key, value := range input {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		if m, ok := value.(map[string]interface{}); ok {
			for dataType, val := range m {
				dataType = strings.TrimSpace(dataType)
				var transformedValue interface{}
				switch dataType {
				case "S":
					transformedValue = transformString(val.(string))
				case "N":
					transformedValue = transformNumber(val.(string))
				case "BOOL":
					transformedValue = transformBool(val.(string))
				case "NULL":
					transformedValue = transformNull(val.(string))
				case "L":
					transformedValue = transformList(val)
				case "M":
					transformedValue = transformMap(val.(map[string]interface{}))
				}

				if transformedValue != nil {
					if dataType == "NULL" {
						transformed[key] = nil

					} else {
						transformed[key] = transformedValue

					}
				}
			}
		}
	}
	if len(transformed) == 0 {
		return nil
	}
	return transformed
}
