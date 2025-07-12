package diff

import (
	"bytes"
	"io/ioutil"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// Diff represents a diff
type Diff struct {
	Path    string
	Patches []diffmatchpatch.Patch
}

// New creates a new diff
func New(path string, oldContent, newContent []byte) *Diff {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(oldContent), string(newContent), false)
	patches := dmp.PatchMake(diffs)

	return &Diff{
		Path:    path,
		Patches: patches,
	}
}

// Apply applies the diff
func (d *Diff) Apply() error {
	content, err := ioutil.ReadFile(d.Path)
	if err != nil {
		return err
	}

	dmp := diffmatchpatch.New()
	newContent, _ := dmp.PatchApply(d.Patches, string(content))

	return ioutil.WriteFile(d.Path, []byte(newContent), 0644)
}

// ToBytes returns the diff as bytes
func (d *Diff) ToBytes() []byte {
	dmp := diffmatchpatch.New()
	return []byte(dmp.PatchToText(d.Patches))
}

// FromBytes returns a diff from bytes
func FromBytes(path string, data []byte) (*Diff, error) {
	dmp := diffmatchpatch.New()
	patches, err := dmp.PatchFromText(string(data))
	if err != nil {
		return nil, err
	}

	return &Diff{
		Path:    path,
		Patches: patches,
	}, nil
}

// HasChanges returns true if the diff has changes
func (d *Diff) HasChanges() bool {
	return len(d.ToBytes()) > 0
}

// ApplyToContent applies the diff to the given content
func (d *Diff) ApplyToContent(content []byte) ([]byte, error) {
	dmp := diffmatchpatch.New()
	newContent, _ := dmp.PatchApply(d.Patches, string(content))
	return []byte(newContent), nil
}

// CreatePatch creates a patch from the given old and new content
func CreatePatch(oldContent, newContent []byte) []byte {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(oldContent), string(newContent), false)
	patches := dmp.PatchMake(diffs)
	return []byte(dmp.PatchToText(patches))
}

// ApplyPatch applies a patch to the given content
func ApplyPatch(content, patch []byte) ([]byte, error) {
	dmp := diffmatchpatch.New()
	patches, err := dmp.PatchFromText(string(patch))
	if err != nil {
		return nil, err
	}
	newContent, _ := dmp.PatchApply(patches, string(content))
	return []byte(newContent), nil
}

// Merge merges the given old and new content
func Merge(oldContent, newContent []byte) ([]byte, error) {
	patch := CreatePatch(oldContent, newContent)
	return ApplyPatch(oldContent, patch)
}

// MergeFiles merges the given old and new files
func MergeFiles(oldPath, newPath, outputPath string) error {
	oldContent, err := ioutil.ReadFile(oldPath)
	if err != nil {
		return err
	}

	newContent, err := ioutil.ReadFile(newPath)
	if err != nil {
		return err
	}

	mergedContent, err := Merge(oldContent, newContent)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputPath, mergedContent, 0644)
}

// Compare returns true if the given files are different
func Compare(path1, path2 string) (bool, error) {
	content1, err := ioutil.ReadFile(path1)
	if err != nil {
		return false, err
	}

	content2, err := ioutil.ReadFile(path2)
	if err != nil {
		return false, err
	}

	return !bytes.Equal(content1, content2), nil
}
