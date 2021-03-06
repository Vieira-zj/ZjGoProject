package demos

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/jmcvetta/randutil"
	"gopkg.in/mgo.v2/bson"
)

// demo, var "HelloMessage" init before init() function
func init() {
	fmt.Println("[demo04.go] init") // #2
}

// HelloMsg public var, invoked from main.go.
var HelloMsg = sayHello()

func sayHello() string {
	fmt.Println("[demo04.go] start run sayHello()") // #1
	return "hello world!"
}

// demo, get file base and full name
func testGetFileName() {
	srcPath := os.Getenv("HOME") + "/Downloads/tmp_files"
	f, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("\nfile full name:", f.Name())

	info, err := f.Stat()
	if err != nil {
		panic(err)
	}
	fmt.Println("file base name:", info.Name())
}

// demo, verify go version
func testVerifyGoVersion() {
	curVersion := runtime.Version()
	fmt.Printf("\n%s >= go1.15: %v\n", curVersion, isGoVersionOK("1.15"))
	fmt.Printf("%s >= go1.10: %v\n", curVersion, isGoVersionOK("1.10"))
	fmt.Printf("%s >= go1.9.3: %v\n", curVersion, isGoVersionOK("1.9.3"))
}

func isGoVersionOK(baseVersion string) bool {
	curVersion := runtime.Version()[2:]
	curArr := strings.Split(curVersion, ".")
	baseArr := strings.Split(baseVersion, ".")

	for i := 0; i < 2; i++ { // check first 2 digits
		cur, _ := strconv.ParseInt(curArr[i], 10, 32)
		base, _ := strconv.ParseInt(baseArr[i], 10, 32)
		if cur == base {
			continue
		}
		return cur > base
	}
	return true // cur == base
}

// demo, time calculation
func testTimeOpSub() {
	start := time.Now()
	time.Sleep(time.Duration(2) * time.Second)
	duration := time.Now().Sub(start)
	fmt.Printf("time duration: %.2f\n", duration.Seconds())

	for int(time.Now().Sub(start).Seconds()) < 5 {
		fmt.Println("wait 1 second ...")
		time.Sleep(time.Second)
	}
}

// demo, test get random strings
func testRandomValues() {
	if num, err := randutil.IntRange(1, 10); err == nil {
		fmt.Println("get a random number 1-10:", num)
	}

	if str1, err := randutil.String(10, randutil.Numerals); err == nil {
		fmt.Println("get string of 10 chars (random number):", str1)
	}
	if str2, err := randutil.String(10, randutil.Alphabet); err == nil {
		fmt.Println("get string of 10 chars (random alphabet):", str2)
	}
	if str3, err := randutil.String(10, randutil.Alphanumeric); err == nil {
		fmt.Println("get string of 10 chars (random number and alphabet):", str3)
	}
}

// demo, init random bytes
func testInitBytes() {
	buf := initBytesBySize(32)
	fmt.Printf("init bytes print as numbers: %d\n", buf)
	fmt.Printf("init bytes print as chars: %c\n", buf)

	str := base64.StdEncoding.EncodeToString(buf)
	fmt.Printf("init bytes print as base64 string: %s\n", str)
}

func initBytesBySize(size int) []byte {
	// init []byte "buf" with size of zero
	buf := make([]byte, size)
	for i := 0; i < len(buf); i++ {
		buf[i] = uint8(i % 16)
	}
	return buf
}

