package gigi

type File struct {
	Name       string
	AddedCount int
}

type Result struct {
	TotalAddedCount int
	Files           []File
}
