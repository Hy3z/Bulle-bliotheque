package util

const(
	BrowsePath = "/browse"

	BrowseTagPath = BrowsePath+"/tag/:"+TagParam
	BrowseAllPath = BrowsePath +"/all"
	BrowseLatestPath = BrowsePath +"/latest"

	BookPath = "/book/:"+IsbnParam
	BookCoverPath = BookPath+"/cover"

	SeriePath = "serie/:"+SerieParam
)
