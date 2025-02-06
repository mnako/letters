package parser

import (
	"net/mail"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestOptVerbose(t *testing.T) {
	opt := WithVerbose()
	p := NewParser(opt)
	if p.verbose == false {
		t.Error("expected p.verbose to be true")
	}
}

func TestOptWithoutAttachments(t *testing.T) {
	opt := WithoutAttachments()
	p := NewParser(opt)
	if p.processType != noAttachments {
		t.Errorf("processType got %s want %s", p.processType, noAttachments)
	}
}

func TestOptAddressCustomFunc(t *testing.T) {

	b, a := "Bart Simpson", "<bart@example.com>"
	addressFunc := func(s string) (*mail.Address, error) {
		return &mail.Address{Name: b, Address: a}, nil
	}

	opt := WithCustomAddressFunc(addressFunc)
	p := NewParser(opt)

	addr, err := p.parseAddress("Ronny Burke <ronnie@example.com>")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := addr.Address, a; got != want {
		t.Errorf("got %s want %s", got, want)
	}
	if got, want := addr.Name, b; got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestOptAddressesCustomFunc(t *testing.T) {
	addresses := [][]string{
		[]string{"Bart Simpson", "<bart@example.com>"},
		[]string{"Darth Vader", "<darth@example.com>"},
	}
	addressesFunc := func(stringList string) ([]*mail.Address, error) {
		return []*mail.Address{
			&mail.Address{Name: addresses[0][0], Address: addresses[0][1]},
			&mail.Address{Name: addresses[1][0], Address: addresses[1][1]},
		}, nil
	}

	opt := WithCustomAddressesFunc(addressesFunc)
	p := NewParser(opt)

	results, err := p.parseAddresses("Ronny Burke <ronnie@example.com>, A. S. Byatt <ab@oxon.ac.uk")
	if err != nil {
		t.Fatal(err)
	}
	for i, a := range results {
		name, addr := addresses[i][0], addresses[i][1]
		if got, want := a.Address, addr; got != want {
			t.Errorf("got %s want %s", got, want)
		}
		if got, want := a.Name, name; got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}
}

func TestOptDateCustomFunc(t *testing.T) {
	dateFunc := func(string) (time.Time, error) {
		return time.Time{}, nil
	}
	opt := WithCustomDateFunc(dateFunc)
	p := NewParser(opt)

	tt, err := p.dateFunc("anything")
	if err != nil {
		t.Fatal(err)
	}
	if !tt.IsZero() {
		t.Error("time should be time.IsZero()")
	}
}

func TestOptSaveFilesToDirectory(t *testing.T) {

	expectedNames := []string{"cat1.jpg", "cat2.png", "cat3.jpg"}

	tempDir, err := os.MkdirTemp("", "letters-tmpdir-*")
	if err != nil {
		t.Fatal(err)
	}
	dirOpt := WithSaveFilesToDirectory(tempDir)
	p := NewParser(dirOpt)

	c, err := os.Open("testdata/cats.eml")
	if err != nil {
		t.Fatal(err)
	}
	_, err = p.Parse(c)
	if err != nil {
		t.Fatal(err)
	}

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
