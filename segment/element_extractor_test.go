package segment

import (
	"reflect"
	"testing"
)

func Test_ExtractElements(t *testing.T) {
	type testData struct {
		in  string
		out []string
		err error
	}
	tests := []testData{
		{
			"abcde+123+012+de+'",
			[]string{
				"abcde",
				"123",
				"012",
				"de",
			},
			nil,
		},
		{
			"abcde:123:012+de+'",
			[]string{
				"abcde:123:012",
				"de",
			},
			nil,
		},
	}

	for _, test := range tests {
		extracted, err := ExtractElements([]byte(test.in))

		if err != nil {
			t.Logf("Expected no error, got %T:%v\n", err, err)
			t.Fail()
		}

		if extracted == nil {
			t.Logf("Expected result not to be nil\n")
			t.Fail()
		}

		actual := make([]string, len(extracted))
		for i, b := range extracted {
			actual[i] = string(b)
		}

		if !reflect.DeepEqual(test.out, actual) {
			t.Logf("Extract: \n%q\n", extracted)
			t.Logf("Expected result to equal\n%q\n\tgot\n%q\n", test.out, actual)
			t.Fail()
		}

	}
}
