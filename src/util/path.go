package util

const (
	BrowsePath = "/browse"

	BrowseTagPath    = BrowsePath + "/tag/:" + TagParam
	BrowseAuthorPath = BrowsePath + "/author/:" + AuthorParam
	BrowseAllPath    = BrowsePath + "/all"
	BrowseLatestPath = BrowsePath + "/latest"

	BookPath      = "/book/:" + BookParam
	BookCoverPath = BookPath + "/cover"

	SeriePath      = "serie/:" + SerieParam
	SerieCoverPath = SeriePath + "/cover"

	ContactPath       = "/contact"
	ContactTicketPath = ContactPath + "/ticket"

	AuthPath = "/auth"
)
