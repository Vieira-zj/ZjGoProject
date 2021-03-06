package demos

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// demo, interface
type abser interface {
	abs() float64
	string() string
}

type myFloat float64

func (f myFloat) abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

func (f *myFloat) string() string {
	return fmt.Sprintf("%.2f", f.abs())
}

func testInterface01() {
	var a abser
	f := myFloat(-math.Sqrt2)
	fmt.Println("\nby value:")
	fmt.Println("abs value:", f.abs())
	fmt.Println("float string:", f.string())

	fmt.Println("\nby reference:")
	a = &f
	fmt.Println("abs value:", a.abs())
	fmt.Println("float string:", a.string())
}

// demo, interface
type iMyGetter interface {
	myGet() string
}

type iMySetter interface {
	mySet(string)
}

type iMyGetterAndSetter interface {
	iMyGetter
	iMySetter
}

type zjGetterAndSetter struct {
	name string
	age  int
	desc string
}

func (instance *zjGetterAndSetter) mySet(val string) {
	instance.desc = val
}

func (instance zjGetterAndSetter) myGet() string {
	if len(instance.desc) == 0 {
		instance.desc = "empty"
	}
	return fmt.Sprintf("name: %s, age: %d, desc: %s", instance.name, instance.age, instance.desc)
}

func testInterface02() {
	testStruct := zjGetterAndSetter{
		name: "vieira",
		age:  30,
	}

	var testRef iMyGetterAndSetter = &testStruct
	fmt.Println(testRef.myGet())

	initAndPrintInfoByStruct(testStruct)
	initAndPrintInfoByPointer(&testStruct)
	initAndPrintInfoByInterface(&testStruct)
}

func initAndPrintInfoByStruct(arg zjGetterAndSetter) {
	arg.mySet("this is a struct.")
	fmt.Println(arg.myGet())
}

func initAndPrintInfoByPointer(arg *zjGetterAndSetter) {
	arg.mySet("this is a pointer.")
	fmt.Println(arg.myGet())
}

func initAndPrintInfoByInterface(arg iMyGetterAndSetter) {
	arg.mySet("this is an interface.")
	fmt.Println(arg.myGet())
}

// demo, OO inherit
type super struct {
	Name string
}

func (s super) Print() {
	fmt.Println("name:", s.Name)
}

type sub struct {
	super
	Desc string
}

func (s sub) PrintDesc() {
	fmt.Println("desc:", s.Desc)
}

func testOOInherit() {
	s := super{Name: "parent_1"}
	sub1 := sub{super: s, Desc: "child_1 from parent_1."}
	sub1.Print()
	sub1.PrintDesc()

	sub2 := new(sub)
	sub2.Name = "sub_2"
	sub2.Desc = "this is child_2."
	sub2.Print()
	sub2.PrintDesc()
}

// demo, rw mutex
func testRwMutex() {
	mutex := new(sync.RWMutex)
	fmt.Println("\nready in main")
	mutex.Lock()
	fmt.Println("mutex locked in main")

	chs := make([]chan int, 4)
	for i := 0; i < 4; i++ {
		chs[i] = make(chan int)
		go func(i int, ch chan<- int) {
			fmt.Println("ready in routine:", i)
			mutex.RLock()
			fmt.Println("mutex read locked in routine:", i)
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			fmt.Println("mutex read unlocked in routine:", i)
			mutex.RUnlock()
			ch <- i
		}(i, chs[i])
	}

	time.Sleep(time.Second)
	fmt.Println("mutex unlocked in main")
	mutex.Unlock()

	for _, ch := range chs {
		<-ch
	}
}

// demo, panic and recover()
func testPanicRecover() {
	defer func() {
		fmt.Println("\nrecover:")
		if r := recover(); r != nil {
			myLog(r)
		}
	}()
	myWork(true)
}

func myWork(isOccur bool) {
	myLog("start", "myWork")
	if isOccur {
		panic("mock error")
	}
	myLog("end", "myWork")
}

func myLog(args ...interface{}) {
	fmt.Printf("args type: %T\n", args)
	fmt.Println(args...)
}

// demo, Error
type myError struct {
	infoa string
	infob string
	err   error
}

func (merr *myError) Error() string {
	errInfo := fmt.Sprintf(
		"infoa: %s, infob: %s, original error info: %s", merr.infoa, merr.infob, merr.err.Error())
	return errInfo
}

func testCreateError() {
	// #1
	var err = errors.New("new mock error")
	fmt.Println("\nerror:", err.Error())

	// #2
	err = fmt.Errorf("%s", "error from fmt.Errorf()")
	fmt.Println("error:", err)

	// #3
	myErr := &myError{
		infoa: "error info a",
		infob: "error info b",
		err:   errors.New("test custom error"),
	}
	fmt.Println("custom error:", myErr)
}

