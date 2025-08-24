package letters_test

import (
	"testing"
	"time"

	"github.com/mnako/letters"
)

type dateParsingTestCase struct {
	name         string
	dateHeader   string
	expectedDate time.Time
}

func TestParseDateHeader(t *testing.T) {
	t.Parallel()

	utcMinus0600FixedLocation := time.FixedZone("UTC-0600", -6*60*60)
	utcMinus0330FixedLocation := time.FixedZone("UTC-0330", -3.5*60*60)
	edtLocation := time.FixedZone("EDT", -4*60*60)
	gmtLocation, _ := time.LoadLocation("GMT")
	utcPlus0200FixedLocation := time.FixedZone("UTC+0200", 2*60*60)

	testCases := []dateParsingTestCase{
		{
			name:       "RFC822 A.3.1.",
			dateHeader: "26 Aug 76 1429 GMT",
			expectedDate: time.Date(
				1976,
				time.August,
				26,
				14,
				29,
				0,
				0,
				time.UTC,
			),
		},
		{
			name:       "RFC822 A.3.1.",
			dateHeader: "26 Aug 76 1429 EDT",
			expectedDate: time.Date(
				1976,
				time.August,
				26,
				14,
				29,
				0,
				0,
				edtLocation,
			),
		},
		{
			name:       "RFC5322 A.1.1.",
			dateHeader: "Fri, 21 Nov 1997 09:55:06 -0600",
			expectedDate: time.Date(
				1997,
				time.November,
				21,
				9,
				55,
				6,
				0,
				utcMinus0600FixedLocation,
			),
		},
		{
			name:       "RFC5322 A.1.2.",
			dateHeader: "Tue, 1 Jul 2003 10:52:37 +0200",
			expectedDate: time.Date(
				2003,
				time.July,
				1,
				10,
				52,
				37,
				0,
				utcPlus0200FixedLocation,
			),
		},
		{
			name:       "RFC5322 A.1.3.",
			dateHeader: "Thu, 13 Feb 1969 23:32:54 -0330",
			expectedDate: time.Date(
				1969,
				time.February,
				13,
				23,
				32,
				54,
				0,
				utcMinus0330FixedLocation,
			),
		},
		{
			name:       "RFC5322 Appendix A.5.",
			dateHeader: "Thu, 13 Feb 1969 23:32 -0330 (Newfoundland Time)",
			// From RFC5322 Appendix A.1.3:
			// The above example is aesthetically displeasing, but perfectly legal.
			// Note particularly [...] the missing seconds in the time of the date field;
			expectedDate: time.Date(
				1969,
				time.February,
				13,
				23,
				32,
				0,
				0,
				utcMinus0330FixedLocation,
			),
		},
		{
			name:       "RFC5322 Appendix A.6.2. Obsolete Dates",
			dateHeader: "21 Nov 97 09:55:06 GMT",
			expectedDate: time.Date(
				1997,
				time.November,
				21,
				9,
				55,
				6,
				0,
				gmtLocation,
			),
		},
	}

	for _, tc := range testCases {
		testCase := tc

		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				parsedDate := letters.ParseDateHeader(testCase.dateHeader)

				if !parsedDate.Equal(testCase.expectedDate) {
					t.Errorf("dates are not equal")
					t.Errorf("Got  %#v %d", parsedDate, parsedDate.Unix())
					t.Errorf(
						"Want %#v %d",
						testCase.expectedDate,
						testCase.expectedDate.Unix(),
					)
				}
			},
		)
	}
}
