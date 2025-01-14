package base64toraw

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var phrases = []struct {
	encoding     string
	want         string
	encodedBytes []byte
}{
	{
		encoding:     "std",
		want:         `Bonjour, joyeux lion`,
		encodedBytes: []byte(`Qm9uam91ciwgam95ZXV4IGxpb24=`),
	},
	{
		encoding:     "raw",
		want:         `Bonjour, joyeux lion`,
		encodedBytes: []byte(`Qm9uam91ciwgam95ZXV4IGxpb24`),
	},
	{
		encoding:     "std",
		want:         `Bonjour joyeux lion`,
		encodedBytes: []byte(`Qm9uam91ciBqb3lldXggbGlvbg==`),
	},
	{
		encoding:     "raw",
		want:         `Bonjour joyeux lion`,
		encodedBytes: []byte(`Qm9uam91ciBqb3lldXggbGlvbg`),
	},
	{
		encoding:     "std",
		want:         `Bonjour lion`,
		encodedBytes: []byte(`Qm9uam91ciBsaW9u`),
	},
	{
		encoding:     "raw",
		want:         `Bonjour lion`,
		encodedBytes: []byte(`Qm9uam91ciBsaW9u`),
	},
	{
		encoding: "std",
		want: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		encodedBytes: []byte(`VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wcyBvdmVyIGEgbGF6eSBkb2cuCkdsaWIgam9ja3MgcXVp
eiBueW1waCB0byB2ZXggZHdhcmYuClNwaGlueCBvZiBibGFjayBxdWFydHosIGp1ZGdlIG15IHZv
dy4KSG93IHZleGluZ2x5IHF1aWNrIGRhZnQgemVicmFzIGp1bXAhClRoZSBmaXZlIGJveGluZyB3
aXphcmRzIGp1bXAgcXVpY2tseS4KSmFja2Rhd3MgbG92ZSBteSBiaWcgc3BoaW54IG9mIHF1YXJ0
ei4KUGFjayBteSBib3ggd2l0aCBmaXZlIGRvemVuIGxpcXVvciBqdWdzLg==`),
	},
	{
		encoding: "std",
		want: `The quick brown fox jumps over a lazy dog.
Glib jocks quiz nymph to vex dwarf.
Sphinx of black quartz, judge my vow.
How vexingly quick daft zebras jump!
The five boxing wizards jump quickly.
Jackdaws love my big sphinx of quartz.
Pack my box with five dozen liquor jugs.`,
		encodedBytes: []byte(`VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wcyBvdmVyIGEgbGF6eSBkb2cuCkdsaWIgam9ja3MgcXVp
eiBueW1waCB0byB2ZXggZHdhcmYuClNwaGlueCBvZiBibGFjayBxdWFydHosIGp1ZGdlIG15IHZv
dy4KSG93IHZleGluZ2x5IHF1aWNrIGRhZnQgemVicmFzIGp1bXAhClRoZSBmaXZlIGJveGluZyB3
aXphcmRzIGp1bXAgcXVpY2tseS4KSmFja2Rhd3MgbG92ZSBteSBiaWcgc3BoaW54IG9mIHF1YXJ0
ei4KUGFjayBteSBib3ggd2l0aCBmaXZlIGRvemVuIGxpcXVvciBqdWdzLg==`),
	},
}

func TestDecodingToRaw(t *testing.T) {
	for i, p := range phrases {
		b64 := NewBase64ToRaw(bytes.NewReader(p.encodedBytes))
		b, err := io.ReadAll(base64.NewDecoder(base64.RawStdEncoding, b64))
		if err != nil {
			ee := bytes.ReplaceAll(p.encodedBytes, []byte(" "), []byte("*"))
			t.Fatalf("err %s\n for %s", err, string(ee))
		}
		if got, want := string(b), p.want; cmp.Diff(want, got) != "" {
			t.Errorf("case %d diff\n%s\n", i, cmp.Diff(want, got))
		}
	}
}

func TestFromFileWithEmptyLines(t *testing.T) {

	// files with empty lines
	for _, jf := range []string{"testdata/j1.b64", "testdata/j2.b64", "testdata/j3.b64", "testdata/j4.b64", "testdata/e1.b64"} {

		input, err := ioutil.ReadFile(jf)
		if err != nil {
			t.Fatal("file", jf, err)
		}

		contentReader := base64.NewDecoder(
			base64.RawStdEncoding,
			base64toraw.NewBase64ToRaw(input),
		)
		bb, err := io.ReadAll(contentReader)
		if err != nil {
			t.Fatal(err)
		}
	}
}
