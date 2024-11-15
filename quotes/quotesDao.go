package quotes

type QuotesDao interface {
	GetQuotes() ([]CustomQuoteResponse , error)
	AddQuote(quote Quote) error
}