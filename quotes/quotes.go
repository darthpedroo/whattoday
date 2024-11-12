package quotes

import ("time"
	api "whattoday/web-service-gin/users")


type Quote struct {
	Id int
	Text string
	User api.User
	PublishDate time.Time

}

func NewQuote(id int, text string, user api.User, publishDate time.Time) Quote {

	return Quote{
		Id: id,
		Text: text,
		User: user,
		PublishDate: publishDate,
	}
}