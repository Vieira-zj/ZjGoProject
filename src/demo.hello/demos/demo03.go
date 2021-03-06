package demos

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// demo, map entry
func testCheckMapEntry() {
	m := map[int]string{
		1: "one",
		2: "two",
	}

	idx := 2
	fmt.Printf("\nentry[%d] value: %s\n", idx, m[idx])
	fmt.Printf("entry[%d] length: %d\n", idx, len(m[idx]))
	fmt.Printf("entry[%d] 1st char: %c\n", idx, m[idx][0])

	if entry, ok := m[3]; ok {
		fmt.Println("entry[3] value:", entry)
	}
}

// demo, iterator for chars
func testIteratorChars() {
	s := "hello"
	for _, c := range s {
		fmt.Printf("%c", c)
	}
	fmt.Println()

	b := []byte("world")
	fmt.Printf("b type: %T\n", b) // type: []uint8
	for _, c := range b {
		fmt.Printf("%c", c)
	}
	fmt.Println()
}

// demo, value and reference variable
func testValueAndRefVar() {
	arr := [5]int{1, 2, 3, 4, 5}
	fmt.Printf("\narray: addr=%p, val_addr=%p, val=%v\n", &arr, arr, arr)
	fmt.Printf("array item[0]: addr=%p, val=%d\n", &arr[0], arr[0])

	s := []int{1, 2, 3, 4, 5}
	fmt.Printf("\nslice: addr=%p, val_addr=%p, val=%v\n", &s, s, s)
	fmt.Printf("slice item[0]: addr=%p, val=%d\n", &s[0], s[0])

	m := make(map[int]string, 2)
	fmt.Printf("\nmap: addr=%p, val_addr=%p, val=%v\n", &m, m, m)
}

// demo, custom reader (override Read())
type alphaReader1 struct {
	src string
	cur int
}

// Read reads bytes from current position, and copy to p.
func (a *alphaReader1) Read(p []byte) (int, error) {
	if a.cur >= len(a.src) {
		return 0, io.EOF
	}

	x := len(a.src) - a.cur
	bound := 0
	if x >= len(p) {
		bound = len(p)
	} else {
		bound = x
	}

	buf := make([]byte, bound)
	for n := 0; n < bound; n++ {
		if char := alpha(a.src[a.cur]); char != 0 {
			buf[n] = char
		}
		a.cur++
	}
	copy(p, buf)
	return bound, nil
}

func alpha(r byte) byte {
	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
		return r
	}
	return 0
}

func newAlphaReader1(src string) *alphaReader1 {
	return &alphaReader1{src: src}
}

func testCustomAlphaReader1() {
	reader := newAlphaReader1("Hello! It's 9am, where is the sun?")
	p := make([]byte, 4)
	var b []byte

	for {
		n, err := reader.Read(p)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		// fmt.Print(string(p[:n]))
		b = append(b, p[:n]...)
	}
	fmt.Println("\noutput:", string(b))
}

// demo, custom reader (override Read())
type alphaReader2 struct {
	reader io.Reader
}

func (a *alphaReader2) Read(p []byte) (int, error) {
	n, err := a.reader.Read(p)
	if err != nil {
		return n, err
	}

	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		if char := alpha(p[i]); char != 0 {
			buf[i] = char
		}
	}
	copy(p, buf)
	return n, nil
}

func newAlphaReader2(reader io.Reader) *alphaReader2 {
	return &alphaReader2{reader: reader}
}

func testCustomAlphaReader2() {
	reader := newAlphaReader2(strings.NewReader("Hello! It's 9am, where is the sun?"))
	p := make([]byte, 4)

	for {
		n, err := reader.Read(p)
		if err != nil {
			if err == io.EOF {
				fmt.Print(string(p[:n]))
				break
			}
			panic(err.Error())
		}
		fmt.Print(string(p[:n]))
	}
	fmt.Println()
}

// demo, custom writer (override Write())
type chanWriter struct {
	ch chan byte
}

func (w *chanWriter) Chan() <-chan byte {
	return w.ch
}

func (w *chanWriter) Write(p []byte) (int, error) {
	n := 0
	for _, b := range p {
		w.ch <- b
		n++
	}
	return n, nil
}

func (w *chanWriter) Close() error {
	close(w.ch)
	return nil
}

func newChanWriter() *chanWriter {
	// return &chanWriter{ch: make(chan byte, 256)}
	return &chanWriter{make(chan byte, 256)}
}