// demo, []string in array and map
func testArrayInMap() {
	fmt.Println("\n#1: by map:")
	m := make(map[int][]string, 3)
	m[1] = []string{"a1", "a2", "a3"}
	m[2] = []string{"b1", "b2", "b3"}
	fmt.Printf("map length: %d\n", len(m))
	for k, v := range m {
		fmt.Printf("%d=%v\n", k, strings.Join(v, ","))
	}

	fmt.Println("\n#2: by slice:")
	arr := [...][3]string{
		{"a1", "a2", "a3"},
		{"b1", "b2", "b3"},
		{"c1", "c2", "c3"},
	}
	for idx, item := range arr {
		fmt.Printf("%d=%v\n", idx, strings.Join(item[:], ":"))
	}

	fmt.Println("\n#3: by slice:")
	var s [][]string
	for i := 0; i < 3; i++ {
		var sTmp []string
		for j := 0; j < 3; j++ {
			sTmp = append(sTmp, strconv.Itoa(i+j))
		}
		s = append(s, sTmp)
	}
	for idx, item := range s {
		fmt.Printf("%d=%v\n", idx, strings.Join(item, "|"))
	}
}

// demo, slice is sequence
func testSliceOrder() {
	s := make([]string, 0, 10)
	s = append(s, "one")
	s = append(s, "two")
	s = append(s, "three")
	s = append(s, "3")
	s = append(s, "2")
	s = append(s, "4")
	s = append(s, "1")

	fmt.Println("\nslice values:")
	for i, v := range s {
		fmt.Printf("%d=%s\n", i, v)
	}
}

// demo, map is not sequence
func testMapOrder() {
	m := make(map[int]string)
	m[1] = "one"
	m[5] = "five"
	m[2] = "two"
	m[4] = "four"
	m[3] = "three"

	fmt.Println("\nmap values:")
	for k, v := range m {
		fmt.Printf("%d=%s\n", k, v)
	}
}

// demo, get map value
func testGetMapValue() {
	s := make([]int, 5, 10)
	fmt.Printf("\ninit slice: %v, length: %d, cap: %d\n", s, len(s), cap(s))

	m := make(map[string]string, 5)
	m["1"] = "one"
	m["2"] = "two"
	m["3"] = "three"
	fmt.Printf("init map: %v, length: %d\n", m, len(m))

	key := "2"
	if val, ok := m[key]; ok {
		fmt.Printf("key[%s]: value[%s]\n", key, val)
	} else {
		fmt.Printf("key[%s]: value not found!\n", key)
	}
}

// demo, init reference variable
func testInitRefVar() {
	var i *int   // 声明变量
	i = new(int) // 分配内存空间
	fmt.Printf("\ndefault int value: %d\n", *i)
	*i = 10
	fmt.Printf("int value: %d\n", *i)
}

// demo, init slice and map by make()
func testInitVarByMake() {
	s1 := make([]int, 10, 10)
	fmt.Printf("\nmake slice: %v, len: %d, cap: %d\n", s1, len(s1), cap(s1))
	s2 := []int{}
	fmt.Printf("init slice: %v, len: %d, cap: %d\n", s2, len(s2), cap(s2))

	m1 := make(map[int]string, 5)
	fmt.Printf("make map: %+v, len: %d\n", m1, len(m1))
	m2 := map[int]string{}
	fmt.Printf("init map: %+v, len: %d\n", m2, len(m2))
}

// demo, printf format
func testPrintfFormat() {
	type point struct {
		x, y int
	}
	p := point{1, 2}

	fmt.Println("\nprintf format:")
	fmt.Printf("string: %s\n", "\"string\"")
	fmt.Printf("raw string (go): %q\n", "\"string\"")

	fmt.Printf("p addr: %p\n", &p)
	fmt.Printf("p type: %T\n", p)
	fmt.Printf("p struct: %v\n", p)
	fmt.Printf("p struct with fields: %+v\n", p)
	// print struct in golang code format
	fmt.Printf("p struct (go): %#v\n", p)

	fmt.Printf("binary: %b\n", 14)
	fmt.Printf("char: %c\n", 33)
	fmt.Printf("hex: %x\n", 456)
	fmt.Printf("float: %.1f\n", 78.91)

	fmt.Printf("e+: %e\n", 123400000.0)
	fmt.Printf("E+: %E\n", 123400000.0)
}

