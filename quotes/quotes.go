package quotes

import (
	"time"
)

type Quote struct {
	Id          int
	Text        string
	UserId      int
	PublishDate time.Time
}

func NewQuote(id int, text string, userId int, publishDate time.Time) Quote {

	return Quote{
		Id:          id,
		Text:        text,
		UserId:      userId,
		PublishDate: publishDate,
	}
}
