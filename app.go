package main

import (
	termbox "github.com/nsf/termbox-go"
	"github.com/nico385412/goreader/epub"
)

// app is used to store the current state of the application.
type app struct {
	pager   pager
	book    *epub.Rootfile
	chapter int
}

// run opens a book, renders its contents within the pager, and polls for
// terminal events until an error occurs or an exit event is detected.
func (a *app) run() error {
	if err := termbox.Init(); err != nil {
		return err
	}
	defer termbox.Flush()
	defer termbox.Close()

	if err := a.openChapter(); err != nil {
		return err
	}

	for {
		if err := a.pager.draw(); err != nil {
			return err
		}
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return nil
			case termbox.KeyArrowDown:
				a.pager.scrollDown()
			case termbox.KeyArrowUp:
				a.pager.scrollUp()
			case termbox.KeyArrowRight:
				a.pager.scrollRight()
			case termbox.KeyArrowLeft:
				a.pager.scrollLeft()
			default:
				switch ev.Ch {
				case 'q':
					return nil
				case 'j':
					a.pager.scrollDown()
				case 'k':
					a.pager.scrollUp()
				case 'h':
					a.pager.scrollLeft()
				case 'l':
					a.pager.scrollRight()
				case 'f':
					if a.pager.pageDown() || a.chapter >= len(a.book.Spine.Itemrefs)-1 {
						continue
					}

					// Go to the next chapter if we reached the end.
					if err := a.nextChapter(); err != nil {
						return err
					}
					a.pager.toTop()
				case 'b':
					if a.pager.pageUp() || a.chapter <= 0 {
						continue
					}

					// Go to the previous chapter if we reached the beginning.
					if err := a.prevChapter(); err != nil {
						return err
					}
					a.pager.toBottom()
				case 'g':
					a.pager.toTop()
				case 'G':
					a.pager.toBottom()
				case 'L':
					if a.chapter >= len(a.book.Spine.Itemrefs)-1 {
						continue
					}

					if err := a.nextChapter(); err != nil {
						return err
					}
					a.pager.toTop()
				case 'H':
					if a.chapter <= 0 {
						continue
					}

					if err := a.prevChapter(); err != nil {
						return err
					}
					a.pager.toTop()
				}
			}
		}
	}
}

// openChapter opens the current chapter and renders it within the pager.
func (a *app) openChapter() error {
	f, err := a.book.Spine.Itemrefs[a.chapter].Open()
	if err != nil {
		return err
	}
	doc, err := parseText(f, a.book.Manifest.Items)
	if err != nil {
		return err
	}
	a.pager.doc = doc

	return nil
}

// nextChapter opens the next chapter.
func (a *app) nextChapter() error {
	a.chapter++
	return a.openChapter()
}

// prevChapter opens the previous chapter.
func (a *app) prevChapter() error {
	a.chapter--
	return a.openChapter()
}
