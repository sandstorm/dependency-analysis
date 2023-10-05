package parsing

import (
	"bytes"
	"testing"
)

func TestParsePhpSourceUnit(t *testing.T) {
	testCases := []struct {
		name        string
		fileContent string
		expected    fullyQualifiedType
	}{
		{
			name: "simple class without imports",
			fileContent: `<?php

			namespace App\Controller;
			
			public class MyController extends AbstractController
			{

			}`,
			expected: fullyQualifiedType{"App", "Controller", "MyController"},
		},
		{
			name: "simple class with imports",
			fileContent: `<?php

			namespace App\Controller;

			use App\ArrayUtils;
			use App\Entity\MyEntity;
			use Symfony\Component\HttpFoundation\Request;
			use Symfony\Component\HttpFoundation\Response;
			use Symfony\Component\Routing\Annotation\Route;
			
			#[Route('/my')]
			class MyController extends AbstractController
			{

			}`,
			expected: fullyQualifiedType{"App", "Controller", "MyController"},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			file := bytes.NewBufferString(testCase.fileContent)
			AssertEquals(t,
				[]fullyQualifiedType{testCase.expected},
				ParsePhpSourceUnit(file),
			)
		})
	}
}

func TestParsePhpImports(t *testing.T) {
	testCases := []struct {
		name        string
		fileContent string
		expected    []fullyQualifiedType
	}{
		{
			name: "simple class without imports",
			fileContent: `<?php

			namespace App\Controller;
			
			public class MyController extends AbstractController
			{

			}`,
			expected: []fullyQualifiedType{},
		},
		{
			name: "simple class with imports",
			fileContent: `<?php

			use App\ArrayUtils;
			use App\Entity\MyEntity   ;
			use Symfony\Component\HttpFoundation\Request;
			use    Symfony\Component\HttpFoundation\Response;
			use Symfony\Component\Routing\Annotation\Route;
			
			#[Route('/my')]
			public class MyController extends AbstractController
			{

			}`,
			expected: []fullyQualifiedType{
				{"App", "ArrayUtils"},
				{"App", "Entity", "MyEntity"},
				{"Symfony", "Component", "HttpFoundation", "Request"},
				{"Symfony", "Component", "HttpFoundation", "Response"},
				{"Symfony", "Component", "Routing", "Annotation", "Route"},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			file := bytes.NewBufferString(testCase.fileContent)
			actual, _ := ParsePhpImports(
				file)
			AssertEquals(t,
				testCase.expected,
				actual)
		})
	}
}
