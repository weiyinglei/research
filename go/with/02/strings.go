/**
 * working-with-strings.go
 *
 * Working with strings, manipulating, creating etc
 * Most strings functions are stored in standard library "strings"
 * See: http://golang.org/pkg/strings/
 */

// standard main package
package main

// Note: if you include a package but don't use it, the Go compiler will barf
import (
	"fmt"     // for standard output
	s "strings"// for manipulating strings
)

var p = fmt.Println

func main()  {
	main1()
	main2()
}

// this function gets run on program execution
func main1() {

	// create a string variable
	str := "HI, I'M UPPER CASE"

	// convert to lower case
	lower := s.ToLower(str)

	// output to show its really lower case
	p(lower)

	// check if string contains another string
	if s.Contains(lower, "case") {
		p("Yes, exists!")
	}

	// strings are arrays of characters
	// printing out characters 3 to 9
	p("Characters 3-9: " + str[3:9])

	// printing out first 5 characters
	p("First Five: " + str[:5])

	// split a string on a specific character or word
	sentence := "I'm a sentence made up of words"
	words := s.Split(sentence, " ")
	fmt.Printf("%v \n", words)

	// If you were splitting on whitespace, using Fields is better because
	// it will split on more than just the space, but all whitespace chars
	fields := s.Fields(sentence)
	fmt.Printf("%v \n", fields)

}

func main2() {
	//这是一些 strings 中的函数例子。注意他们都是包中的函数，不是字符串对象自身的方法，这意味着我们需要考虑在调用时传递字符作为第一个参数进行传递。
	p("Contains:  ", s.Contains("test", "es"))
	p("Count:     ", s.Count("test", "t"))
	p("HasPrefix: ", s.HasPrefix("test", "te"))
	p("HasSuffix: ", s.HasSuffix("test", "st"))
	p("Index:     ", s.Index("test", "e"))
	p("Join:      ", s.Join([]string{"a", "b"}, "-"))
	p("Repeat:    ", s.Repeat("a", 5))
	p("Replace:   ", s.Replace("foo", "o", "0", -1))
	p("Replace:   ", s.Replace("foo", "o", "0", 1))
	p("Split:     ", s.Split("a-b-c-d-e", "-"))
	p("ToLower:   ", s.ToLower("TEST"))
	p("ToUpper:   ", s.ToUpper("test"))
	p()
	p("Len: ", len("hello"))
	p("Char:", "hello"[1])
	p(s.NewReader("Hello World!"))
}

// run program in your terminal using
// $ go run 02-working-with-strings.go