package email

import (
	"errors"
	"fmt"
	"mime"
	"net/textproto"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
)

// Relevant RFCs and related information concerning Content* fields
// which are used to describe MIME (Multipurpose Internet Mail
// Extensions) email messages.
//
// Content-Type:
// https://datatracker.ietf.org/doc/html/rfc2045#page-10
// The purpose of the Content-Type field is to describe the data
// contained in the body fully enough that the receiving user agent can
// pick an appropriate agent or mechanism to present the data to the
// user, or otherwise deal with the data in an appropriate manner. The
// value in this field is called a media type.
// ...
// Parameters are modifiers of the media subtype, and as such do not
// fundamentally affect the nature of the content.  The set of
// meaningful parameters depends on the media type and subtype.
// ...
// For example, the "charset" parameter is applicable to any subtype of
// "text", while the "boundary" parameter is required for any subtype of
// the "multipart" media type.
//
// Content-Type is not required for a document to conform with RFC 2045.
// Content-Type defaults to text/plain. Content-Type defines the type of
// data in each part as a type/subtype.
//
//	content := "Content-Type" ":" type "/" subtype
//	           *(";" parameter)
//	           ; Matching of media type and subtype
//	           ; is ALWAYS case-insensitive.
//	type := discrete-type / composite-type
//	discrete-type := "text" / "image" / "audio" / "video" /
//	                 "application" / extension-token
//	composite-type := "message" / "multipart" / extension-token
//	extension-token := ietf-token / x-token
//	ietf-token := <An extension token defined by a
//	               standards-track RFC and registered
//	               with IANA.>
//	x-token := <The two characters "X-" or "x-" followed, with
//	            no intervening white space, by any token>
//	subtype := extension-token / iana-token
//	iana-token := <A publicly-defined extension token. Tokens
//	               of this form must be registered with IANA
//	               as specified in RFC 2048.>
//	parameter := attribute "=" value
//	attribute := token
//	             ; Matching of attributes
//	             ; is ALWAYS case-insensitive.
//	value := token / quoted-string
//	token := 1*<any (US-ASCII) CHAR except SPACE, CTLs,
//	            or tspecials>
//	tspecials :=  "(" / ")" / "<" / ">" / "@" /
//	              "," / ";" / ":" / "\" / <">
//	              "/" / "[" / "]" / "?" / "="
//	              ; Must be in quoted-string,
//	              ; to use within parameter values
//
// Content-Disposition:
// (RFC RFC2183,
// https://www.iana.org/assignments/cont-disp/cont-disp.xhtml#cont-disp-1)
//	disposition := "Content-Disposition" ":"
//	               disposition-type
//	               *(";" disposition-parm)
//
//	disposition-type := "inline"
//	                  / "attachment"
//	                  / extension-token
//	                  ; values are not case-sensitive
//
//	disposition-parm := filename-parm
//	                  / creation-date-parm
//	                  / modification-date-parm
//	                  / read-date-parm
//	                  / size-parm
//	                  / parameter
//
//	filename-parm := "filename" "=" value
//	creation-date-parm := "creation-date" "=" quoted-date-time
//	modification-date-parm := "modification-date" "=" quoted-date-time
//	read-date-parm := "read-date" "=" quoted-date-time
//	size-parm := "size" "=" 1*DIGIT
//	quoted-date-time := quoted-string
//	                 ; contents MUST be an RFC 822 `date-time'
//	                 ; numeric timezones (+HHMM or -HHMM) MUST be used
//
// Content-Transfer-Encoding:
// https://www.ietf.org/rfc/rfc2045.txt
// Many Content-Types are represented as 8-bit character or binary data,
// and can include XML, which typically uses UTF-8 or UTF-16 encoding.
// This type of data cannot be transmitted over some transport
// protocols, and might be encoded to 7-bit. The
// Content-Transfer-Encoding header field is used to indicate the type
// of transformation that has been used for encoding this type of data
// into a 7-bit format.
//
//	 encoding := "Content-Transfer-Encoding" ":" mechanism
//
//	 mechanism := "7bit" / "8bit" / "binary" /
//	              "quoted-printable" / "base64" /
//	              ietf-token / x-token
//
// Content-ID:
// https://datatracker.ietf.org/doc/html/rfc2392
// Optional. This field enables parts to be labeled, and referenced from
// other parts of the message. These parts are typically referenced from
// part 0 (the first) of the message
//
//	content-id    = url-addr-spec
//	url-addr-spec = addr-spec  ; URL encoding of RFC 822 addr-spec
//
// ...the addr-spec in a Content-ID [MIME] or Message-ID [822] header is
// enclosed in angle brackets (<>).  Since addr-spec in a Message-ID or
// Content-ID might contain characters not allowed within a URL; any
// such character (including "/", which is reserved within the "mid"
// scheme) must be hex-encoded using the %hh escape mechanism in [URL].
// e.g.
//	Content-ID: <foo4%25foo1@bar.net>