// demo, slice append and copy
func testSliceAppend() {
	// #1
	s := []int{5}
	fmt.Printf("slice len=%d, cap=%d, val=%p\n", len(s), cap(s), s)
	s = append(s, 7)
	fmt.Printf("slice len=%d, cap=%d, val=%p\n", len(s), cap(s), s)
	s = append(s, 9)
	fmt.Printf("slice len=%d, cap=%d, val=%p\n", len(s), cap(s), s)
	x := append(s, 11)
	fmt.Printf("slice x len=%d, cap=%d, val=%p\n", len(x), cap(x), x)
	fmt.Println("items in x:")
	for i := 0; i < len(x); i++ {
		fmt.Printf("item %d: addr=%p, val=%d\n", i, &x[i], x[i])
	}

	// #2
	y := append(s, 12)
	fmt.Printf("\nslice y len=%d, cap=%d, val=%p\n", len(y), cap(y), y)
	fmt.Println("new items in x:")
	for i := 0; i < len(x); i++ {
		fmt.Printf("item %d: addr=%p, val=%d\n", i, &x[i], x[i])
	}
	fmt.Println("new items in y:")
	for i := 0; i < len(y); i++ {
		fmt.Printf("item %d: addr=%p, val=%d\n", i, &y[i], y[i])
	}

	// #3
	z := make([]int, 4, 4)
	copy(z, y)
	fmt.Printf("\nslice z len=%d, cap=%d, val=%p\n", len(z), cap(z), z)
	fmt.Println("items in copied z:")
	for i := 0; i < len(y); i++ {
		fmt.Printf("item %d: addr=%p, val=%d\n", i, &z[i], z[i])
	}

	// #4
	printSliceInfo := func(s []int) {
		fmt.Printf("\n[func] slice: addr=%p, val=%p\n", &s, s)
		fmt.Println("[func] slice items:")
		for i := 0; i < len(s); i++ {
			fmt.Printf("item %d: addr=%p, val=%d\n", i, &s[i], s[i])
		}
	}
	printSliceInfo(z)
}

// demo, struct reference
type mySuperStruct struct {
	id  uint
	val string
}

type mySubStruct struct {
	super mySuperStruct // by value
	desc  string
}

type mySubStructRef struct {
	super *mySuperStruct // by refrence
	desc  string
}

func testStructRefValue() {
	s := mySuperStruct{
		id:  10,
		val: "test10",
	}

	subVal := mySubStruct{
		super: s,
		desc:  "inherit from super by value",
	}
	fmt.Printf("before => sub struct: %+v\n", subVal)

	subRef := mySubStructRef{
		super: &s,
		desc:  "inherit from super by reference",
	}
	fmt.Printf("before => sub struct ref: %+v\n", subRef.super)

	s.val = strings.ToUpper(s.val)
	fmt.Printf("after => sub struct: %+v\n", subVal)
	fmt.Printf("after => sub struct Ref: %+v\n", subRef.super)
}

// demo, if and map
var fnPrintMsgID = func(id string) {
	fmt.Println("message id:", id)
}

var fnPrintMsgName = func(name string) {
	fmt.Println("message name:", name)
}

func testPrintMsgByCond() {
	tag := "id"
	name := "message01"
	printMsgByIf(tag, name)
	printMsgByMap(tag, name)
}

func printMsgByIf(tag, input string) {
	fmt.Println("\nprint message by if condition.")
	if tag == "id" {
		fnPrintMsgID(input)
	} else if tag == "name" {
		fnPrintMsgName(input)
	} else {
		fmt.Println("invalid argument!")
	}
}

func printMsgByMap(tag, input string) {
	fmt.Println("\nprint message by map.")
	fns := make(map[string]func(string))
	fns["id"] = fnPrintMsgID
	fns["name"] = fnPrintMsgName
	fns[tag](input)
}

// demo, print info by log
func testLogInfoToStdout() {
	logger := log.New(os.Stdout, "test_", log.Ldate|log.Ltime|log.Lshortfile)
	// logger.SetPrefix("test_");
	// logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile);
	fmt.Println("\nlogger flags:", logger.Flags())
	fmt.Println("logger prefix:", logger.Prefix())

	log.Printf("at %d line, output: %s\n", 21, "this is a default log")
	logger.Printf("at %d line, output: %s\n", 22, "this is a custom logger")

	var (
		isError = true
		isPanic = false
	)
	if isError {
		//print();os.Exit(1);
		logger.Fatal("this is a error")
	}
	if isPanic {
		//print();panic();
		logger.Panic("this is a panic")
	}
}

