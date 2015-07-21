package segment

import (
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

var HKVVBSegmentNumber = -1

func NewCommonBankParameterSegment(
	bpdVersion int,
	bankId domain.BankId,
	bankName string,
	businessTransactionCount int,
	supportedLanguages []int,
	supportedHBCIVersions []int,
	maxMessageSize int) *CommonBankParameterSegment {
	c := &CommonBankParameterSegment{
		BPDVersion:               element.NewNumber(bpdVersion, 3),
		BankID:                   element.NewBankIndentification(bankId),
		BankName:                 element.NewAlphaNumeric(bankName, 60),
		BusinessTransactionCount: element.NewNumber(businessTransactionCount, 3),
		SupportedLanguages:       element.NewSupportedLanguages(supportedLanguages...),
		SupportedHBCIVersions:    element.NewSupportedHBCIVersions(supportedHBCIVersions...),
		MaxMessageSize:           element.NewNumber(maxMessageSize, 4),
	}
	header := element.NewReferencingSegmentHeader("HIBPA", 1, 2, HKVVBSegmentNumber)
	c.Segment = NewBasicSegmentWithHeader(header, c)
	return c
}

type CommonBankParameterSegment struct {
	Segment
	BPDVersion               *element.NumberDataElement
	BankID                   *element.BankIdentificationDataElement
	BankName                 *element.AlphaNumericDataElement
	BusinessTransactionCount *element.NumberDataElement
	SupportedLanguages       *element.SupportedLanguagesDataElement
	SupportedHBCIVersions    *element.SupportedHBCIVersionsDataElement
	MaxMessageSize           *element.NumberDataElement
}

func (c *CommonBankParameterSegment) version() int         { return 2 }
func (c *CommonBankParameterSegment) id() string           { return "HIBPA" }
func (c *CommonBankParameterSegment) referencedId() string { return "HKVVB" }
func (c *CommonBankParameterSegment) sender() string       { return senderBank }

func (c *CommonBankParameterSegment) elements() []element.DataElement {
	return []element.DataElement{
		c.BPDVersion,
		c.BankID,
		c.BankName,
		c.BusinessTransactionCount,
		c.SupportedLanguages,
		c.SupportedHBCIVersions,
		c.MaxMessageSize,
	}
}

func (c *CommonBankParameterSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 || len(elements) < 7 {
		return fmt.Errorf("Malformed marshaled value")
	}
	segment, err := SegmentFromHeaderBytes(elements[0], c)
	if err != nil {
		return err
	}
	c.Segment = segment
	version, err := strconv.Atoi(string(elements[1]))
	if err != nil {
		return err
	}
	c.BPDVersion = element.NewNumber(version, 3)
	bankId := &element.BankIdentificationDataElement{}
	err = bankId.UnmarshalHBCI(elements[2])
	if err != nil {
		return err
	}
	c.BankID = bankId
	c.BankName = element.NewAlphaNumeric(string(elements[3]), 60)
	transactionCount, err := strconv.Atoi(string(elements[4]))
	if err != nil {
		return err
	}
	c.BusinessTransactionCount = element.NewNumber(transactionCount, 3)
	languages := &element.SupportedLanguagesDataElement{}
	err = languages.UnmarshalHBCI(elements[5])
	if err != nil {
		return err
	}
	c.SupportedLanguages = languages
	versions := &element.SupportedHBCIVersionsDataElement{}
	err = versions.UnmarshalHBCI(elements[6])
	if err != nil {
		return err
	}
	c.SupportedHBCIVersions = versions
	if len(elements) == 8 {
		maxSize, err := strconv.Atoi(string(elements[7]))
		if err != nil {
			return err
		}
		c.MaxMessageSize = element.NewNumber(maxSize, 4)
	}
	return nil
}

func (c *CommonBankParameterSegment) BankParameterData() domain.BankParameterData {
	return domain.BankParameterData{
		Version:                   c.BPDVersion.Val(),
		BankID:                    c.BankID.Val(),
		BankName:                  c.BankName.Val(),
		MaxTransactionsPerMessage: c.BusinessTransactionCount.Val(),
	}
}

type SecurityMethodSegment struct {
	Segment
	MixAllowed       *element.BooleanDataElement
	SupportedMethods *element.SupportedSecurityMethodDataElement
}

func (s *SecurityMethodSegment) version() int         { return 2 }
func (s *SecurityMethodSegment) id() string           { return "HISHV" }
func (s *SecurityMethodSegment) referencedId() string { return "HKVVB" }
func (s *SecurityMethodSegment) sender() string       { return senderBank }

func (s *SecurityMethodSegment) elements() []element.DataElement {
	return []element.DataElement{
		s.MixAllowed,
		s.SupportedMethods,
	}
}

type CompressionMethodSegment struct {
	Segment
	SupportedCompressionMethods *element.SupportedCompressionMethodsDataElement
}

func (c *CompressionMethodSegment) version() int         { return 1 }
func (c *CompressionMethodSegment) id() string           { return "HIKPV" }
func (c *CompressionMethodSegment) referencedId() string { return "HKVVB" }
func (c *CompressionMethodSegment) sender() string       { return senderBank }

func (c *CompressionMethodSegment) elements() []element.DataElement {
	return []element.DataElement{
		c.SupportedCompressionMethods,
	}
}

type BusinessTransactionParamsSegment struct {
	Segment
	ID            string
	Version       int
	MaxJobs       *element.NumberDataElement
	MinSignatures *element.NumberDataElement
	Params        *element.BusinessTransactionParameter
}

func (b *BusinessTransactionParamsSegment) version() int         { return b.Version }
func (b *BusinessTransactionParamsSegment) id() string           { return b.ID }
func (b *BusinessTransactionParamsSegment) referencedId() string { return "HKVVB" }
func (b *BusinessTransactionParamsSegment) sender() string       { return senderBank }

func (b *BusinessTransactionParamsSegment) elements() []element.DataElement {
	return []element.DataElement{
		b.MaxJobs,
		b.MinSignatures,
		b.Params,
	}
}
