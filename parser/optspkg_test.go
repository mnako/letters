package parser_test // this is a package test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mnako/letters/email"
	"github.com/mnako/letters/parser"
)

func TestPkgOptVerbose(t *testing.T) {
	opt := parser.WithVerbose()
	_ = parser.NewParser(opt)
}

// TestPkgOptFileFilter shows an example of a file filter to only save
// jpeg images from an email.
func TestPkgOptFileFilter(t *testing.T) {

	expectedNames := []string{"cat1.jpg", "cat3.jpg"}

	tempDir, err := os.MkdirTemp("", "letters-pkg-tmpdir-*")
	if err != nil {
		t.Fatal(err)
	}

	// customJPGFileFunc selects only attached or inline files with the
	// ContentType of image/jpeg
	customJPGFileFunc := func(ef *email.File) error {
		fcc := strings.ToLower(ef.ContentInfo.Type)
		if fcc != "image/jpeg" {
			return nil
		}
		f, err := os.Create(filepath.Join(tempDir, ef.Name))
		if err != nil {
			return fmt.Errorf("file creation error %w", err)
		}
		defer f.Close()
		_, err = io.Copy(f, ef.Reader)
		if err != nil {
			return fmt.Errorf("file saving error %w", err)
		}
		return nil
	}

	// make option, add to parser constructor, parse email
	dirOpt := parser.WithCustomFileFunc(customJPGFileFunc)
	p := parser.NewParser(dirOpt)
	c, err := os.Open("testdata/cats.eml")
	if err != nil {
		t.Fatal(err)
	}
	_, err = p.Parse(c)
	if err != nil {
		t.Fatal(err)
	}

	// check the output
	dir, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	files := []string{}
	for _, f := range dir {
		files = append(files, f.Name())
	}
	slices.Sort(files)
	if diff := cmp.Diff(expectedNames, files); diff != "" {
		t.Error(diff)
	}
	err = os.RemoveAll(tempDir)
	if err != nil {
		t.Fatal(err)
	}
}
