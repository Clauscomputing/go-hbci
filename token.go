package hbci

import "fmt"

// Token represents a token returned from the scanner.
type Token struct {
	typ TokenType // Type, such as FLOAT
	val string    // Value, such as "23.2".
	pos int       // position of token in input
}

func (t Token) String() string {
	switch t.typ {
	case EOF:
		return "EOF"
	case ERROR:
		return t.val
	}
	if len(t.val) > 10 {
		return fmt.Sprintf("%.10q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}

const (
	ERROR TokenType = iota // error occurred;
	// value is text of error
	DATA_ELEMENT                 // Datenelement (DE)
	DATA_ELEMENT_SEPARATOR       // Datenelement (DE)-Trennzeichen
	GROUP_DATA_ELEMENT           // Gruppendatenelement (GD)
	GROUP_DATA_ELEMENT_SEPARATOR // Gruppendatenelement (GD)-Trennzeichen
	SEGMENT                      // Segment
	SEGMENT_HEADER               // Segmentende-Zeichen
	SEGMENT_END_MARKER           // Segmentende-Zeichen
	ESCAPE_SEQUENCE              // Freigabeabfolge
	ESCAPE_CHARACTER             // Freigabezeichen
	ESCAPED_CHARACTER            // Freigegebenes Zeichen
	BINARY_DATA_LENGTH           // Binärdaten Länge
	BINARY_DATA                  // Binärdaten
	BINARY_DATA_MARKER           // Binärdatenkennzeichen
	ALPHA_NUMERIC                // an
	TEXT                         // txt
	DTAUS_CHARSET                // dta
	NUMERIC                      // num: 0-9 without leading 0
	DIGIT                        // dig: 0-9 with optional leading 0
	FLOAT                        // float
	YES_NO                       // jn
	DATE                         // dat
	VIRTUAL_DATE                 // vdat
	TIME                         // tim
	IDENTIFICATION               // id
	COUNTRY_CODE                 // ctr: ISO 3166-1 numeric
	CURRENCY                     // cur: ISO 4217
	VALUE                        // wrt
	EOF
)

// TokenType identifies the type of lex tokens.
type TokenType int

var tokenName = map[TokenType]string{
	ERROR: "error",
	// value is text of error
	DATA_ELEMENT:                 "dataElement",
	DATA_ELEMENT_SEPARATOR:       "dataElementSeparator",
	GROUP_DATA_ELEMENT:           "groupDataElement",
	GROUP_DATA_ELEMENT_SEPARATOR: "groupDataElementSeparator",
	SEGMENT:            "segment",
	SEGMENT_HEADER:     "segmentHeader",
	SEGMENT_END_MARKER: "segmentEndMarker",
	ESCAPE_SEQUENCE:    "escapeSequence",
	ESCAPE_CHARACTER:   "escapeCharacter",
	ESCAPED_CHARACTER:  "escapedCharacter",
	BINARY_DATA_LENGTH: "binaryDataLength",
	BINARY_DATA:        "binaryData",
	BINARY_DATA_MARKER: "binaryDataMarker",
	ALPHA_NUMERIC:      "alphaNumeric",
	TEXT:               "text",
	DTAUS_CHARSET:      "dtausCharset",
	NUMERIC:            "numeric",
	DIGIT:              "digit",
	FLOAT:              "float",
	YES_NO:             "yesNo",
	DATE:               "date",
	VIRTUAL_DATE:       "virtualDate",
	TIME:               "time",
	IDENTIFICATION:     "identification",
	COUNTRY_CODE:       "countryCode",
	CURRENCY:           "currency",
	VALUE:              "value",
	EOF:                "eof",
}

func (t TokenType) String() string {
	s := tokenName[t]
	if s == "" {
		return fmt.Sprintf("Token%d", int(t))
	}
	return s
}
