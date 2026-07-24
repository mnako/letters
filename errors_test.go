package letters_test

import (
	"errors"
	"net/mail"
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
		t.Errorf("unexpected error message: got %q, want %q", err, expectedMessage)
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
		t.Errorf("unexpected error message: got %q, want %q", err, expectedMessage)
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
		t.Errorf("unexpected error message: got %q, want %q", err, expectedMessage)
	}
}
