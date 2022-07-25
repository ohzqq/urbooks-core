package urbooks

//func MediaMetaToBook(lib string, m *avtools.Media) *Book {
//  book := NewBook(lib)
//  titleRegex := regexp.MustCompile(`(?P<title>.*) \[(?P<series>.*), Book (?P<position>.*)\]$`)
//  titleAndSeries := titleRegex.FindStringSubmatch(m.GetTag("title"))

//  book.NewColumn("title").SetValue(titleAndSeries[titleRegex.SubexpIndex("title")])
//  book.NewItem("series").
//    SetValue(titleAndSeries[titleRegex.SubexpIndex("series")]).
//    Set("position", titleAndSeries[titleRegex.SubexpIndex("position")])
//  book.NewCategory("authors").Split(m.GetTag("artist"), true)
//  book.NewCategory("narrators").Split(m.GetTag("composer"), true)
//  book.NewColumn("description").SetValue(m.GetTag("comment"))
//  book.NewCategory("tags").Split(m.GetTag("genre"), false)
//  return book
//}
