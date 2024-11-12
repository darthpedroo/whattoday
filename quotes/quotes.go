package quotes

import ("time")


type Quote struct {
	Id int
	Text string
	PublishDate time.Time
}

func NewQuote(id int, text string, publishDate time.Time) Quote {

	return Quote{
		Id: id,
		Text: text,
		PublishDate: publishDate,
	}
}