package utils

type Index struct {
	Value   string
}

var INDEX_START byte = 48
var INDEX_END byte = 122

func (element *Index) Compare(index Index) int {
	var elements []byte = []byte(ReverseString(element.Value))
	var elements2 []byte = []byte(ReverseString(index.Value))
	for i := 0; i < len(element.Value); i++ {
		if elements[i] > elements2[i] {
			return 1
		} else if elements[i] < elements2[i] {
			return -1
		}
	}
	return 0
}

func (element *Index) New(size int) {
	var elements []byte = []byte{}
	for i := 0; i < size; i++ {
		elements = append(elements, INDEX_START)
	}
	element.Value = string(elements)
}

func (element *Index) Next() Index {
	if element.IsMaxValue() {
		byteArray := []byte{}
		for i := 0; i < len(element.Value); i++ {
			byteArray = append(byteArray, INDEX_END)
		}
		newIndex := Index{Value:string(byteArray)}
		return newIndex
	}
	var elements []byte = []byte(element.Value)
	var length int = len(elements)
	var report bool = true
	for i := length-1; i >= 0; i-- {
		if report {
			if elements[i] == INDEX_END {
				elements[i] = INDEX_START
				report = true
			} else {
				report = false
				elements[i]=elements[i] + 1
			}
		}
	}
	return Index{
		Value: string(elements),
	}
}

func (element *Index) IsZero() bool {
	for i := 0; i < len(element.Value); i++ {
			if element.Value[i] > INDEX_START {
				return false
			}
	}
	return true
}

func (element *Index) IsMaxValue() bool {
	for i := 0; i < len(element.Value); i++ {
		if element.Value[i] < INDEX_END {
			return false
		}
	}
	return true
}

func (element *Index) Previous() Index {
	if element.IsZero() {
		newIndex := Index{Value:""}
		newIndex.New(len(element.Value))
		return newIndex
	}
	var elements []byte = []byte(element.Value)
	var length int = len(element.Value)
	var report bool = true
	for i := length-1; i >= 0; i-- {
		if report {
			if elements[i] == INDEX_START {
				elements[i] = INDEX_END
				report = true
			} else {
				report = false
				elements[i]--
			}
		}
	}
	return Index{
		Value: string(elements),
	}
}

func (element *Index) FromInt(value int, length int) {
	element.New(length)
	var elements []byte = []byte(element.Value)
	var report byte = byte(value)
	var RATIONAL byte = INDEX_END - INDEX_START
	for i := length-1; i >= 0; i-- {
		var newValue byte = (elements[i] + report)
		if newValue >= INDEX_END {
			value := newValue % RATIONAL
			elements[i] = value
			report = (newValue - value) / RATIONAL
		} else {
			elements[i] = newValue
			report = 0
		}
	}
	element.Value = string(elements)
}



func (element *Index) Sum(index Index) {
	var elements []byte = []byte(element.Value)
	var elements2 []byte = []byte(index.Value)
	var length int = len(element.Value)
	var report byte = 0
	var RATIONAL byte = INDEX_END - INDEX_START
	for i := length-1; i >= 0; i-- {
		var newValue byte = (elements[i] + (elements2[i]-INDEX_START) + report)
		if newValue >= INDEX_END {
			value := newValue % RATIONAL
			elements[i] = value
			report = (newValue - value) / RATIONAL
		} else {
			elements[i] = newValue
			report = 0
		}
	}
	element.Value = string(elements)
}

func (element *Index) Subtract(index Index) {
	var elements []byte = []byte(element.Value)
	var elements2 []byte = []byte(index.Value)
	var length int = len(element.Value)
	var report byte = 0
	for i := length-1; i >= 0; i-- {
		var newValue byte = (elements[i] - (elements2[i]-INDEX_START) - report)
		if newValue < INDEX_START {
			var value byte = 0
			value = INDEX_START - newValue
			elements[i] = INDEX_END - value
			report = 1
		} else {
			elements[i] = newValue
			report = 0
		}
	}
	element.Value = string(elements)
}
