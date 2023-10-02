package parsing

import (
	"bytes"
	"testing"
)

func TestParseJavaScriptSourceUnit(t *testing.T) {
	testCases := []struct {
		name       string
		sourcePath string
		expected   [][]string
	}{
		{
			name:       "file path without dots",
			sourcePath: "src/Components/Button/button.js",
			expected:   [][]string{[]string{"src", "Components", "Button", "button"}},
		},
		{
			name:       "file path with dots",
			sourcePath: "a/.././src/a/../b/c/../.././Components/Button/button.js",
			expected:   [][]string{[]string{"src", "Components", "Button", "button"}},
		},
		{
			name:       "file path with index.js",
			sourcePath: "a/.././src/a/../b/c/../.././Components/Button/index.js",
			expected:   [][]string{[]string{"src", "Components", "Button"}},
		},
		{
			name:       "ignore node_modules",
			sourcePath: "node_modules/some/lib.js",
			expected:   [][]string{},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			AssertEquals(t,
				testCase.expected,
				ParseJavaScriptSourceUnit(testCase.sourcePath),
			)
		})
	}
}

func TestParseJavaScriptImports(t *testing.T) {
	sourceRoot := "my/projects/src/root/path"
	testCases := []struct {
		name        string
		fileContent string
		expected    [][]string
	}{
		{
			name: "no imports",
			fileContent: `package dataStructures

			function myCoolFunction() {}`,
			expected: [][]string{},
		},
		{
			name: "library imports",
			fileContent: `package analysis

			import React from 'react';
			
			function myCoolFunction() {}`,
			expected: [][]string{
				[]string{"react"},
			},
		},
		{
			name: "local import",
			fileContent: `package analysis

			import User from './Model/user';

			function myCoolFunction() {}`,
			expected: [][]string{
				[]string{"my", "projects", "src", "root", "path", "Model", "user"},
			},
		},
		{
			name: "local and library imports",
			fileContent: `package analysis

			import React from 'react';
			import User from './Model/user';

			function myCoolFunction() {}`,
			expected: [][]string{
				[]string{"react"},
				[]string{"my", "projects", "src", "root", "path", "Model", "user"},
			},
		},
		{
			name: "import of folder (simple import since source-units do not have index suffix)",
			fileContent: `package analysis

			import Model from './Model';
			
			function myCoolFunction() {}`,
			expected: [][]string{
				[]string{"my", "projects", "src", "root", "path", "Model"},
			},
		},
		{
			name: "star import",
			fileContent: `package analysis

			import * from './Model';
			
			function myCoolFunction() {}`,
			expected: [][]string{
				[]string{"my", "projects", "src", "root", "path", "Model"},
			},
		},
		{
			name: "named import",
			fileContent: `package analysis

			import User from './Model';
			
			function myCoolFunction() {}`,
			expected: [][]string{
				[]string{"my", "projects", "src", "root", "path", "Model"},
			},
		},
		{
			name: "member import",
			fileContent: `package analysis

			import { User } from './Model';
			
			function myCoolFunction() {}`,
			expected: [][]string{
				[]string{"my", "projects", "src", "root", "path", "Model"},
			},
		},
		{
			name: "named and member mixed import",
			fileContent: `package analysis

			import Model, { User } from './Model';
			
			function myCoolFunction() {}`,
			expected: [][]string{
				[]string{"my", "projects", "src", "root", "path", "Model"},
			},
		},
		{
			name: "alias import",
			fileContent: `package analysis

			import { User as Nutzer } from './Model';
			
			function myCoolFunction() {}`,
			expected: [][]string{
				[]string{"my", "projects", "src", "root", "path", "Model"},
			},
		},
		{
			name: "from-only import",
			fileContent: `package analysis

			from "react";
			
			function myCoolFunction() {}`,
			expected: [][]string{
				[]string{"react"},
			},
		},
		{
			name: "optional semicolon",
			fileContent: `package analysis

			from "react"
			
			function myCoolFunction() {}`,
			expected: [][]string{
				[]string{"react"},
			},
		},
		{
			name: "single ticks semicolon",
			fileContent: `package analysis

			from 'react';
			
			function myCoolFunction() {}`,
			expected: [][]string{
				[]string{"react"},
			},
		},
		{
			name: "complex path",
			fileContent: `package analysis

			import User from '../../Model/x/./a/./../..';
			
			function myCoolFunction() {}`,
			expected: [][]string{
				[]string{"my", "projects", "src", "Model"},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			file := bytes.NewBufferString(testCase.fileContent)
			actual, _ := ParseJavaScriptImports(
				sourceRoot+"/myFile.js",
				file)
			AssertEquals(t,
				testCase.expected,
				actual)
		})
	}
}
