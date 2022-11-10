package golang118

import (
	"errors"
	"fmt"
	"testing"
)

func Test_Definition(t *testing.T) {
	//not works - this is not allowed for basic types
	//type Number[T int|float32|float64] T

	//works for complex types, e.g., slice, map ...
	type Slice1[T int | float32 | float64 | string] []T
	intSlice1 := Slice1[int]{1, 2, 3}
	stringSlice1 := Slice1[string]{"a", "b", "c"}
	fmt.Println(intSlice1)
	fmt.Println(stringSlice1)

	// ~ symbol will check the underlying type
	type MyInt int
	//not works - type(MyInt) != type(int)
	//intSlice1 = Slice1[MyInt]{}
	type Slice2[X ~int] []X
	intSlice2 := Slice2[MyInt]{1, MyInt(2)}
	fmt.Println(intSlice2)
	//not works - ~ symbol does not work for interface
	//type Slice3[Y ~error] []Y
}

func Test_NestDefinition(t *testing.T) {
	type Slice[T int | float32 | float64 | string] []T

	type IntAndStringSlice1[T int | string] Slice[T]
	intAndStrings1 := IntAndStringSlice1[string]{"abc"}
	fmt.Println(intAndStrings1)

	//not works
	//type IntAndStringSlice2[X T | string] Slice[X]

	type Map1[X uint | string, T int | string] map[X]Slice[T]
	//not works
	//m := Map1[uint, float32]{}

	m := Map1[uint, int]{}
	m[1] = Slice[int]{123}
	fmt.Println(m[1])

	type Words[Y int | string, T float32 | string] struct {
		Count    Y
		Contents Slice[T]
	}
	words1 := Words[int, string]{2, Slice[string]{"abc", "xyz"}}
	fmt.Println(words1)
	words2 := Words[int, float32]{1, Slice[float32]{float32(1.0)}}
	fmt.Println(words2)
}

func Test_AnyDefinition(t *testing.T) {
	type Slice[T any] []T
	intSlice := Slice[int]{1, 2, 3}
	fmt.Println(intSlice)
	errorSlice := Slice[error]{errors.New("abc")}
	fmt.Println(errorSlice)
}

func Min[T int | int32 | int64 | float32 | float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Test_FunctionMethods(t *testing.T) {
	fmt.Println(Min[int](10, 20))
	fmt.Println(Min(10, 20)) //auto type inference
	fmt.Println(Min[float64](10.0, 20.0))

	// not works - anonymous function
	//fnAdd := func[T int | float32](a, b T) T {
	//	return a + b
	//}

	type A struct {
	}
	// not works - method
	//func (receiver A) Add[T int | float32 | float64](a T, b T) T {
	//return a + b
	//}
}

type MySlice[T int | float32] []T

func (s MySlice[T]) Sum() T {
	var sum T
	for _, value := range s {
		sum += value
	}
	return sum
}

func Test_Receiver(t *testing.T) {
	var myS1 MySlice[int] = []int{1, 2, 3, 4}
	fmt.Println(myS1.Sum())

	var myS2 MySlice[float32] = []float32{1.0, 2.0, 3.0, 4.0}
	fmt.Println(myS2.Sum())
}

// Interface before 1.18: An interface type specifies a method set called its interface.
// Interface after 1.18: An interface type defines a type set.
func Test_Interface(t *testing.T) {
	type Float interface {
		~float32 | ~float64 //union
	}
	type Int interface {
		~int32 | ~int64
	}
	// not works - cannot be used to create variable
	//var a Float
	//fmt.Println(a)

	// union
	type FloatOrInt interface {
		Float | Int
	}

	// Intersect
	type FloatAndInt interface {
		Float
		Int
	}

	var num any // any = interface{}
	num = 10
	fmt.Println(num)
}

// first kind of interface
func Test_BasicInterface(t *testing.T) {
	// only methods in interface
	type MyError interface {
		Error() string
	}
	fmt.Errorf("hello world")
}

// second kind of interface
func Test_GeneralInterface(t *testing.T) {
	// types in interface
	type Uint interface {
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
	}

	type ReadWriter interface {
		~string | ~[]rune

		Read(p []byte) (n int, err error)
		Write(p []byte) (n int, err error)
	}

}

type DataProcessor[T any] interface {
	Process(oriData T) (newData T)
	Save(data T) error
}

type CSVProcessor struct {
}

func (c CSVProcessor) Process(oriData string) (newData string) {
	fmt.Println("process", oriData)
	newData = oriData
	return
}

func (c CSVProcessor) Save(oriData string) error {
	fmt.Println("save", oriData)
	return nil
}

func Test_GenericInterface1(t *testing.T) {
	// DataProcessor[string] is a basic interface
	var processor1 DataProcessor[string]
	processor1 = CSVProcessor{}
	processor1.Process("this is raw data")

	var processor2 DataProcessor[string] = CSVProcessor{}
	processor2.Process("this is raw data")
	// not works
	//var processor3 DataProcessor[int] = CSVProcessor{}
}

type DataProcessor2[T any] interface {
	string | []byte | ~struct{ Data interface{} }

	Process(data T) (newData T)
	Save(data T) error
}

type NumberProcessor int

func (c NumberProcessor) Process(oriData string) (newData string) {
	fmt.Println("process", oriData)
	newData = oriData
	return
}

func (c NumberProcessor) Save(oriData string) error {
	fmt.Println("save", oriData)
	return nil
}

type JsonProcessor struct {
	Data interface{}
}

func (c JsonProcessor) Process(oriData string) (newData string) {
	fmt.Println("process", oriData)
	newData = oriData
	return
}

func (c JsonProcessor) Save(oriData string) error {
	fmt.Println("save", oriData)
	return nil
}

func Test_GenericInterface2(t *testing.T) {
	// not works - DataProcessor2[string] is a general interface
	//var processor1 DataProcessor2[string]

	type ProcessorList[T DataProcessor2[string]] []T
	pls := ProcessorList[JsonProcessor]{JsonProcessor{}}
	for _, pl := range pls {
		fmt.Println(pl.Process("abc"))
	}

	// not works - underlying data of NumberProcessor is int
	// pls = ProcessorList[NumberProcessor]{NumberProcessor{}}
}
