package segment

import (
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewCommunicationAccessRequestSegment(fromBank domain.BankID, toBank domain.BankID, maxEntries int, continuationReference string) *CommunicationAccessRequestSegment {
	c := &CommunicationAccessRequestSegment{
		FromBankID: element.NewBankIdentification(fromBank),
		ToBankID:   element.NewBankIdentification(toBank),
		MaxEntries: element.NewNumber(maxEntries, 4),
	}
	if continuationReference != "" {
		c.ContinuationReference = element.NewAlphaNumeric(continuationReference, 35)
	}
	c.ClientSegment = NewBasicSegment(2, c)
	return c
}

type CommunicationAccessRequestSegment struct {
	ClientSegment
	FromBankID            *element.BankIdentificationDataElement
	ToBankID              *element.BankIdentificationDataElement
	MaxEntries            *element.NumberDataElement
	ContinuationReference *element.AlphaNumericDataElement
}

func (c *CommunicationAccessRequestSegment) Version() int         { return 3 }
func (c *CommunicationAccessRequestSegment) ID() string           { return "HKKOM" }
func (c *CommunicationAccessRequestSegment) referencedId() string { return "" }
func (c *CommunicationAccessRequestSegment) sender() string       { return senderUser }

func (c *CommunicationAccessRequestSegment) elements() []element.DataElement {
	return []element.DataElement{
		c.FromBankID,
		c.ToBankID,
		c.MaxEntries,
		c.ContinuationReference,
	}
}

const HKKOMSegmentNumber = -1

func NewCommunicationAccessResponseSegment(bankId domain.BankID, language int, params domain.CommunicationParameter) *CommunicationAccessResponseSegment {
	c := &CommunicationAccessResponseSegment{
		BankID:              element.NewBankIdentification(bankId),
		StandardLanguage:    element.NewNumber(language, 3),
		CommunicationParams: element.NewCommunicationParameter(params),
	}
	header := element.NewReferencingSegmentHeader("HIKOM", 4, 3, HKKOMSegmentNumber)
	c.Segment = NewBasicSegmentWithHeader(header, c)
	return c
}

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment CommunicationAccessResponseSegment

type CommunicationAccessResponseSegment struct {
	Segment
	BankID              *element.BankIdentificationDataElement
	StandardLanguage    *element.NumberDataElement
	CommunicationParams *element.CommunicationParameterDataElement
}

func (c *CommunicationAccessResponseSegment) Version() int         { return 3 }
func (c *CommunicationAccessResponseSegment) ID() string           { return "HIKOM" }
func (c *CommunicationAccessResponseSegment) referencedId() string { return "HKKOM" }
func (c *CommunicationAccessResponseSegment) sender() string       { return senderBank }

func (c *CommunicationAccessResponseSegment) elements() []element.DataElement {
	return []element.DataElement{
		c.BankID,
		c.StandardLanguage,
		c.CommunicationParams,
	}
}
