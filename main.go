package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	usage = `jsonencoder - A CLI tool to encode and decode JSON strings

Usage:
  %s [options] <command> <input>

Commands:
  encode    Encode JSON (escape for embedding)
  decode    Decode JSON (unescape)

Options:
  -f, --file    Read input from file instead of command line argument
  -h, --help    Show this help message

Examples:
  %s encode '{"key": "value"}'
  %s decode '"{\"key\": \"value\"}"'
  %s encode -f input.json
  %s decode -f encoded.json
`
)

func main() {
	var fileInput bool
	flag.BoolVar(&fileInput, "f", false, "Read input from file")
	flag.BoolVar(&fileInput, "file", false, "Read input from file")

	flag.Usage = func() {
		progName := os.Args[0]
		fmt.Fprintf(os.Stderr, usage, progName, progName, progName, progName, progName)
	}

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 || (len(args) < 2 && !fileInput) {
		flag.Usage()
		os.Exit(1)
	}

	command := args[0]
	var input string

	if len(args) > 1 {
		input = args[1]
	}

	var jsonData string
	var err error

	if fileInput {
		if input == "" {
			fmt.Fprintf(os.Stderr, "Error: file name required when using -f flag\n")
			os.Exit(1)
		}
		jsonData, err = readFromFile(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
	} else {
		if input == "" {
			fmt.Fprintf(os.Stderr, "Error: JSON input required\n")
			os.Exit(1)
		}
		jsonData = input
	}

	switch strings.ToLower(command) {
	case "encode":
		result, err := encodeJSON(jsonData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(result)
	case "decode":
		result, err := decodeJSON(jsonData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(result)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		flag.Usage()
		os.Exit(1)
	}
}

// readFromFile reads the entire content of a file
func readFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

// encodeJSON takes a JSON string and encodes it for safe embedding
// This validates the JSON and then marshals it as a string
func encodeJSON(jsonStr string) (string, error) {
	// First, validate and minify the input JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		return "", fmt.Errorf("invalid JSON input: %v", err)
	}

	// Marshal the input as minified JSON (no extra whitespace)
	minified, err := json.Marshal(jsonData)
	if err != nil {
		return "", fmt.Errorf("failed to minify JSON: %v", err)
	}

	// Use strconv.Quote to escape special characters for safe embedding
	quoted := strconv.Quote(string(minified))
	return quoted, nil
}

// decodeJSON takes an encoded JSON string and decodes it
func decodeJSON(encodedStr string) (string, error) {
	var decoded string
	if err := json.Unmarshal([]byte(encodedStr), &decoded); err != nil {
		return "", fmt.Errorf("failed to decode JSON: %v", err)
	}

	// Validate that the decoded result is valid JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(decoded), &jsonData); err != nil {
		return "", fmt.Errorf("decoded result is not valid JSON: %v", err)
	}

	return decoded, nil
}
