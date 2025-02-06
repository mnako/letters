package decoders

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mnako/letters/email"
	// "github.com/google/go-cmp/cmp/cmpopts"
)

func TestDecodeHeader(t *testing.T) {
	tests := []struct {
		header string
		want   string
	}{
		{
			header: "Some One <someone@example.com>",
			want:   "Some One <someone@example.com>",
		},
		{
			header: `"=?UTF-8?B?SHViZXJ0IFNjaMO2bG5hc3Q=?=" <localpart@domain.tld>`,
			want:   `"Hubert Schölnast" <localpart@domain.tld>`,
		},
		// https://nerderati.com/mime-encoded-words-in-email-headers/
		{
			header: `=?utf-8?Q?Andreas_Birkeb=C3=A6k?=`,
			want:   `Andreas Birkebæk`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			got, err := DecodeHeader(tt.header)
			if err != nil {
				t.Fatal(err)
			}
			if got, want := got, tt.want; got != want {
				t.Errorf("got %s want %s", got, want)
			}
		})
	}
}

func TestDecodeContent(t *testing.T) {
	tests := []struct {
		name           string
		encodingString string
		ci             *email.ContentInfo
		content        string
		want           []byte
	}{
		// test charset decoding
		{
			name:           "text/plain from test_thai_multipart_mixed_tis-620_over_base64.txt",
			encodingString: "tis-620",
			ci:             &email.ContentInfo{TransferEncoding: "base64"},
			content: `
4LvnucG52MnC7MrYtLvD0ODKw9Sw4MXUyKTYs6To0iChx+jSusPDtNK92afK0bXH7OC0w9GoqdK5
IKinvejSv9G5vtGyudLH1KrSodLDIM3C6NLF6dKnvMXSrcTl4KLouabo0rrVsdLjpMMg5MHottfN
4rfJ4qHDuOGq6Ker0bTO1rTO0bS06NIgy9G0zcDRwuDLwdfNuaHVzNLN0aqs0srRwiC7r9S60bXU
u8PQvsS11KGuodPLubTjqCC+2bSo0uPL6ajq0OYgqOvS5iC56NK/0afgzcLPCgq50sLK0aemwNGz
sewg4M6nvtS30aHJ7L3R6KcgvNnp4LLo0qvW6KfB1c3SqtW+4LvnuaW5otLCo8e0ILbZobXTw8eo
u6/UutG11KHSw6jRur/pzafI0sUgsNK5xdGhudLM1KHSpNizy63Up6nRtcOqrtIgrNK5ysHSuNQ`,
			want: []byte(`เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า หัดอภัยเหมือนกีฬาอัชฌาสัย ปฏิบัติประพฤติกฎกำหนดใจ พูดจาให้จ๊ะๆ จ๋าๆ น่าฟังเอยฯ

นายสังฆภัณฑ์ เฮงพิทักษ์ฝั่ง ผู้เฒ่าซึ่งมีอาชีพเป็นฅนขายฃวด ถูกตำรวจปฏิบัติการจับฟ้องศาล ฐานลักนาฬิกาคุณหญิงฉัตรชฎา ฌานสมาธิ`),
		},
		// base64 std encoding
		{
			name:           "inline-jpg-image-without-disposition.jpg from test_thai_multipart_mixed_tis-620_over_base64.txt",
			encodingString: "",
			ci:             &email.ContentInfo{TransferEncoding: "base64"},
			content: `
/9j/2wBDAAMCAgICAgMCAgIDAwMDBAYEBAQEBAgGBgUGCQgKCgkICQkKDA8MCgsOCwkJDRENDg8Q
EBEQCgwSExIQEw8QEBD/yQALCAABAAEBAREA/8wABgAQEAX/2gAIAQEAAD8A0s8g/9k=`,
			want: []byte{
				255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4, 6, 4, 4, 4, 4, 4, 8, 6,
				6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9, 13, 17, 13, 14, 15, 16, 16,
				17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8, 0, 1, 0, 1, 1, 1, 17, 0,
				255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207, 32, 255, 217,
			},
		},
		// base64 raw encoding
		{
			name:           "raw inline-jpg-image-without-disposition.jpg from test_thai_multipart_mixed_tis-620_over_base64.txt",
			encodingString: "",
			ci:             &email.ContentInfo{TransferEncoding: "base64"},
			content: `
/9j/2wBDAAMCAgICAgMCAgIDAwMDBAYEBAQEBAgGBgUGCQgKCgkICQkKDA8MCgsOCwkJDRENDg8Q
EBEQCgwSExIQEw8QEBD/yQALCAABAAEBAREA/8wABgAQEAX/2gAIAQEAAD8A0s8g/9k`, // removed = sign
			want: []byte{
				255, 216, 255, 219, 0, 67, 0, 3, 2, 2, 2, 2, 2, 3, 2, 2, 2, 3, 3, 3, 3, 4, 6, 4, 4, 4, 4, 4, 8, 6,
				6, 5, 6, 9, 8, 10, 10, 9, 8, 9, 9, 10, 12, 15, 12, 10, 11, 14, 11, 9, 9, 13, 17, 13, 14, 15, 16, 16,
				17, 16, 10, 12, 18, 19, 18, 16, 19, 15, 16, 16, 16, 255, 201, 0, 11, 8, 0, 1, 0, 1, 1, 1, 17, 0,
				255, 204, 0, 6, 0, 16, 16, 5, 255, 218, 0, 8, 1, 1, 0, 0, 63, 0, 210, 207, 32, 255, 217,
			},
		},
		// quoted-printable encoding
		{
			name:           "inline part from cats.eml",
			encodingString: "UTF-8",
			ci:             &email.ContentInfo{TransferEncoding: "quoted-printable"},
			content: `<div dir=3D"ltr"><div>Pictures of cats!</div><img src=3D"cid:ii_m6s9zyhs0" =
alt=3D"cat2.png" width=3D"50" height=3D"50">=C2=A0<img src=3D"cid:ii_m6s9zy=
ja1" alt=3D"cat3.jpg" width=3D"50" height=3D"50">=C2=A0<img src=3D"cid:ii_m=
6s9zyka2" alt=3D"cat1.jpg" width=3D"50" height=3D"50"><br></div>`,
			want: []byte{
				60, 100, 105, 118, 32, 100, 105, 114, 61, 34, 108, 116, 114, 34, 62, 60, 100, 105, 118, 62, 80, 105,
				99, 116, 117, 114, 101, 115, 32, 111, 102, 32, 99, 97, 116, 115, 33, 60, 47, 100, 105, 118, 62, 60,
				105, 109, 103, 32, 115, 114, 99, 61, 34, 99, 105, 100, 58, 105, 105, 95, 109, 54, 115, 57, 122, 121,
				104, 115, 48, 34, 32, 97, 108, 116, 61, 34, 99, 97, 116, 50, 46, 112, 110, 103, 34, 32, 119, 105,
				100, 116, 104, 61, 34, 53, 48, 34, 32, 104, 101, 105, 103, 104, 116, 61, 34, 53, 48, 34, 62, 194,
				160, 60, 105, 109, 103, 32, 115, 114, 99, 61, 34, 99, 105, 100, 58, 105, 105, 95, 109, 54, 115, 57,
				122, 121, 106, 97, 49, 34, 32, 97, 108, 116, 61, 34, 99, 97, 116, 51, 46, 106, 112, 103, 34, 32,
				119, 105, 100, 116, 104, 61, 34, 53, 48, 34, 32, 104, 101, 105, 103, 104, 116, 61, 34, 53, 48, 34,
				62, 194, 160, 60, 105, 109, 103, 32, 115, 114, 99, 61, 34, 99, 105, 100, 58, 105, 105, 95, 109, 54,
				115, 57, 122, 121, 107, 97, 50, 34, 32, 97, 108, 116, 61, 34, 99, 97, 116, 49, 46, 106, 112, 103,
				34, 32, 119, 105, 100, 116, 104, 61, 34, 53, 48, 34, 32, 104, 101, 105, 103, 104, 116, 61, 34, 53,
				48, 34, 62, 60, 98, 114, 62, 60, 47, 100, 105, 118, 62,
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("test_%s", tt.name), func(t *testing.T) {
			reader := strings.NewReader(tt.content)
			tt.ci.Charset = tt.encodingString
			tt.ci.ExtractEncoding()
			got, err := io.ReadAll(DecodeContent(reader, tt.ci))
			if err != nil {
				t.Fatal(err)
			}
			// fmt.Println(got)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDecodeContentBody(t *testing.T) {

	// source tests/test_chinese_plaintext_gb18030_over_base64.txt
	ci := &email.ContentInfo{TransferEncoding: "base64"}
	ci.Charset = "GB18030"
	// ci.ExtractEncoding()

	content := `
yq/K0sqryr/KqcrPo6zKyMqoo6zKxMqzyq7KqKGjCsrPyrHKscrKytDK08qooaMKyq7KsaOsysrK
rsqoysrK0KGjCsrHyrGjrMrKyqnKz8rKytChowrKz8rTysfKrsqoo6zK0cq4ysajrMq5ysfKrsqo
ysXKwKGjCsrPyrDKx8quyqjKrKOsysrKr8rSoaMKyq/K0sqqo6zKz8q5yszKw8qvytKhowrKr8rS
ysOjrMrPyrzK1MqzysfKrsqooaMKyrPKsaOsyrzKtsrHyq7KqMqso6zKtcquyq/KqMqsoaMKytTK
zcrHysKhow==`

	want := `石室诗士施氏，嗜狮，誓食十狮。
氏时时适市视狮。
十时，适十狮适市。
是时，适施氏适市。
氏视是十狮，恃矢势，使是十狮逝世。
氏拾是十狮尸，适石室。
石室湿，氏使侍拭石室。
石室拭，氏始试食是十狮。
食时，始识是十狮尸，实十石狮尸。
试释是事。`

	got, err := io.ReadAll(DecodeContent(strings.NewReader(content), ci))
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(got), want); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
	if got, want := (ci.Encoding == nil), false; got != want {
		t.Errorf("encoding should not be nil, got %t", got)
	}
}
