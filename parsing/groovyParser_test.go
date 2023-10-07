package parsing

import (
	"bytes"
	"testing"
)

func TestParseGroovySourceUnit(t *testing.T) {
	testCases := []struct {
		name        string
		fileContent string
		expected    fullyQualifiedType
	}{
		{
			name: "simple class without imports",
			fileContent: `package de.sandstorm.test;
			
			public class Main {

			}`,
			expected: fullyQualifiedType{"de", "sandstorm", "test", "Main"},
		},
		{
			name: "simple class without imports with default visibility",
			fileContent: `package de.sandstorm.test;
			
			class Main {

			}`,
			expected: fullyQualifiedType{"de", "sandstorm", "test", "Main"},
		},
		{
			name: "simple class no semicolons",
			fileContent: `package de.sandstorm.test
			
			class Main {

			}`,
			expected: fullyQualifiedType{"de", "sandstorm", "test", "Main"},
		},
		{
			name: "simple class with imports",
			fileContent: `package de.sandstorm.test
			
			import static org.junit.Assert.assertEquals // some comment

			public class Main {

			}`,
			expected: fullyQualifiedType{"de", "sandstorm", "test", "Main"},
		},
		// TODO: test private static final class
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			file := bytes.NewBufferString(testCase.fileContent)
			AssertEquals(t,
				[]fullyQualifiedType{testCase.expected},
				ParseGroovySourceUnit(file),
			)
		})
	}
}

func TestParseGroovyImports(t *testing.T) {
	testCases := []struct {
		name        string
		fileContent string
		expected    []fullyQualifiedType
	}{
		{
			name: "simple class with mixed imports",
			fileContent: `
				package de.sandstorm.test;

				import java.util.List;
				import        de.sandstorm.greetings.HelloWorld;
				import        de.sandstorm.http.requests.GetGreetingRequest;
				import  static  org.junit.Assert.assertEquals; // some comment

				public class Main {

				}
			`,
			expected: []fullyQualifiedType{
				{"java", "util", "List"},
				{"de", "sandstorm", "greetings", "HelloWorld"},
				{"de", "sandstorm", "http", "requests", "GetGreetingRequest"},
				{"org", "junit", "Assert", "assertEquals"},
			},
		},
		{
			name: "simple class with mixed imports - no semicolons",
			fileContent: `
				package de.sandstorm.test

				import java.util.List
				import        de.sandstorm.greetings.HelloWorld
				import        de.sandstorm.http.requests.GetGreetingRequest
				import  static  org.junit.Assert.assertEquals // some comment

				class Main {

				}
			`,
			expected: []fullyQualifiedType{
				{"java", "util", "List"},
				{"de", "sandstorm", "greetings", "HelloWorld"},
				{"de", "sandstorm", "http", "requests", "GetGreetingRequest"},
				{"org", "junit", "Assert", "assertEquals"},
			},
		},
		// TODO: test with different parameters
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			file := bytes.NewBufferString(testCase.fileContent)
			actual, _ := ParseGroovyImports(
				file)
			AssertEquals(t,
				testCase.expected,
				actual)
		})
	}
}
