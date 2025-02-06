package parser

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/mnako/letters/decoders"
)

var (
	errorEmptyAddress error = errors.New("Empty Address")
	errorEmptyDate    error = errors.New("Empty Date")
)

// explicitHeaders are those headers stored in their own field in
// email.Headers, rather than in email.Headers.ExtraHeaders
var explicitHeaders = []string{
	"Date",
	"Sender",
	"From",
	"Reply-To",
	"To",
	"Cc",
	"Bcc",
	"Message-Id",
	"In-Reply-To",
	"References",
	"Received",
	"Subject",
	"Comments",
	"Keywords",
	"Resent-Date",
	"Resent-From",
	"Resent-Sender",
	"Resent-To",
	"Resent-Cc",
	"Resent-Bcc",
	"Resent-Message-Id",
	"Content-Transfer-Encoding",
	"Content-Type",
	"Content-Disposition",
}

// isExplicitHeader checks if the header is to be registered as a field.
// This slice search is much the same speed as a map lookup for small
// slices.
func isExplicitHeader(s string) bool {
	for _, e := range explicitHeaders {
		if e == s {
			return true
		}
	}
	return false
}

// idTrimCutset is the set of characters to trim around a message ID
const idTrimCutset string = "<> \n"

// parseAddresses parses a list of email addresses. Note that
// net/mail.Header[param] gets a list of addresses rather than slice.
func (p *Parser) parseAddresses(s string) ([]*mail.Address, error) {
	if s == "" {
		return nil, errorEmptyAddress
	}
	addresses := []*mail.Address{}
	decodedHeader, err := decoders.DecodeHeader(s)
	if err != nil {
		return addresses, fmt.Errorf("cannot decode address %q: %w", s, err)
	}
	// plug point for custom address parsing
	return p.addressesFunc(decodedHeader)
}

// parseAddress parses a single *mail.Address from a string using
// parseAddresses
func (p *Parser) parseAddress(s string) (*mail.Address, error) {
	if s == "" {
		return nil, errorEmptyAddress
	}
	decodedHeader, err := decoders.DecodeHeader(s)
	if err != nil {
		return nil, fmt.Errorf("cannot decode address %q: %w", s, err)
	}
	// plug point for custom address parsing
	return p.addressFunc(decodedHeader)
}

