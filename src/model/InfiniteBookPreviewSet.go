package model

type InfiniteBookPreviewSet struct {
	BookPreviewSet BookPreviewSet
	Url string
	Params []PathParameter
}

const InfiniteBookPreviewSetTemplate = "infinite-book-set"
