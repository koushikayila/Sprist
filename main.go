package main

import (
	"bufio"        // Importing necessary packages
	"encoding/json" // Package for JSON encoding and decoding
	"fmt"           // Package for formatted I/O
	"os"            // Package for operating system functionalities
	"strconv"       // Package for string conversion
	"strings"       // Package for string manipulation
	"time"          // Package for time functionalities
)

// Struct to hold input JSON data
type JSONInput struct {
	Content map[string]interface{} `json:""`
}

// Struct to hold output JSON data
type JSONOutput struct {
	Content map[string]interface{} `json:""`
}

func main() {
	processingStart := time.Now() // Record the start time of processing

	reader := bufio.NewReader(os.Stdin) // Create a reader to read input from standard input
	fmt.Print("Please provide the JSON file path: ")
	inputPath, _ := reader.ReadString('\n')      // Read the input path from user
	cleanedPath := strings.TrimSpace(inputPath)  // Clean the input path by trimming spaces
	result, processErr := processJSON(cleanedPath) // Process the JSON file
	if processErr != nil { // If there is an error, print it and return
		fmt.Println(processErr)
		return
	}

	jsonEncoder := json.NewEncoder(os.Stdout) // Create a JSON encoder to encode the result to standard output
	if encodeErr := jsonEncoder.Encode(result.Content); encodeErr != nil { // Encode the result content and check for errors
		fmt.Println(encodeErr)
	}

	processingTime := time.Since(processingStart) // Calculate the processing time
	fmt.Printf("\nProcessing time: %s\n", processingTime) // Print the processing time
}

// Function to process the JSON file
func processJSON(path string) (JSONOutput, error) {
	fileContent, fileErr := os.ReadFile(path) // Read the file content
	if fileErr != nil { // If there is an error, return an empty JSONOutput and the error
		return JSONOutput{}, fileErr
	}

	var inputData map[string]interface{} // Declare a map to hold the input data
	if jsonErr := json.Unmarshal(fileContent, &inputData); jsonErr != nil { // Unmarshal the JSON content into the map
		return JSONOutput{}, jsonErr
	}

	transformedData := JSONOutput{Content: make(map[string]interface{})} // Create a JSONOutput with an empty map
	for key, val := range inputData { // Iterate through the input data
		trimmedKey := strings.TrimSpace(key) // Trim the key
		if trimmedKey == "" { // If the key is empty, continue to the next item
			continue
		}

		if nestedMap, valid := val.(map[string]interface{}); valid { // If the value is a nested map, proceed
			for dataType, data := range nestedMap { // Iterate through the nested map
				if dataStr, ok := data.(string); ok { // If the data is a string, proceed
					cleanData := strings.TrimSpace(dataStr) // Trim the data
					if cleanData == "" { // If the data is empty, continue to the next item
						continue
					}

					if transformed, valid := processDataType(dataType, cleanData); valid { // Process the data based on its type
						transformedData.Content[trimmedKey] = transformed // Add the transformed data to the result map
					}
				}
			}
		}
	}

	return transformedData, nil // Return the transformed data
}

// Function to process data based on its type
func processDataType(typeKey, value string) (interface{}, bool) {
	switch typeKey { // Switch based on the type key
	case "S": // If the type is string
		if parsedTime, err := time.Parse(time.RFC3339, value); err == nil { // Try to parse the string as time
			return parsedTime.Unix(), true // If successful, return the Unix timestamp
		} else { // If not, return the original string
			return value, true
		}
	case "N": // If the type is number
		if num, err := strconv.ParseFloat(value, 64); err == nil { // Try to parse the string as float
			return num, true // If successful, return the float
		}
	case "BOOL": // If the type is boolean
		if trueVals := map[string]bool{"1": true, "t": true, "true": true}; trueVals[strings.ToLower(value)] { // Check for true values
			return true, true // Return true if matched
		}
		if falseVals := map[string]bool{"0": true, "f": true, "false": true}; falseVals[strings.ToLower(value)] { // Check for false values
			return false, true // Return false if matched
		}
	case "NULL": // If the type is null
		if value == "null" { // Check if the value is "null"
			return nil, true // Return nil
		}
	}

	return nil, false // Return nil if none of the cases matched
}
