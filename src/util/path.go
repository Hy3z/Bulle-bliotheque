package util

const (
	BrowsePath = "/browse"

	BrowseTagPath    = BrowsePath + "/tag/:" + TagParam
	BrowseAuthorPath = BrowsePath + "/author/:" + AuthorParam
	BrowseAllPath    = BrowsePath + "/all"
	BrowseLatestPath = BrowsePath + "/latest"

	BookPath      = "/book/:" + IsbnParam
	BookCoverPath = BookPath + "/cover"

	SeriePath = "serie/:" + SerieParam

	ContactPath       = "/contact"
	ContactTicketPath = ContactPath + "/ticket"

	AuthPath = "/auth"
)
