package quotes

import (
	"time"
	"whattoday/web-service-gin/users"
)

type Quote struct {
	Id          int
	Text        string
	UserId      int
	PublishDate time.Time
}

type CustomQuoteResponse struct {
	Quote Quote
	User  users.User
}

func NewQuote(id int, text string, userId int, publishDate time.Time) Quote {

	return Quote{
		Id:          id,
		Text:        text,
		UserId:      userId,
		PublishDate: publishDate,
	}
}
