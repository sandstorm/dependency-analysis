package parsing

import (
	"bytes"
	"testing"
)

func TestParseCppSourceUnit(t *testing.T) {
	testCases := []struct {
		name        string
		fileContent string
		expected    []string
	}{
		{
			name: "simple class without imports",
			fileContent: `package de.sandstorm.test;
			
			public class Main {

			};`,
			expected: []string{"de", "sandstorm", "test", "Main"},
		},
		{
			name: "simple class without imports with default visibility",
			fileContent: `package de.sandstorm.test;
			
			class Main {

			};`,
			expected: []string{"de", "sandstorm", "test", "Main"},
		},
		{
			name: "simple class with internal dependencies",
			fileContent: `package de.sandstorm.test;
			
			#include "path/dependency.h" // some comment

			public class Main {

			};`,
			expected: [][]string{
			[]string{"de", "sandstorm", "test", "Main"},
			[]string{"path", "dependency.h"}
			},
		},
		{
			name: "simple class with external dependencies",
			fileContent: `package de.sandstorm.test;
			
			#include <dependency.h> // some comment

			public class Main {

			};`,
			expected: [][]string{
			[]string{"de", "sandstorm", "test", "Main"},
			[]string{"dependency.h"}
			},
		},
		{
			name: "simple class with external dependencies",
			fileContent: `package de.sandstorm.test;
			
			#include <dependency> // links to a file which contains the actual dependencies

			public class Main {

			};`,
			expected: [][]string{
			[]string{"de", "sandstorm", "test", "Main"},
			[]string{"dependency"}
			},
		}
		// TODO: test private static final class
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

func TestParseJavaImports(t *testing.T) {
	testCases := []struct {
		name        string
		fileContent string
		expected    [][]string
	}{
		{
			name: "simple class with mixed imports",
			fileContent: `
				package de.sandstorm.test;

				#include <QList>;
				#include <path/QList2.h>;
				#include "myQList.h";
				#include "path/myQList.h";

				public class Main {

				}
			`,
			expected: [][]string{
				[]string{"de", "sandstorm", "test", "Main"},
				[]string{"QList"},
				[]string{"path", "QList2.h"},
				[]string{"myList.h"},
				[]string{"path", "myList.h"}
			},
		},
		// TODO: test with different parameters
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			file := bytes.NewBufferString(testCase.fileContent)
			actual, _ := ParseJavaImports(
				file)
			AssertEquals(t,
				testCase.expected,
				actual)
		})
	}
}