func testCustomChanWriter() {
	writer := newChanWriter()

	for i := 0; i < 10; i++ {
		go func(idx int) {
			writer.Write([]byte(fmt.Sprintf("Stream%d:", idx)))
			writer.Write([]byte("data\n"))
		}(i)
	}

	go func() {
		time.Sleep(time.Duration(2) * time.Second)
		fmt.Println("close chan writer")
		writer.Close()
	}()

	for c := range writer.Chan() {
		fmt.Printf("%c", c)
	}
	fmt.Println()
}

// demo, time ticker in select block
func testSelectTimeTicker01() {
	ticker := time.NewTicker(time.Duration(3) * time.Second)
	for i := 0; i < 10; i++ {
		select {
		case time := <-ticker.C:
			fmt.Printf("ticker time: %v\n", time)
		default: // not block
			fmt.Println("wait 1 sec ...")
			time.Sleep(time.Second)
		}
	}
	ticker.Stop()
}

func testSelectTimeTicker02() {
	tick := time.Tick(time.Duration(3) * time.Second)
	for i := 0; i < 10; i++ {
		select {
		case time := <-tick:
			fmt.Printf("tick time: %d:%d\n", time.Hour(), time.Minute())
		default: // not block
			fmt.Println("wait 1 sec ...")
			time.Sleep(time.Second)
		}
	}
}

// demo, time after in select block
func testSelectTimeAfter() {
	ch := make(chan string)
	go func() {
		wait := 10
		fmt.Printf("wait %d second in go routine...\n", wait)
		time.Sleep(time.Duration(wait) * time.Second)
		ch <- "done"
	}()

	select {
	case ret := <-ch:
		fmt.Println("return from routine:", ret)
	case <-time.After(time.Duration(3) * time.Second):
		fmt.Printf("3 seconds timeout!\n")
	}
}

// demo, channel as pipeline
func testChanPipeline() {
	naturals := make(chan int)
	squares := make(chan int)

	go func() {
		for i := 0; i < 20; i++ {
			naturals <- i
			time.Sleep(time.Duration(200) * time.Millisecond)
		}
		close(naturals)
	}()

	go func() {
		for {
			x, ok := <-naturals
			if !ok {
				break
			}
			squares <- x * x
		}
		close(squares)
	}()

	fmt.Println("\nsquares:")
	for i := range squares {
		fmt.Printf("%d,", i)
	}
	fmt.Println()
}

// demo, channel as queue
func testChanQueue() {
	const cap = 5
	queue := make(chan int, cap)
	for i := 0; i < cap; i++ {
		queue <- rand.Intn(10)
		time.Sleep(time.Duration(300) * time.Millisecond)
	}

	go func() {
		for i := 0; i < 10; i++ {
			queue <- rand.Intn(20)
			time.Sleep(time.Duration(300) * time.Millisecond)
		}
		close(queue)
	}()

	fmt.Println("queue value:")
	for v := range queue {
		fmt.Println(v)
	}
}

// demo, bufferred channel
func testBufferedChan() {
	queue := make(chan int, 10)
	go func() {
		producers(queue)
	}()
	go func() {
		consumer(queue)
	}()

	for i := 0; i < 15; i++ {
		fmt.Println("queue size:", len(queue))
		time.Sleep(time.Second)
	}
	fmt.Println("close queue")
	close(queue)
}

func producers(queue chan<- int) {
	for {
		select {
		case queue <- rand.Intn(10):
			fmt.Println("true => enqueued without blocking")
		default:
			fmt.Println("false => not enqueued, would have blocked because of queue is full")
		}
		time.Sleep(time.Duration(500) * time.Millisecond)
	}
}

func consumer(queue <-chan int) {
	// OUTER:
	for {
		select {
		case item, valid := <-queue:
			if valid {
				fmt.Println("ok && valid => item is good, use it")
				fmt.Printf("pop off item: %d\n", item)
			} else {
				fmt.Println("ok && !valid => channel closed, quit polling")
			}
			// break OUTER
		default:
			fmt.Println("!ok => channel open, but empty, try later")
		}
		time.Sleep(time.Second)
	}
}

