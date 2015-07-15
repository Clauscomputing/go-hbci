package segment

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/element"
)

type SegmentSequence []Segment

var validHBCIVersions = []int{201, 210, 220}

func NewReferencingMessageHeaderSegment(size int, hbciVersion int, dialogId string, number int, referencedMessage *element.ReferenceMessage) *MessageHeaderSegment {
	m := NewMessageHeaderSegment(size, hbciVersion, dialogId, number)
	m.Ref = referencedMessage
	return m
}

func NewMessageHeaderSegment(size int, hbciVersion int, dialogId string, number int) *MessageHeaderSegment {
	m := &MessageHeaderSegment{
		Size:        element.NewDigit(size, 12),
		HBCIVersion: element.NewNumber(hbciVersion, 3),
		DialogID:    element.NewIdentification(dialogId),
		Number:      element.NewNumber(number, 4),
	}
	m.Segment = NewBasicSegment(1, m)
	return m
}

type MessageHeaderSegment struct {
	Segment
	Size        *element.DigitDataElement
	HBCIVersion *element.NumberDataElement
	DialogID    *element.IdentificationDataElement
	Number      *element.NumberDataElement
	Ref         *element.ReferenceMessage
}

func (m *MessageHeaderSegment) init() {
	*m.Size = *new(element.DigitDataElement)
	*m.HBCIVersion = *new(element.NumberDataElement)
	*m.DialogID = *new(element.IdentificationDataElement)
	*m.Number = *new(element.NumberDataElement)
	*m.Ref = *new(element.ReferenceMessage)
}
func (m *MessageHeaderSegment) version() int         { return 3 }
func (m *MessageHeaderSegment) id() string           { return "HNHBK" }
func (m *MessageHeaderSegment) referencedId() string { return "" }
func (m *MessageHeaderSegment) sender() string       { return senderBoth }

func (m *MessageHeaderSegment) UnmarshalHBCI(value []byte) error {
	value = bytes.TrimSuffix(value, []byte("'"))
	elements := bytes.Split(value, []byte("+"))
	elementsLen := len(elements)
	if elementsLen == 0 || elementsLen < 5 {
		return fmt.Errorf("Malformed marshaled value")
	}
	header := elements[0]
	numStr := bytes.Split(header, []byte(":"))[1]
	num, err := strconv.Atoi(string(numStr))
	if err != nil {
		return fmt.Errorf("Malformed segment header")
	}
	m.Segment = NewBasicSegment(num, m)
	size, err := strconv.Atoi(string(elements[1]))
	if err != nil {
		return fmt.Errorf("Error while unmarshaling size: %v", err)
	}
	hbciVersion, err := strconv.Atoi(string(elements[2]))
	if err != nil {
		return fmt.Errorf("Error while unmarshaling hbci version: %v", err)
	}
	dialogId := string(elements[3])
	messageNum, err := strconv.Atoi(string(elements[4]))
	if err != nil {
		return fmt.Errorf("Error while unmarshaling message number: %v", err)
	}
	m.Size = element.NewDigit(size, 12)
	m.HBCIVersion = element.NewNumber(hbciVersion, 3)
	m.DialogID = element.NewIdentification(dialogId)
	m.Number = element.NewNumber(messageNum, 4)
	return nil
}

func (m *MessageHeaderSegment) elements() []element.DataElement {
	return []element.DataElement{
		m.Size,
		m.HBCIVersion,
		m.DialogID,
		m.Number,
		m.Ref,
	}
}

func (m *MessageHeaderSegment) SetSize(size int) {
	*m.Size = *element.NewDigit(size, 12)
}

func NewMessageEndSegment(segmentNumber, messageNumber int) *MessageEndSegment {
	end := &MessageEndSegment{
		Number: element.NewNumber(messageNumber, 4),
	}
	end.Segment = NewBasicSegment(segmentNumber, end)
	return end
}

type MessageEndSegment struct {
	Segment
	Number *element.NumberDataElement
}

func (m *MessageEndSegment) init() {
	*m.Number = *new(element.NumberDataElement)
}
func (m *MessageEndSegment) version() int         { return 1 }
func (m *MessageEndSegment) id() string           { return "HNHBS" }
func (m *MessageEndSegment) referencedId() string { return "" }
func (m *MessageEndSegment) sender() string       { return senderBoth }

func (m *MessageEndSegment) elements() []element.DataElement {
	return []element.DataElement{
		m.Number,
	}
}
