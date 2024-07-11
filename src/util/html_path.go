package util

const (
	BrowsePath = "/browse"

	BrowseTagPath    = BrowsePath + "/tag/:" + TagParam
	BrowseAuthorPath = BrowsePath + "/author/:" + AuthorParam
	BrowseAllPath    = BrowsePath + "/all"
	BrowseLatestPath = BrowsePath + "/latest"

	BookPath       = "/book/:" + BookParam
	BookCoverPath  = BookPath + "/cover"
	BookBorrowPath = BookPath + "/borrow"
	BookReturnPath = BookPath + "/return"
	//BookEditPath   = BookPath + "/edit"

	SeriePath      = "serie/:" + SerieParam
	SerieCoverPath = SeriePath + "/cover"

	ContactPath       = "/contact"
	ContactTicketPath = ContactPath + "/ticket"

	LoginPath         = "/login"
	CallbackLoginPath = LoginPath + "/callback"
	LogoutPath        = "/logout"
)
