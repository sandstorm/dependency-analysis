package dataStructures

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
	for key, _ := range this.content {
		result[i] = key
		i++
	}
	return result
}
