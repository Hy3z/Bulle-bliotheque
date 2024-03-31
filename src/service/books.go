package service


type BookPreview struct {
	Title string
	Cover string
	Id string
}
type BookPreviewSet []BookPreview
type InfiniteBookPreviewSet struct {
	BookPreviewSet BookPreviewSet
	Url string
	Params []PathParameter
}
const (
	BookPreviewSetTemplate = "book-preview"
	InfiniteBookPreviewSetTemplate = "infinite-book-set"
	BookPreviewTemplate = "book-set"
)

