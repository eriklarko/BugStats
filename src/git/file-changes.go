package git

type Change int
const (
	MODIFIED Change = iota
	CREATED
	REMOVED
)

type FileChange struct {
	FileName string
	FileChange Change
}