package pages

import (
	"bytes"
)

type Page struct {
	pagenumber int
	data       []byte
}

func (p Page) Records() int {
	records := bytes.Split(p.data, []byte("\n"))
	return len(records)
}

type Pages struct {
	pages  []Page
	length int
}

// New Pages
func NewPages(recordsPerPage int, data []byte) Pages {
	pages := Pages{}
	lines := bytes.Split(data, []byte("\n"))
	pages.length = len(lines) / recordsPerPage
	counter := 0
	currpage := 0
	var grouping []byte
	for _, l := range lines {
		if string(l) == "" {
			continue
		}
		if counter >= recordsPerPage {
			currpage += 1
			counter = 0
			page := Page{
				pagenumber: currpage,
				data:       grouping[:len(grouping)-1],
			}
			grouping = []byte{}
			pages.pages = append(pages.pages, page)
		}
		l = append(l, '\n')
		grouping = append(grouping, l...)
		counter += 1
	}

	if pages.length < len(lines) {
		lines = lines[pages.length:]
		lastpage := Page{pagenumber: pages.length + 1}
		for _, l := range lines {
			if string(l) == "" {
				continue
			}
			l = append(l, '\n')
			lastpage.data = append(lastpage.data, l...)
		}
		lastpage.data = lastpage.data[:len(lastpage.data)-1]
		pages.length += 1
		pages.pages = append(pages.pages, lastpage)
	}

	return pages
}
