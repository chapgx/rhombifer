package pages

import (
	"fmt"
	"testing"
)

func TestPages(t *testing.T) {
	data := []byte("Hello Wolrd\n")
	data = append(data, []byte("How Are you\n")...)
	data = append(data, []byte("Bounjoir\n")...)
	data = append(data, []byte("Ciao amore\n")...)
	data = append(data, []byte("Hola como estas\n")...)
	data = append(data, []byte("Gene compre pa\n")...)
	data = append(data, []byte("Foo bar bar foo\n")...)

	pages := NewPages(2, data)

	for _, p := range pages.pages {
		fmt.Printf("records in page %d are %d\n", p.pagenumber, p.Records())
		fmt.Printf("%s\n\n", string(p.data))
	}
}
