package util

// Liste de tous les chemins disponibles à travers le site
const (
	BrowsePath       = "/browse"
	BrowseTagPath    = BrowsePath + "/tag/:" + TagParam
	BrowseAuthorPath = BrowsePath + "/author/:" + AuthorParam
	BrowseAllPath    = BrowsePath + "/all"
	BrowseLatestPath = BrowsePath + "/latest"
	BrowseLikedPath  = BrowsePath + "/liked"

	BookPath       = "/book/:" + BookParam
	BookCoverPath  = BookPath + "/cover"
	BookBorrowPath = BookPath + "/borrow"
	BookReturnPath = BookPath + "/return"
	BookLikePath   = BookPath + "/like"
	BookUnlikePath = BookPath + "/unlike"
	BookReviewPath = BookPath + "/review"
	//BookEditPath   = BookPath + "/edit"

	SeriePath      = "serie/:" + SerieParam
	SerieCoverPath = SeriePath + "/cover"

	ContactPath       = "/contact"
	ContactTicketPath = ContactPath + "/ticket"

	LoginPath         = "/login"
	CallbackLoginPath = LoginPath + "/callback"
	LogoutPath        = "/logout"

	AdminPath            = "/admin"
	AdminSeriePath       = AdminPath + "/serie"
	AdminCreateSeriePath = AdminSeriePath + "/create"

	AccountPath = "/account"
)
