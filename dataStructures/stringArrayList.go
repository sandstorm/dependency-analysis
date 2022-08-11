package dataStructures

// array of dynamic size but with fixed capacity
type StringArrayList struct {
	content []string
	length int
}

func NewStringArrayList(capacity int) *StringArrayList {
	result := &StringArrayList{}
	result.content = make([]string, capacity)
	result.length = 0
    return result
}

func (this *StringArrayList) Push(element string) {
	this.content[this.length] = element
	this.length++
}

func (this *StringArrayList) Pop(element string) string {
	if this.length > 0 {
		this.length--
		return this.content[this.length]
	} else {
		return ""
	}
}

func (this *StringArrayList) IndexOf(element string) int {
	for i := 0; i < this.length; i++ {
		if this.content[i] == element {
			return i
		}
	}
	return -1
}

func (this *StringArrayList) Get(index int) string {
	return this.content[index]
}

func (this *StringArrayList) Length() int {
	return this.length
}

func (this *StringArrayList) ToArray() []string {
	return this.content[0:this.length]
}