// ContentInfo is a struct containing the Content* information
// concerning messages, message-bodies, multi-part message bodies and
// other non-US-ASCII header information.
type ContentInfo struct {
	Type              string            // Content-Type header or mime-part data description
	TypeParams        map[string]string // Content-Type parameters
	Disposition       string            // Content-Disposition header or mime-part data description
	DispositionParams map[string]string // Content-Disposition parameters
	TransferEncoding  string            // Content-Transfer-Encoding header or mime-part data description
	ID                string            // ContentID part labelling
	// additional fields
	Charset  string            // the charset extracted from the content type
	Encoding encoding.Encoding // the encoding determined by the charset
}

// contentDispositions is a slice of valid content
// dispositions.
var contentDispositions = []string{"attachment", "inline"}

func inSlice(s []string, q string) bool {
	for _, o := range s {
		if o == q {
			return true
		}
	}
	return false
}

// contentTransfers is a slice of valid transfer encoding
var contentTransferEncodings = []string{
	"7bit",
	"8bit",
	"binary",
	"quoted-printable",
	"base64",
}

// ExtractContentInfo extracts information from a headers map from
// either a net/mail.Message.Header or mime/multipart.Part.Header, whose
// underlying type is a map[string][]string.
// Fallback information may be provided by a parent ConentInfo instance.
func ExtractContentInfo(headers map[string][]string, parentCI *ContentInfo) (*ContentInfo, error) {
	get := func(key string) string {
		if headers == nil {
			return ""
		}
		v := headers[textproto.CanonicalMIMEHeaderKey(key)]
		if len(v) == 0 {
			return ""
		}
		return v[0]
	}

	c := &ContentInfo{}
	err := c.extractType(get("Content-Type"))
	if err != nil {
		return c, err
	}
	err := c.extractCharset(parentCI)
	if err != nil {
		return c, err
	}
	err := c.extractTransferEncoding(get("Content-Transfer-Encoding"))
	if err != nil {
		return c, err
	}
	err := c.extractDisposition(get("Content-Disposition"))
	if err != nil {
		return c, err
	}
	c.ExtractID(get("Content-ID"))
	return c, nil
}

// IsInlineFile reports if the content type describes an inline file.
func (c *ContentInfo) IsInlineFile(parentCI *ContentInfo) bool {
	switch {
	case c.Disposition == "inline":
		return true
	case c.Type == "text/plain", c.Type == "text/enriched", c.Type == "text/html":
		return false
	case parentCI != nil && parentCI.Type == "multipart/related":
		return true
	}
	return false
}

// IsAttachedFile reports if the content type describes an attached file.
func (c *ContentInfo) isAttachedFile(parentCI *ContentInfo) bool {
	switch {
	case c.Disposition == "attached":
		return true
	case c.Type == "text/plain", c.Type == "text/enriched", c.Type == "text/html":
		return false
	case parentCI != nil && parentCI.Type == "multipart/mixed":
		return true
	case parentCI != nil && parentCI.Type == "multipart/parallel":
		return true
	}
	return false
}

// extractType extracts the Content-Type and Parameter information
func (c *ContentInfo) extractType(s string) error {
	if s == "" {
		s = "text/plain"
	}
	var err error
	c.Type, c.TypeParams, err = mime.ParseMediaType(s)
	if err != nil {
		return fmt.Errorf("cannot extract Content-Type %q: %w", s, err)
	}
	for _, param := range []string{"charset", "micalg", "protocol"} {
		if v, ok := c.Params[param]; ok {
			c.Params[param] = strings.ToLower(v)
		}
	}
	return nil
}

// extractCharset extracts the charset from the Content Type or parent
// Content Type
func (c *ContentInfo) extractCharset(parentCI *ContentInfo) error {
	if c.Type == "" {
		return errors.New("content type empty, cannot extract charset")
	}
	c.Charset = c.TypeParams["charset"]
	if charsetLabel == "" && parentCI != nil {
		c.Charset = parentCI.TypeParams["charset"]
	}
	if c.Charset == "" {
		return nil
	}
	c.Encoding, _ := charset.Lookup(c.Charset)
	return nil
}

// extractDisposition extracts the Content-Disposition and Parameter information
func (c *ContentInfo) extractDisposition(s string) error {
	if s == "" {
		return nil
	}
	var err error
	c.Disposition, c.DispositionParams, err = mime.ParseMediaType(s)
	if err != nil {
		return fmt.Errorf("cannot extract Content-Disposition %q: %w", s, err)
	}
	if !inSlice(contentDispositions, c.Disposition) {
		return fmt.Errorf("unknown Content-Disposition %q", c.Disposition)
	}
	return nil
}

// extractTransferEncoding extracts the Content-Transfer-Encoding
func (c *ContentInfo) ExtractTransferEncoding(s string) error {
	c.TransferEncoding := strings.ToLower(strings.TrimSpace(s))
	if label == "" {
		c.TransferEncoding = "7bit"
		return nil
	}

	if !inSlice(contentTransferEncodings, c.TransferEncoding) {
		return fmt.Errorf("unknown Content-Transfer-Encoding %q", label)
	}
	return nil
}

// extractID extracts the ContentID
func (c *ContentInfo) ExtractID(s string) {
	c.ID = strings.TrimSpace(strings.Trim(s, "<>"))
}