// demo, regexp
func testRegExpMore() {
	fmt.Println("\n#1: First match")
	re := regexp.MustCompile(`foo.?`)
	fmt.Printf("%q\n", re.FindString("seafood fool"))
	fmt.Printf("%q\n", re.FindString("meat"))

	fmt.Println("\n#2: Location")
	re = regexp.MustCompile(`ab?`)
	fmt.Println(re.FindStringIndex("tablett"))
	fmt.Println(re.FindStringIndex("foo") == nil)

	fmt.Println("\n#3: All matches")
	// if n >= 0, the function returns at most n matches
	re = regexp.MustCompile(`a.`)
	fmt.Printf("%v\n", re.FindAllString("paranormal", -1))
	fmt.Printf("%v\n", re.FindAllString("paranormal", 2))
	fmt.Printf("%v\n", re.FindAllString("graal", -1))
	fmt.Printf("%v\n", re.FindAllString("none", -1))

	fmt.Println("\n#4: Replace")
	re = regexp.MustCompile(`ab*`)
	fmt.Printf("%q\n", re.ReplaceAllString("-a-abb-", "T"))

	fmt.Println("\n#5: Split")
	a := regexp.MustCompile(`a`)
	fmt.Printf("%q\n", a.Split("banana", -1))
	fmt.Printf("%q\n", a.Split("banana", 0))
	fmt.Printf("%q\n", a.Split("banana", 1))
	fmt.Printf("%q\n", a.Split("banana", 2))

	zp := regexp.MustCompile(`z+`)
	fmt.Printf("%q\n", zp.Split("pizza", -1))
}

// demo, rune type
func testRuneType() {
	s := "测试 hello world"
	fmt.Println("\nstring length:", len(s))
	b := []byte(s)
	fmt.Println("bytes length:", len(b))
	r := []rune(s)
	fmt.Println("runes count:", len(r))

	fmt.Println("top 2 cn words:", string(r[:2]))
}

// demo, base64 encode and decode
func testBase64Code() {
	const str = "Go 言语编程 "
	base64EncodeAndDecode(base64.StdEncoding, str)
	base64EncodeAndDecode(base64.URLEncoding, str)
	base64EncodeAndDecode(base64.RawStdEncoding, str)
	base64EncodeAndDecode(base64.RawURLEncoding, str)
}

func base64EncodeAndDecode(enc *base64.Encoding, input string) {
	encStr := enc.EncodeToString([]byte(input))
	fmt.Printf("\nbase64 encoded string: %s\n", encStr)

	decStr, err := enc.DecodeString(encStr)
	if err != nil {
		panic("base64 decode error!")
	}
	fmt.Printf("base64 decoded string: %s\n", decStr)

	if input != string(decStr) {
		panic(errors.New("not equal"))
	}
}

// demo, io writer
func testIOWriter() {
	// #1
	// bytes.Buffer实现了io.Writer和io.Reader接口
	var b1 bytes.Buffer
	n, err := b1.Write([]byte("hello world\n"))
	if err != nil {
		panic(err)
	}
	fmt.Println("\nwrite bytes size:", n)
	// (b *Buffer)实现了io.Writer接口, 传入&b
	fmt.Fprint(&b1, "this is a test.\n")
	b1.WriteTo(os.Stdout)

	// #2
	b2 := bytes.NewBuffer([]byte("hello world\n"))
	fmt.Fprint(b2, "this is a test2.\n")
	fmt.Println("\nwrite bytes:")
	b2.WriteTo(os.Stdout)
}

// demo, io reader
func testIOReader() {
	// #1
	var b1 bytes.Buffer
	var data1 []byte
	b1.Write([]byte("hello world\n"))
	fmt.Fprint(&b1, "this is a test.\n")
	for {
		var p [8]byte
		n, err := b1.Read(p[:])
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		fmt.Printf("read size=%d, bytes: %v\n", n, string(p[:]))
		data1 = append(data1, p[:]...)
	}
	fmt.Printf("read bytes data:\n%s", string(data1))

	// #2
	b2 := bytes.NewBuffer([]byte("hello world\n"))
	fmt.Fprint(b2, "this is a test2.\n")
	data2, err := ioutil.ReadAll(b2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nread bytes data:\n%s", string(data2))
}

// demo, SetFinalizer for object
type testPeople struct {
	id  string
	age int
	job string
}

func (p testPeople) String() string {
	return fmt.Sprintf("id=%s,age=%d,job=%s", p.id, p.age, p.job)
}

func createPeople() {
	fmt.Println("\ninit and SetFinalizer for people")
	p := &testPeople{
		id:  "001",
		age: 20,
		job: "tester",
	}
	runtime.SetFinalizer(p, clearUp)
	fmt.Println("new people:", *p)
}

func clearUp(p *testPeople) {
	fmt.Printf("clearup people: %+v\n", p)
}

func testFinalizerInObject() {
	createPeople()
	// 局部变量p已经不可达, 被GC
	for i := 0; i < 3; i++ {
		fmt.Println("runtime GC...")
		runtime.GC()
		time.Sleep(time.Second)
	}
}

// MainDemo02 main for golang demo02.
func MainDemo02() {
	// testInterface01()
	// testInterface02()
	// testOOInherit()
	// testRwMutex()

	// testPanicRecover()
	// testCreateError()

	// testArrayInMap()
	// testSliceOrder()
	// testMapOrder()
	// testGetMapValue()

	// testInitRefVar()
	// testInitVarByMake()

	// testPrintfFormat()
	// testRegExpMore()
	// testRuneType()
	// testBase64Code()

	// testIOWriter()
	// testIOReader()

	// testFinalizerInObject()

	fmt.Println("golang demo02 DONE.")
}
