package vtproto

import (
	"math/rand"
	reflect "reflect"
	"testing"
	"unsafe"
)

var testStrings [][]byte

func init() {
	testStrings = make([][]byte, 1000)
	for i := range testStrings {
		bytestring := make([]byte, rand.Intn(16)+4)
		for j := range bytestring {
			bytestring[j] = byte(rand.Intn(256))
		}
		testStrings[i] = bytestring
	}
}

func TestIntern(t *testing.T) {
	var strings []string
	// Intern all the strings
	for _, bs := range testStrings {
		strings = append(strings, Intern(bs))
	}
	// Intern them again and verify that we always get the same string
	for i, bs := range testStrings {
		s1 := strings[i]
		s2 := Intern(bs)
		if s1 != s2 {
			t.Fatal("Strings are not equal")
		}

		sp1 := (*reflect.StringHeader)(unsafe.Pointer(&s1)).Data
		sp2 := (*reflect.StringHeader)(unsafe.Pointer(&s1)).Data
		if sp1 != sp2 {
			t.Fatal("Strings do not share memory")
		}
	}
}

func BenchmarkIntern(b *testing.B) {
	// Intern all the strings once before starting
	for _, bs := range testStrings {
		Intern(bs)
	}
	b.ResetTimer()
	// Now test how long it takes to get them
	var s string
	for i := 0; i < b.N; i++ {
		for _, bs := range testStrings {
			s = Intern(bs)
		}
	}
	if s != string(testStrings[len(testStrings)-1]) {
		b.Fatal("final strings do not match")
	}
}

// The alterntive to Intern: just convert bytes to a newly allocated string
func Convert(bytes []byte) string {
	return string(bytes)
}

func BenchmarkNoIntern(b *testing.B) {
	// Now test how long it takes to allocate a non-interned string
	var s string
	for i := 0; i < b.N; i++ {
		for _, bs := range testStrings {
			s = Convert(bs)
		}
	}
	if s != string(testStrings[len(testStrings)-1]) {
		b.Fatal("final strings do not match")
	}
}
