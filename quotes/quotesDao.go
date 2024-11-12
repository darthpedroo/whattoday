package quotes

type QuotesDao interface {
	GetQuotes() ([]Quote , error)
	AddQuote(quote Quote) error
}