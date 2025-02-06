package parser

import (
	"fmt"
	"net/mail"
	"strings"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/mnako/letters/email"

	"github.com/google/go-cmp/cmp"
)

func TestIsExplicitHeader(t *testing.T) {
	tests := []struct {
		name             string
		isExplicitHeader bool
	}{
		{"From", true},
		{"X-Received", false},
		{"Donald", false},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			got := isExplicitHeader(tt.name)
			if got, want := got, tt.isExplicitHeader; got != want {
				t.Errorf("got %t want %t for %s", got, want, tt.name)
			}
		})
	}
}

func TestParseAddressSingle(t *testing.T) {
	p := NewParser()
	tests := []struct {
		stringAddress string
		name          string
		address       string
	}{
		{"<test@test.com>", "", "test@test.com"},
		{"Ronny Burke <ronnie@example.com>", "Ronny Burke", "ronnie@example.com"},
		{`"=?UTF-8?B?SHViZXJ0IFNjaMO2bG5hc3Q=?=" <localpart@domain.tld>`, "Hubert Schölnast", "localpart@domain.tld"},
		{`=?ISO-8859-2?Q?"Odbieraj=B1ca,_Karolina"?= <karolina.odbierajaca@example.net>`, "Odbierająca, Karolina", "karolina.odbierajaca@example.net"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			addr, err := p.parseAddress(tt.stringAddress)
			if err != nil {
				t.Fatal(err)
			}
			if got, want := addr.Name, tt.name; got != want {
				t.Errorf("got %s want %s", got, want)
			}
			if got, want := addr.Address, tt.address; got != want {
				t.Errorf("got %s want %s", got, want)
			}
		})
	}
}

func TestParseAddresses(t *testing.T) {
	p := NewParser()
	tests := []struct {
		stringAddress string
		name          []string
		address       []string
	}{
		{
			stringAddress: `=?UTF-8?q?Bob_Recipient?= <bob.recipient@example.com>, =?UTf-8?q?Carol_Recipient?= <carol.recipient@example.com>`,
			name:          []string{"Bob Recipient", "Carol Recipient"},
			address:       []string{"bob.recipient@example.com", "carol.recipient@example.com"},
		},
		{
			stringAddress: `=?Iso-8859-2?q?"Odbieraj=B1cy,_Bob"?= <bob.odbierajacy@example.net>, =?ISO-8859-2?Q?"Odbieraj=B1ca,_Karolina"?= <karolina.odbierajaca@example.net>`,
			name:          []string{"Odbierający, Bob", "Odbierająca, Karolina"},
			address:       []string{"bob.odbierajacy@example.net", "karolina.odbierajaca@example.net"},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			addr, err := p.parseAddresses(tt.stringAddress)
			if err != nil {
				t.Fatal(err)
			}
			for ii, name := range tt.name {
				if got, want := addr[ii].Name, name; got != want {
					t.Errorf("got %s want %s", got, want)
				}

			}
			for ii, address := range tt.address {
				if got, want := addr[ii].Address, address; got != want {
					t.Errorf("got %s want %s", got, want)
				}
			}
		})
	}
}

