package decoders

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mnako/letters/email"
	"golang.org/x/net/html/charset"
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
		name                    string
		encodingString          string
		contentTransferEncoding email.ContentTransferEncoding
		content                 string
		want                    []byte
	}{
		// test charset decoding
		{
			name:                    "text/plain from test_thai_multipart_mixed_tis-620_over_base64.txt",
			encodingString:          "tis-620",
			contentTransferEncoding: email.CTEBase64,
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
			name:                    "inline-jpg-image-without-disposition.jpg from test_thai_multipart_mixed_tis-620_over_base64.txt",
			encodingString:          "",
			contentTransferEncoding: email.CTEBase64,
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
			name:                    "raw inline-jpg-image-without-disposition.jpg from test_thai_multipart_mixed_tis-620_over_base64.txt",
			encodingString:          "",
			contentTransferEncoding: email.CTEBase64,
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
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("test_%s", tt.name), func(t *testing.T) {
			reader := strings.NewReader(tt.content)
			encoding, _ := charset.Lookup(tt.encodingString)
			cte := tt.contentTransferEncoding
			got, err := io.ReadAll(DecodeContent(reader, encoding, cte))
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