// parseHeaders parses the headers in the net/mail.Header at p.msg into
// p.email.Headers field values.
func (p *Parser) parseHeaders() error {

	// get is a shortcut to net/mail.Header.Get, which returns the first
	// value (if any) for a header field. Note that all lists of email
	// addresses are returned as single string, so should be retrieved
	// using "Get" rather than by map lookup.
	get := func(field string) string {
		return p.msg.Header.Get(field)
	}

	// getAll is shortcut to get the net/mail.Header []string elements
	getAll := func(field string) []string {
		return p.msg.Header[field]
	}

	// getID returns a cleaned message id
	getID := func(s string) string { return strings.Trim(s, idTrimCutset) }

	// getIDs returns a slice of cleaned message ids
	getIDs := func(s string) []string {
		ids := []string{}
		for _, id := range strings.Split(s, " ") {
			id := strings.TrimSpace(strings.Trim(id, idTrimCutset))
			if id == "" {
				continue
			}
			ids = append(ids, id)
		}
		return ids
	}

	callDateFunc := func(s string) (time.Time, error) {
		if s == "" {
			return time.Time{}, errorEmptyDate
		}
		// plug point for custom address parsing
		return p.dateFunc(s)
	}

	// getDecodedString decodes and trims a string header
	getDecodedString := func(s string) (string, error) {
		return decoders.DecodeHeader(strings.TrimSpace(s))
	}

	// getCSV gets parts of a comma delimited string
	getCSV := func(s string) []string {
		o := []string{}
		parts := strings.Split(s, ",")
		for _, p := range parts {
			pp := strings.TrimSpace(p)
			if len(pp) > 0 {
				o = append(o, pp)
			}
		}
		return o
	}

	// alias headers for easy reference
	h := &p.email.Headers

	// set contentInfo from parser
	h.ContentInfo = p.contentInfo

	h.ExtraHeaders = map[string][]string{}
	for key, value := range p.msg.Header {
		if isExplicitHeader(key) {
			continue
		}
		h.ExtraHeaders[key] = []string{}
		for _, val := range value {
			val, _ := decoders.DecodeHeader(val)
			h.ExtraHeaders[key] = append(h.ExtraHeaders[key], val)
		}
	}

	var err error
	if h.Sender, err = p.parseAddress(get("Sender")); err != nil {
		if !errors.Is(errorEmptyAddress, err) {
			return fmt.Errorf("cannot parse Sender header: %w", err)
		}
	}

	// Get email address lists via get. See get function comments.
	if h.From, err = p.parseAddresses(get("From")); err != nil {
		if !errors.Is(errorEmptyAddress, err) {
			return fmt.Errorf("cannot parse From header: %w", err)
		}
	}

	if h.ReplyTo, err = p.parseAddresses(get("Reply-To")); err != nil {
		if !errors.Is(errorEmptyAddress, err) {
			return fmt.Errorf("cannot parse Reply-To header: %w", err)
		}
	}

	if h.To, err = p.parseAddresses(get("To")); err != nil {
		if !errors.Is(errorEmptyAddress, err) {
			return fmt.Errorf("cannot parse To header: %w", err)
		}
	}

	if h.Cc, err = p.parseAddresses(get("Cc")); err != nil {
		if !errors.Is(errorEmptyAddress, err) {
			return fmt.Errorf("cannot parse Cc header: %w", err)
		}
	}

	if h.Bcc, err = p.parseAddresses(get("Bcc")); err != nil {
		if !errors.Is(errorEmptyAddress, err) {
			return fmt.Errorf("cannot parse Bcc header: %w", err)
		}
	}

	if h.ResentFrom, err = p.parseAddresses(get("Resent-From")); err != nil {
		if !errors.Is(errorEmptyAddress, err) {
			return fmt.Errorf("cannot parse Resent-From header: %w", err)
		}
	}

	if h.ResentSender, err = p.parseAddress(get("Resent-Sender")); err != nil {
		if !errors.Is(errorEmptyAddress, err) {
			return fmt.Errorf("cannot parse Resent-Sender header: %w", err)
		}
	}

	if h.ResentTo, err = p.parseAddresses(get("Resent-To")); err != nil {
		if !errors.Is(errorEmptyAddress, err) {
			return fmt.Errorf("cannot parse Resent-To header: %w", err)
		}
	}

	if h.ResentCc, err = p.parseAddresses(get("Resent-Cc")); err != nil {
		if !errors.Is(errorEmptyAddress, err) {
			return fmt.Errorf("cannot parse Resent-Cc header: %w", err)
		}
	}

	if h.ResentBcc, err = p.parseAddresses(get("Resent-Bcc")); err != nil {
		if !errors.Is(errorEmptyAddress, err) {
			return fmt.Errorf("cannot parse Resent-Bcc header: %w", err)
		}
	}

	if h.Date, err = callDateFunc(get("Date")); err != nil {
		if !errors.Is(errorEmptyDate, err) {
			return fmt.Errorf("cannot parse Date header: %w", err)
		}
	}

	if h.ResentDate, err = callDateFunc(get("Resent-Date")); err != nil {
		if !errors.Is(errorEmptyDate, err) {
			return fmt.Errorf("cannot parse Resent-Date header: %w", err)
		}
	}

	if h.Subject, err = getDecodedString(get("Subject")); err != nil {
		return fmt.Errorf("cannot parse Subject header: %w", err)
	}

	if h.Comments, err = getDecodedString(get("Comments")); err != nil {
		return fmt.Errorf("cannot parse Comments header: %w", err)
	}

	// consider parsing this into []Received
	if re := getAll("Received"); len(re) > 0 {
		h.Received = re
	}

	if id := getID(get("Message-ID")); id != "" {
		h.MessageID = id
	}

	if ids := getIDs(get("In-Reply-To")); len(ids) > 0 {
		h.InReplyTo = ids
	}

	if ids := getIDs(get("References")); len(ids) > 0 {
		h.References = ids
	}

	if kw := getCSV(get("Keywords")); len(kw) > 0 {
		h.Keywords = kw
	}

	if id := getID(get("Resent-Message-ID")); id != "" {
		h.ResentMessageID = id
	}

	return nil
}
