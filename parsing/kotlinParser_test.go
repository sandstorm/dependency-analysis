package parsing

import (
	"bytes"
	"testing"
)

func TestParseKotlinSourceUnit(t *testing.T) {
	testCases := []struct {
		name        string
		fileContent string
		expected    []fullyQualifiedType
	}{
		{
			name: "simple class without imports",
			fileContent: `package de.sandstorm.test;
			
			class TheApp

			fun main(args: Array<String>) {
				SpringApplication.run(TheApp::class.java, *args)
			}`,
			expected: []fullyQualifiedType{{"de", "sandstorm", "test", "TheApp"}},
		},
		{
			name: "simple class without imports with default visibility",
			fileContent: `package de.sandstorm.test;
			
			class Main {

			}`,
			expected: []fullyQualifiedType{{"de", "sandstorm", "test", "Main"}},
		},
		{
			name: "simple class with imports",
			fileContent: `package de.sandstorm.test
			
			import org.springframework.boot.SpringApplication
			import org.springframework.boot.autoconfigure.SpringBootApplication
			import org.springframework.boot.context.properties.ConfigurationPropertiesScan
			import org.springframework.scheduling.annotation.EnableScheduling

			@EnableScheduling
			@SpringBootApplication
			@ConfigurationPropertiesScan
			class TheApp

			fun main(args: Array<String>) {
				SpringApplication.run(TheApp::class.java, *args)
			}`,
			expected: []fullyQualifiedType{{"de", "sandstorm", "test", "TheApp"}},
		},
		{
			name: "simple class with imports with ';'",
			fileContent: `package de.sandstorm.test;
			
			import org.springframework.boot.SpringApplication;
			import org.springframework.boot.autoconfigure.SpringBootApplication;
			import org.springframework.boot.context.properties.ConfigurationPropertiesScan;
			import org.springframework.scheduling.annotation.EnableScheduling;

			@EnableScheduling
			@SpringBootApplication
			@ConfigurationPropertiesScan
			class TheApp

			fun main(args: Array<String>) {
				SpringApplication.run(TheApp::class.java, *args)
			}`,
			expected: []fullyQualifiedType{{"de", "sandstorm", "test", "TheApp"}},
		},
		{
			name: "multiple classes and data classes",
			fileContent: `package de.sandstorm.test;

			import org.springframework.stereotype.Service
			import java.math.BigInteger
			import java.security.MessageDigest
			
			data class Person(val name: String)

			data class Item(val price: Int)

			@Service
			class MyService(
				private val settings: MySettings
			) {
				fun do(name: String): String {
				}
			} 
			`,
			expected: []fullyQualifiedType{
                {"de", "sandstorm", "test", "Person"},
                {"de", "sandstorm", "test", "Item"},
                {"de", "sandstorm", "test", "MyService"},
			},
		},
		{
			name: "open class",
			fileContent: `package de.sandstorm.test;
			
				open class Person(val name: String)`,
			expected: []fullyQualifiedType{
				{"de", "sandstorm", "test", "Person"},
			},
		},
		{
			name:       "extension functions only",
			fileContent: `package de.sandstorm.test;

			import java.net.URLEncoder
			
			fun String.trimLines() = this.lines().map { it.trim() }.joinTo(StringBuilder(this.length), "").toString()
			fun String.urlEncode() = URLEncoder.encode(this, "UTF-8")!!`,
			expected: []fullyQualifiedType{{"de", "sandstorm", "test"}},
		},
		{
			name: "constants only",
			fileContent: `package de.sandstorm.test

			const val ADDON_KEY = "io.exply.jira.plugin"
			val EJC_INSTANCE_NAME= "app"`,
			expected: []fullyQualifiedType{{"de", "sandstorm", "test"}},
		},
		// TODO: interfaces, enums, const, fun
		// TODO: class visibilities
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			content := bytes.NewBufferString(testCase.fileContent)
			AssertEquals(t,
				testCase.expected,
				ParseKotlinSourceUnit(content),
			)
		})
	}
}

func TestParseKotlinImports(t *testing.T) {
	testCases := []struct {
		name        string
		fileContent string
		expected    []fullyQualifiedType
	}{
		{
			name: "simple class with mixed imports",
			fileContent: `
				package de.sandstorm.test;

				import java.util.List
				import        de.sandstorm.greetings.HelloWorld
				import        de.sandstorm.http.requests.GetGreetingRequest
				import    org.junit.Assert.assertEquals; // some comment

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
			actual, _ := ParseKotlinImports(
				file)
			AssertEquals(t,
				testCase.expected,
				actual)
		})
	}
}
