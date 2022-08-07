package parsing

import (
    "bytes"
	"reflect"
	"testing"
)

func AssertEquals(t *testing.T, expected interface{}, actual interface{}) {
	if reflect.DeepEqual(actual, expected) {
		return
	}
	t.Errorf("expected %v (type %v), received %v (type %v)",
		expected, reflect.TypeOf(expected),
		actual, reflect.TypeOf(actual))

}

func TestParseJavaSourceUnit(t *testing.T) {
	testCases := []struct {
		name        string
		fileContent string
		expected    []string
	}{
		{
			name: "simple class without imports",
			fileContent: `package de.sandstorm.test;
			
			public class Main {

			}`,
			expected: []string{"de", "sandstorm", "test", "Main"},
		},
		{
			name: "simple class with imports",
			fileContent: `package de.sandstorm.test;
			
			import static org.junit.Assert.assertEquals; // some comment

			public class Main {

			}`,
			expected: []string{"de", "sandstorm", "test", "Main"},
		},
		// TODO: test private static final class
		// TODO: implement and test commented classes
		// TODO: implement and test comments in general
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			file := bytes.NewBufferString(testCase.fileContent)
			AssertEquals(t,
				testCase.expected,
				ParseJavaSourceUnit(file),
			)
		})
	}
}

func TestParseJavaFile(t *testing.T) {
	testCases := []struct {
		name         string
		rootPackage  string
		segmentLimit int
		fileContent  string
		expected     []string
	}{
		{
			name: "simple class with mixed imports",
			rootPackage: "de.sandstorm",
			segmentLimit: 3,
			fileContent: `
				package de.sandstorm.test;

				import java.util.List;
				import        de.sandstorm.greetings.HelloWorld;
				import        de.sandstorm.http.requests.GetGreetingRequest;
				import  static  org.junit.Assert.assertEquals; // some comment

				public class Main {

				}
			`,
			expected: []string{
				"de.sandstorm.greetings",
				"de.sandstorm.http",
			},
		},
		// TODO: test with different parameters
		// TODO: implement and test comments in general
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			file := bytes.NewBufferString(testCase.fileContent)
			actual, _ := ParseJavaImports(
				testCase.rootPackage,
				testCase.segmentLimit,
				file)
			AssertEquals(t,
				testCase.expected,
				actual)
		})
	}
}
