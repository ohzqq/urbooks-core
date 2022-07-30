package urbooks

import (
	"time"

	"github.com/ohzqq/urbooks-core/book"
)

func NewBookFeed(title string, resp BookResponse) *book.RSS {
	rss := book.NewFeed()
	channel := rss.SetChannel()
	channel.SetLanguage("en")
	channel.SetTitle(title)
	channel.SetLink(resp.GetResponseLink("self"))
	channel.SetPubdate(time.Now().String())
	channel.SetDescription(title)

	rss.AddCategory("Audiobooks")

	for _, b := range resp.Books {
		channel.AddItem(book.BookToRssItem(b))
	}

	return rss
}
