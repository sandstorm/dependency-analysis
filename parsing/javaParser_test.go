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

func TestParseJavaPackage(t *testing.T) {
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
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			file := bytes.NewBufferString(testCase.fileContent)
			AssertEquals(t,
				testCase.expected,
				ParseJavaPackage(file),
			)
		})
	}
}

/*
func TestParseJavaFile(t *testing.T) {
	t.Run("mixed imports", func(t *testing.T) {
		actual := ParseJavaFile(`
			package de.sandstorm.test;

			import java.util.List;
			import        de.sandstorm.greetings.HelloWorld;
			import        de.sandstorm.http.requests.GetGreetingRequest;
			import  static  org.junit.Assert.assertEquals; // some comment

			public class Main {

			}
		`)
		lvl1 := []string{"java", "de", "org"}
		lvl2 := []string{"java.util", "de.sandstorm", "org.junit"}
		lvl3 := []string{"java.util.List", "de.sandstorm.greetings", "de.sandstorm.http", "org.junit.Assert"}
		sort.Strings(lvl1)
		sort.Strings(lvl2)
		sort.Strings(lvl3)
		expected := [][]string{lvl1, lvl2, lvl3}
		AssertEquals(t, expected, actual)
	})


}
*/

