package assets

//go:generate mockgen -source=./finder.go -package=mocks -destination=../mocks/mock_finder.go Finder

// Finder finds assets
type Finder interface {
	// Find finds the files and return the paths
	Find(regex string) ([]string, error)
}
