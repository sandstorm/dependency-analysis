package parsing

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

// Matches a file line by line and provides the value of the first capturing
// regex group of the first match in the first matching line.
func getFirstLineMatchInPath(filePath string, regexp *regexp.Regexp) (string, error) {
	reader, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	result := getFirstLineMatchInReader(reader, regexp)
	return result, nil
}

// Matches a reader line by line and provides the value of the first capturing
// regex group of the first match in the first matching line.
func getFirstLineMatchInReader(reader io.Reader, regexp *regexp.Regexp) string {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	return getFirstLineMatchInScanner(scanner, regexp)
}

// Matches a scanner item by item and provides the value of the first capturing
// regex group of the first match in the first matching line.
func getFirstLineMatchInScanner(scanner *bufio.Scanner, regexp *regexp.Regexp) string {
	for scanner.Scan() {
		line := scanner.Text()
		match := regexp.FindStringSubmatch(line)
		if match != nil {
			return match[1]
		}
	}
	return ""
}

// Finds all matches of the given regex in the given content and returns
// the first capturing regex group of each match.
func getAllMatches(content string, regexp *regexp.Regexp) []string {
	matches := regexp.FindAllStringSubmatch(content, -1)
	result := make([]string, len(matches))
	for i, v := range matches {
		result[i] = v[1]
	}
	return result
}

// reads all bytes into a string
func readerToString(reader io.Reader) (string, error) {
	buffer := new(strings.Builder)
	_, err := io.Copy(buffer, reader)
	if err != nil {
		return "", err
	}
	content := buffer.String()
	return content, nil
}

// splits each element using the given delimiter
func splitAll(values []string, delimiter string) [][]string {
	return mapElements(values, func(e string) []string { return strings.Split(e, delimiter) })
}

// applies f to all elements of in and returns the result
func mapElements[TIn any, TOut any](in []TIn, f func(TIn) TOut) []TOut {
	out := make([]TOut, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}