// demo, output info to file by log
func testLogInfoToFile() {
	tmpDir := filepath.Join(os.Getenv("HOME"), "Downloads/tmp_files")
	logFile := filepath.Join(tmpDir, "test_log.txt")
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	logger := log.New(f, "[test] ", log.Ldate|log.Ltime|log.Lshortfile)
	fmt.Println("\nlogger flags:", logger.Flags())
	fmt.Println("logger prefix:", logger.Prefix())

	for i := 0; i < 3; i++ {
		logger.Printf("at %d line, output: %s\n", i, "this is a custom logger")
	}
}

// demo, json keyword "omitempty"
func testJSONOmitEmpty() {
	type project struct {
		Name string `json:"name"`
		URL  string `json:"url"`
		Desc string `json:"desc"`
		Docs string `json:"docs,omitempty"`
	}

	p1 := project{
		Name: "CleverGo",
		URL:  "https://github.com/headwindfly/clevergo",
		Desc: "CleverGo Perf Framework",
		Docs: "https://github.com/headwindfly/clevergo/tree/master/docs",
	}
	if data, err := json.MarshalIndent(p1, "", "  "); err == nil {
		fmt.Println("\nmarshal json string:", string(data))
	}

	p2 := project{
		Name: "CleverGo",
		URL:  "https://github.com/headwindfly/clevergo",
	}
	if data, err := json.MarshalIndent(p2, "", "  "); err == nil {
		fmt.Println("marshal json string with omitempty:", string(data))
	}
}

// demo, bson parser
func testBSONParser() {
	type testStruct struct {
		FH  []byte `bson:"fh"`
		NFH []byte `bson:"nfh"`
	}

	srcFh := "Bpb_fwEAAAB3eK148Y4dFSvzt1ILAAAAMUMVAAAAAAAKqnHPAAAAAAny-rvibYqoFP-lPkI53JfmoIx5"
	srcNfh := "CJYxQxUAAAAAAAny-rvibYqoFP-lPkI53JfmoIx5a29kby10ZXN0LwUAAHJjUUyDsxizWg=="
	fh, err := base64.URLEncoding.DecodeString(srcFh)
	if err != nil {
		panic(err)
	}
	nfh, err := base64.URLEncoding.DecodeString(srcNfh)
	if err != nil {
		panic(err)
	}

	s := testStruct{
		FH:  fh,
		NFH: nfh,
	}
	if data, err := bson.Marshal(&s); err == nil {
		savePath := filepath.Join(os.Getenv("HOME"), "Downloads/tmp_files/fh.test.bson")
		if err := ioutil.WriteFile(savePath, data, 0666); err != nil {
			panic(err)
		}
		fmt.Printf("save bson bin file: %s\nparse bson: 'bsondump fh.test.bson'\n", savePath)
	}
}

// demo, stop routine by chan
func testStopRoutineByChan() {
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("monitor routine is stop")
				return
			case <-time.Tick(time.Second):
				fmt.Println("monitor routine is running ...")
			}
		}
	}()

	time.Sleep(time.Duration(5) * time.Second)
	fmt.Println("stop monitor routine")
	stop <- true
	time.Sleep(time.Duration(3) * time.Second)
	fmt.Println("main routine exit")
}

// demo, stop routine by context
func testStopRoutineByCtx() {
	type ctxKey string
	var key ctxKey = "ctx_1"

	ctx, cancel := context.WithCancel(context.Background())
	valueCtx := context.WithValue(ctx, key, "monitor_1")
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("[%s] monitor routine is cancelled\n", ctx.Value(key))
				return
			case <-time.Tick(time.Second):
				fmt.Printf("[%s] monitor routine is running ...\n", ctx.Value(key))
			}
		}
	}(valueCtx)

	time.Sleep(time.Duration(5) * time.Second)
	fmt.Println("cancel monitor routine")
	cancel()
	time.Sleep(time.Duration(3) * time.Second)
	fmt.Println("main routine exit")
}

