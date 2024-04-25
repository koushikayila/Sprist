package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type JSONInput struct {
	Content map[string]interface{} `json:""`
}

type JSONOutput struct {
	Content map[string]interface{} `json:""`
}

func main() {
	processingStart := time.Now()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please provide the JSON file path: ")
	inputPath, _ := reader.ReadString('\n')
	cleanedPath := strings.TrimSpace(inputPath)
	result, processErr := processJSON(cleanedPath)
	if processErr != nil {
		fmt.Println(processErr)
		return
	}

	jsonEncoder := json.NewEncoder(os.Stdout)
	if encodeErr := jsonEncoder.Encode(result.Content); encodeErr != nil {
		fmt.Println(encodeErr)
	}

	processingTime := time.Since(processingStart)
	fmt.Printf("\nProcessing time: %s\n", processingTime)
}

func processJSON(path string) (JSONOutput, error) {
	fileContent, fileErr := os.ReadFile(path)
	if fileErr != nil {
		return JSONOutput{}, fileErr
	}

	var inputData map[string]interface{}
	if jsonErr := json.Unmarshal(fileContent, &inputData); jsonErr != nil {
		return JSONOutput{}, jsonErr
	}

	transformedData := JSONOutput{Content: make(map[string]interface{})}
	for key, val := range inputData {
		trimmedKey := strings.TrimSpace(key)
		if trimmedKey == "" {
			continue
		}

		if nestedMap, valid := val.(map[string]interface{}); valid {
			for dataType, data := range nestedMap {
				if dataStr, ok := data.(string); ok {
					cleanData := strings.TrimSpace(dataStr)
					if cleanData == "" {
						continue
					}

					if transformed, valid := processDataType(dataType, cleanData); valid {
						transformedData.Content[trimmedKey] = transformed
					}
				}
			}
		}
	}

	return transformedData, nil
}

func processDataType(typeKey, value string) (interface{}, bool) {
	switch typeKey {
	case "S":
		if parsedTime, err := time.Parse(time.RFC3339, value); err == nil {
			return parsedTime.Unix(), true
		} else {
			return value, true
		}
	case "N":
		if num, err := strconv.ParseFloat(value, 64); err == nil {
			return num, true
		}
	case "BOOL":
		if trueVals := map[string]bool{"1": true, "t": true, "true": true}; trueVals[strings.ToLower(value)] {
			return true, true
		}
		if falseVals := map[string]bool{"0": true, "f": true, "false": true}; falseVals[strings.ToLower(value)] {
			return false, true
		}
	case "NULL":
		if value == "null" {
			return nil, true
		}
	}

	return nil, false
}