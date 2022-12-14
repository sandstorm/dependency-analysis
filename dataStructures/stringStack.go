package dataStructures

// array of dynamic size but with fixed capacity
type StringStack struct {
	content []string
	length  int
}

func NewStringStack(capacity int) *StringStack {
	result := &StringStack{}
	result.content = make([]string, capacity)
	result.length = 0
	return result
}

func (this *StringStack) Push(element string) {
	this.content[this.length] = element
	this.length++
}

func (this *StringStack) Pop(element string) string {
	if this.length > 0 {
		this.length--
		return this.content[this.length]
	} else {
		return ""
	}
}

func (this *StringStack) IndexOf(element string) int {
	for i := 0; i < this.length; i++ {
		if this.content[i] == element {
			return i
		}
	}
	return -1
}

func (this *StringStack) Get(index int) string {
	return this.content[index]
}

func (this *StringStack) Length() int {
	return this.length
}

func (this *StringStack) ToArray() []string {
	return this.content[0:this.length]
}
