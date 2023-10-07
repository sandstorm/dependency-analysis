package dataStructures

import (
	"sort"
)

type StringSet struct {
	content map[string]bool
}

func NewStringSet() *StringSet {
	result := &StringSet{}
	result.content = make(map[string]bool)
	return result
}

func (this *StringSet) Add(value string) {
	this.content[value] = true
}

func (this *StringSet) AddSet(other *StringSet) {
	for k := range other.content {
		this.Add(k)
	}
}

func (this *StringSet) Contains(value string) bool {
	return this.content[value]
}

func (this *StringSet) Remove(value string) {
	delete(this.content, value)
}

func (this *StringSet) ToArray() []string {
	result := make([]string, len(this.content))
	var i = 0
	for key := range this.content {
		result[i] = key
		i++
	}
	return result
}

func (this *StringSet) String() string {
	result := "[ "
	keys := make([]string, 0, len(this.content))
	for key := range this.content {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		result += key + " "
	}
	return result + "]"
}
