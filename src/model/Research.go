package model

type Research struct {
	Name string
	IsInfinite bool
	//Use either of the field below depending on boolean value
	BookPreviewSet BookPreviewSet
	InfiniteBookPreviewSet InfiniteBookPreviewSet
}

const ResearchTemplate = "research"