// demo, sync lock and Rlock
func testLockAndRlock() {
	const count = 4
	var mutex sync.RWMutex
	// mutex := new(sync.RWMutex)
	channel := make(chan int, count)

	// lock
	go func(c chan<- int) {
		fmt.Println("\nNot write lock")
		mutex.Lock()
		defer mutex.Unlock()

		fmt.Println("Write Locked")
		time.Sleep(time.Second)
		fmt.Println("Unlock the write lock")
		c <- 10
		fmt.Printf("channel cap=%d, size=%d\n", cap(channel), len(channel))
	}(channel)

	// Rlock
	for i := 0; i < count; i++ {
		go func(i int, c chan<- int) {
			fmt.Println("Not read lock: ", i)
			mutex.RLock()
			defer mutex.RUnlock()

			fmt.Println("Read Locked: ", i)
			time.Sleep(time.Second)
			fmt.Println("Unlock the read lock: ", i)
			c <- i
			fmt.Printf("channel cap=%d, size=%d\n", cap(channel), len(channel))
		}(i, channel)
	}

	time.Sleep(time.Duration(3) * time.Second)
	for i := 0; i < count+1; i++ {
		fmt.Println("output:", <-channel)
	}
}

// demo, function as variable
func testFuncVariable() {
	fmt.Printf("\nadd results: %d\n", myCalculation01(2, 2, funcMyAdd))
	fmt.Printf("min results: %d\n", myCalculation01(2, 8, funcMyMin))

	fmt.Printf("\nadd results: %d\n", myCalculation02(2, 2, funcMyAdd))
	fmt.Printf("min results: %d\n", myCalculation02(2, 8, funcMyMin))
}

func myCalculation01(num1, num2 int, fnCal func(n1, n2 int) int) int {
	return fnCal(num1, num2)
}

type calculateFunc func(n1, n2 int) int

func myCalculation02(num1, num2 int, fnCal calculateFunc) int {
	return fnCal(num1, num2)
}

func funcMyAdd(num1, num2 int) int {
	return num1 + num2
}

func funcMyMin(num1, num2 int) int {
	ret := num1 - num2
	return int(math.Abs(float64(ret)))
}

// demo, function decoration
type apiResponse struct {
	RetCode uint16
	Body    string
	Err     error
}

type apiArgsUser struct {
	UID      uint32
	UserName string
}

func mockAPIPass(args interface{}) *apiResponse {
	info := args.(apiArgsUser)
	content := fmt.Sprintf("user info: Uid=%d, name=%s", info.UID, info.UserName)
	return &apiResponse{
		RetCode: 200,
		Body:    content,
		Err:     nil,
	}
}

type apiArgsGroup struct {
	GID       uint32
	GroupName string
}

func mockAPIFailed(args interface{}) *apiResponse {
	info := args.(apiArgsGroup)
	content := fmt.Sprintf("group not found: Gid=%d, name=%s", info.GID, info.GroupName)
	return &apiResponse{
		RetCode: 204,
		Body:    content,
		Err:     errors.New("EOF"),
	}
}

func testDecorateAPIs() {
	fmt.Println("\n#1. decoration sample: pass")
	{
		args := apiArgsUser{
			UID:      101,
			UserName: "Henry",
		}
		resp := assertAPIs(args, mockAPIPass)
		fmt.Println("pass with resp body:", resp.Body)
	}

	fmt.Println("\n#2. decoration sample: failed")
	{
		args := apiArgsGroup{
			GID:       8,
			GroupName: "QA",
		}
		resp := assertAPIs(args, mockAPIFailed)
		fmt.Println("failed with resp body:", resp.Body)
	}
}

// decoration 装饰器
func assertAPIs(args interface{}, fn func(args interface{}) *apiResponse) *apiResponse {
	resp := fn(args)
	fmt.Printf("response: %+v\n", *resp)
	if resp.RetCode != 200 {
		fmt.Println("failed with ret code:", resp.RetCode)
	}
	if resp.Err != nil {
		fmt.Println("failed with error:", resp.Err.Error())
	}

	return resp
}

// MainDemo03 main for golang demo03.
func MainDemo03() {
	// testCheckMapEntry()
	// testIteratorChars()
	// testValueAndRefVar()

	// testCustomAlphaReader1()
	// testCustomAlphaReader2()
	// testCustomChanWriter()

	// testSelectTimeTicker01()
	// testSelectTimeTicker02()
	// testSelectTimeAfter()

	// testChanPipeline()
	// testChanQueue()
	// testBufferedChan()
	// testLockAndRlock()

	// testFuncVariable()
	// testDecorateAPIs()

	fmt.Println("golang demo03 DONE.")
}