// demo, stop multiple routines by context
func testStopRoutinesByCtx() {
	watcher := func(ctx context.Context, name string) {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("[%s] monitor routine is cancelled\n", name)
				return
			case <-time.Tick(time.Second):
				fmt.Printf("[%s] monitor routine is running ...\n", name)
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < 3; i++ {
		go watcher(ctx, fmt.Sprintf("monitor_%d", i))
	}

	time.Sleep(time.Duration(5) * time.Second)
	fmt.Println("cancel all monitor routines")
	cancel()
	time.Sleep(time.Duration(3) * time.Second)
	fmt.Println("main routine exit")
}

// demo, send/receive nil to/from channel
func testSendNilToChan() {
	ch := make(chan interface{})

	go func() {
		for i := 0; i < 10; i++ {
			ch <- nil
			time.Sleep(time.Second)
		}
		close(ch)
	}()

	fmt.Println("\nreceive chan values:")
	for v := range ch {
		if v == nil {
			fmt.Println("chan val: nil")
		} else {
			fmt.Println("chan val:", v)
		}
	}
}

// demo, use channel as semaphore
func testChanAsSemaphore() {
	fnRoom := func(chToken chan struct{}, name string) {
		chToken <- struct{}{}
		fmt.Println(name, "get token, and in room")
		time.Sleep(time.Duration(2) * time.Second)
		<-chToken
		fmt.Println(name, "release token, and out room")
	}

	chToken := make(chan struct{}, 3)
	for i := 0; i < 10; i++ {
		go func(i int) {
			name := fmt.Sprintf("[routine_%d]", i)
			fmt.Println(name, "is start")
			fnRoom(chToken, name)
			fmt.Println(name, "is end")
		}(i)
	}

	time.Sleep(time.Second)
	for {
		fmt.Println("token size:", len(chToken))
		if len(chToken) == 0 {
			fmt.Println("main done")
			break
		}
		fmt.Println("main sleep ...")
		time.Sleep(time.Second)
	}
}

// demo, get routines count
func testGetGoroutinesCount() {
	printRoutineCount := func() {
		fmt.Println("***** goroutines count:", runtime.NumGoroutine())
	}

	printRoutineCount() // 1
	const waitTime = 5
	ch := make(chan int, 10)
	for i := 0; i < 10; i++ {
		go func(ch chan<- int, num int) {
			sleep, err := randutil.IntRange(2, waitTime)
			if err != nil {
				fmt.Println(err)
				sleep = waitTime
			}
			time.Sleep(time.Duration(sleep) * time.Second)
			ch <- num
		}(ch, i)
	}

	go func(ch chan int) {
		time.Sleep(time.Duration(waitTime+2) * time.Second)
		printRoutineCount() // 2
		fmt.Println("close channel")
		close(ch)
	}(ch)

	time.Sleep(time.Second)
	printRoutineCount() // 12 (10 + 1 + 1)

	for num := range ch {
		fmt.Println("iterator at:", num)
	}
	time.Sleep(time.Second)
	printRoutineCount() // 1
	fmt.Println("testGetGoroutinesCount DONE.")
}

// MainDemo04 main for golang demo04.
func MainDemo04() {
	// testGetFileName()
	// testVerifyGoVersion()
	// testTimeOpSub()
	// testRandomValues()
	// testInitBytes()

	// testSliceAppend()
	// testStructRefValue()
	// testPrintMsgByCond()

	// testLogInfoToStdout()
	// testLogInfoToFile()

	// testJSONOmitEmpty()
	// testBSONParser()

	// testStopRoutineByChan()
	// testStopRoutineByCtx()
	// testStopRoutinesByCtx()

	// testSendNilToChan()
	// testChanAsSemaphore()
	// testGetGoroutinesCount()

	fmt.Println("golang demo04 DONE.")
}
