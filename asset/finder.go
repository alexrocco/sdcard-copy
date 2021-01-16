package asset

//go:generate mockgen -source=./finder.go -package=mock -destination=../mock/mock_finder.go Finder

// Finder finds asset
type Finder interface {
	// Find finds the files and return the paths
	Find(regex string) ([]string, error)
}
