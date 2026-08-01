package letters_test

import (
	"errors"
	"net/mail"
	"os"
	"testing"

	"github.com/mnako/letters"
)

func TestUnknownCharsetError(t *testing.T) {
	t.Parallel()

	const encodedAddress = "=?x-invalid?Q?Alice?= <alice@example.com>"

	header := mail.Header{"From": {encodedAddress}}
	_, err := letters.ParseAddressHeader(header, "From")

	if !errors.Is(err, letters.ErrUnknownCharset) {
		t.Fatalf("expected ErrUnknownCharset, got %v", err)
	}

	const expectedMessage = "letters.parsers.parseAddressHeader: " +
		"cannot decode address header \"=?x-invalid?Q?Alice?= <alice@example.com>\": " +
		"letters.decoders.decodeHeader: " +
		"cannot decode MIME-word-encoded header " +
		"\"=?x-invalid?Q?Alice?= <alice@example.com>\": " +
		"letters.decoders.decodeHeader.CharsetReader: " +
		"cannot lookup encoding x-invalid"

	if err.Error() != expectedMessage {
		t.Errorf(
			"unexpected error message: got %q, want %q",
			err,
			expectedMessage,
		)
	}
}

func TestParseEmailUnknownCharsetError(t *testing.T) {
	t.Parallel()

	rawEmail, err := os.Open(
		"tests/test_english_plaintext_unknown_charset.txt",
	)
	if err != nil {
		t.Fatalf("error while reading email from file: %s", err)
	}

	defer func() {
		if err := rawEmail.Close(); err != nil {
			t.Errorf("error while closing rawEmail: %s", err)
		}
	}()

	_, err = letters.NewEmailParser().Parse(rawEmail)

	if !errors.Is(err, letters.ErrUnknownCharset) {
		t.Fatalf("expected ErrUnknownCharset, got %v", err)
	}
}

func TestUnknownContentDispositionError(t *testing.T) {
	t.Parallel()

	_, err := letters.ParseContentDisposition("unexpected")
	if !errors.Is(err, letters.ErrUnknownContentDisposition) {
		t.Fatalf("expected ErrUnknownContentDisposition, got %v", err)
	}

	const expectedMessage = "letters.parsers.parseContentDisposition: " +
		"unknown Content-Disposition \"unexpected\""

	if err.Error() != expectedMessage {
		t.Errorf(
			"unexpected error message: got %q, want %q",
			err,
			expectedMessage,
		)
	}
}

func TestUnknownContentTransferEncodingError(t *testing.T) {
	t.Parallel()

	_, err := letters.ParseContentTransferEncoding("unexpected")
	if !errors.Is(err, letters.ErrUnknownContentTransferEncoding) {
		t.Fatalf("expected ErrUnknownContentTransferEncoding, got %v", err)
	}

	const expectedMessage = "letters.parsers.parseContentTransferEncoding: " +
		"unknown Content-Transfer-Encoding \"unexpected\""

	if err.Error() != expectedMessage {
		t.Errorf(
			"unexpected error message: got %q, want %q",
			err,
			expectedMessage,
		)
	}
}