func TestParseAddressCustomFunc(t *testing.T) {
	p := NewParser()
	b, a := "Bart Simpson", "<bart@example.com>"
	p.addressFunc = func(s string) (*mail.Address, error) {
		return &mail.Address{Name: b, Address: a}, nil
	}
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

func TestParseAddressesCustomFunc(t *testing.T) {
	p := NewParser()
	addresses := [][]string{
		[]string{"Bart Simpson", "<bart@example.com>"},
		[]string{"Darth Vader", "<darth@example.com>"},
	}
	p.addressesFunc = func(stringList string) ([]*mail.Address, error) {
		return []*mail.Address{
			&mail.Address{Name: addresses[0][0], Address: addresses[0][1]},
			&mail.Address{Name: addresses[1][0], Address: addresses[1][1]},
		}, nil
	}
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

func TestParseDateCustomFunc(t *testing.T) {
	p := NewParser()
	p.dateFunc = func(string) (time.Time, error) {
		return time.Time{}, nil
	}
	ti, err := p.dateFunc("whatever")
	if err != nil {
		t.Fatal(err)
	}
	if !ti.IsZero() {
		t.Errorf("expected %v to be time.Zero", ti)
	}
}

func TestParseHeaders(t *testing.T) {

	rawEmail := `Date: Mon, 01 Apr 2019 07:55:00 +0100 (BST)
From: Alice Sender <alice.sender@example.com>,
 Alice Sender <alice.sender@example.net>
Sender: Alice Sender <alice.sender@example.com>
Reply-To: Alice Sender <alice.sender@example.net>
To: Bob Recipient <bob.recipient@example.com>,
 Carol Recipient <carol.recipient@example.com>
Cc: Dan Recipient <dan.recipient@example.com>,
 Eve Recipient <eve.recipient@example.com>
Bcc: Frank Recipient <frank.recipient@example.com>,
 Grace Recipient <grace.recipient@example.com>
Subject: =?UTF-8?B?8J+Tpw==?= Test
 English Pangrams
In-Reply-To: <Message-Id-0@example.com>
References: <Message-Id-0@example.com>
Message-ID: <Message-Id-1@example.com>
Comments: Message Header Comment
Keywords: Keyword 1, Keyword 2
Resent-Date: Mon, 01 Apr 2019 07:55:00 +0100 (BST)
Resent-From: Alice Sender <alice.sender@example.net>,
 Alice Sender <alice.sender@example.com>
Resent-Sender: Alice Sender <alice.sender@example.net>
Resent-To: Bob Recipient <bob.recipient@example.net>,
 Carol Recipient <carol.recipient@example.net>
Resent-Cc: Dan Recipient <dan.recipient@example.net>,
 Eve Recipient <eve.recipient@example.net>
Resent-Bcc: Frank Recipient <frank.recipient@example.net>,
 Grace Recipient <grace.recipient@example.net>
Resent-Message-ID: <Message-Id-1@example.net>
Content-Type: MULtiPARt/mIXed; Charset="ascII"; bouNDARY="MixedBoundaryString"
Content-Transfer-Encoding: 7BiT
Delivery-date: Tue, 26 May 2020 12:01:38 +0000
Received: from securemail-y17.example.com ([196.35.198.77])
        by anotherexample.net with esmtps
        (TLS1.2:ECDHE_RSA_AES_256_GCM_SHA384:256)
        (envelope-from <amazing@examaple.com>)
        id 1jdYH3-00057X-TF
        for user@anotherexample.net; Mon, 01 Apr 2019 12:01:38 +0000
Received: from [10.1.1.1] (helo=[192.168.0.1])
        by securemail-pl-omx12.eample.com with esmtpa 
        (envelope-from <amazing@example.com>)
        id 1jdYGW-000aQH-Lx; Mon, 01 Apr 2019 14:01:05 +0200

`
	var debug bool = false

	var err error
	p := NewParser()
	p.msg, err = mail.ReadMessage(strings.NewReader(rawEmail))
	if err != nil {
		t.Fatal(err)
	}

	err = p.parseHeaders()
	if err != nil {
		t.Fatal(err)
	}

	if debug {
		spewer := spew.ConfigState{
			DisableCapacities:       true,
			DisablePointerAddresses: true,
			Indent:                  "\t",
		}
		spewer.Dump(p.email.Headers)
	}

	timeLayout := "2006-01-02 15:04:05 -0700 MST"
	toTime := func(s string) time.Time {
		t, _ := time.Parse(timeLayout, s)
		return t
	}

	expectedEmail := email.Email{}

	expectedEmail.Headers = email.Headers{
		Date:   toTime("2019-04-01 07:55:00 +0100 BST"),
		Sender: &mail.Address{Name: "Alice Sender", Address: "alice.sender@example.com"},
		From: []*mail.Address{
			&mail.Address{Name: "Alice Sender", Address: "alice.sender@example.com"},
			&mail.Address{Name: "Alice Sender", Address: "alice.sender@example.net"},
		},
		ReplyTo: []*mail.Address{
			&mail.Address{Name: "Alice Sender", Address: "alice.sender@example.net"},
		},
		To: []*mail.Address{
			&mail.Address{Name: "Bob Recipient", Address: "bob.recipient@example.com"},
			&mail.Address{Name: "Carol Recipient", Address: "carol.recipient@example.com"},
		},
		Cc: []*mail.Address{
			&mail.Address{Name: "Dan Recipient", Address: "dan.recipient@example.com"},
			&mail.Address{Name: "Eve Recipient", Address: "eve.recipient@example.com"},
		},
		Bcc: []*mail.Address{
			&mail.Address{Name: "Frank Recipient", Address: "frank.recipient@example.com"},
			&mail.Address{Name: "Grace Recipient", Address: "grace.recipient@example.com"},
		},
		MessageID: "Message-Id-1@example.com",
		InReplyTo: []string{
			"Message-Id-0@example.com",
		},
		References: []string{
			"Message-Id-0@example.com",
		},
		Subject:    "📧 Test English Pangrams",
		Comments:   "Message Header Comment",
		Keywords:   []string{"Keyword 1", "Keyword 2"},
		ResentDate: toTime("2019-04-01 07:55:00 +0100 BST"),
		ResentFrom: []*mail.Address{
			&mail.Address{Name: "Alice Sender", Address: "alice.sender@example.net"},
			&mail.Address{Name: "Alice Sender", Address: "alice.sender@example.com"},
		},
		ResentSender: &mail.Address{Name: "Alice Sender", Address: "alice.sender@example.net"},
		ResentTo: []*mail.Address{
			&mail.Address{Name: "Bob Recipient", Address: "bob.recipient@example.net"},
			&mail.Address{Name: "Carol Recipient", Address: "carol.recipient@example.net"},
		},
		ResentCc: []*mail.Address{
			&mail.Address{Name: "Dan Recipient", Address: "dan.recipient@example.net"},
			&mail.Address{Name: "Eve Recipient", Address: "eve.recipient@example.net"},
		},
		ResentBcc: []*mail.Address{
			&mail.Address{Name: "Frank Recipient", Address: "frank.recipient@example.net"},
			&mail.Address{Name: "Grace Recipient", Address: "grace.recipient@example.net"},
		},
		ResentMessageID: "Message-Id-1@example.net",
		// newly introduced for tracing
		Received: []string{
			"from securemail-y17.example.com ([196.35.198.77]) by anotherexample.net with esmtps (TLS1.2:ECDHE_RSA_AES_256_GCM_SHA384:256) (envelope-from <amazing@examaple.com>) id 1jdYH3-00057X-TF for user@anotherexample.net; Mon, 01 Apr 2019 12:01:38 +0000",
			"from [10.1.1.1] (helo=[192.168.0.1]) by securemail-pl-omx12.eample.com with esmtpa (envelope-from <amazing@example.com>) id 1jdYGW-000aQH-Lx; Mon, 01 Apr 2019 14:01:05 +0200",
		},
		ExtraHeaders: map[string][]string{
			"Delivery-Date": {"Tue, 26 May 2020 12:01:38 +0000"},
		},
	}

	got, want := p.email.Headers, expectedEmail.Headers
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("emails are not equal\n%s", diff)
	}

}
