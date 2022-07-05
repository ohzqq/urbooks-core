package rss

import (
	"time"

	"github.com/lann/builder"
)

type channelBuilder builder.Builder

func NewChannel() channelBuilder {
	return ChannelBuilder
}

func (c channelBuilder) Title(title string) channelBuilder {
	return builder.Set(c, "Title", title).(channelBuilder)
}

func (c channelBuilder) Link(link string) channelBuilder {
	return builder.Set(c, "Link", link).(channelBuilder)
}

func (c channelBuilder) Description(desc string) channelBuilder {
	d := Description{Body: desc}
	return builder.Set(c, "Description", &d).(channelBuilder)
}

func (c channelBuilder) Language(lang string) channelBuilder {
	return builder.Set(c, "Language", lang).(channelBuilder)
}

func (c channelBuilder) AddCategory(category string) channelBuilder {
	return builder.Append(c, "Category", category).(channelBuilder)
}

func (c channelBuilder) PubDate(pubdate time.Time) channelBuilder {
	return builder.Set(c, "PubDate", Time(pubdate)).(channelBuilder)
}

func (c channelBuilder) Image(image string) channelBuilder {
	i := Image{Href: image}
	return builder.Set(c, "Image", i).(channelBuilder)
}

func (c channelBuilder) AddItem(item Item) channelBuilder {
	return builder.Append(c, "Item", &item).(channelBuilder)
}

func (c channelBuilder) Build() Channel {
	return builder.GetStruct(c).(Channel)
}

// ChannelBuilder is a fluent immutable builder to build OPDS Channels
var ChannelBuilder = builder.Register(channelBuilder{}, Channel{}).(channelBuilder)
