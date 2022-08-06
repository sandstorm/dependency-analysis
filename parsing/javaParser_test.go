package parsing

import (
	"reflect"
	"testing"
	"sort"
)

func AssertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if actual == expected {
		return
	}
	t.Errorf("Received %v (type %v), expected %v (type %v)", 
		actual, reflect.TypeOf(actual),
		expected, reflect.TypeOf(expected))
}

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
		AssertEqual(t, expected, actual)
	})


}

