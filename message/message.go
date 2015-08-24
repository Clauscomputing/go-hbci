package message

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/segment"
)

var bankSegments = map[string]segment.Segment{
	"HIRMG": &segment.MessageAcknowledgement{},
}

type Message interface {
	MessageHeader() *segment.MessageHeaderSegment
	MessageEnd() *segment.MessageEndSegment
	FindSegment(segmentID string) []byte
	FindSegments(segmentID string) [][]byte
	SegmentNumber(segmentID string) int
}

type ClientMessage interface {
	Message
	MarshalHBCI() ([]byte, error)
	Encrypt(provider CryptoProvider) (*EncryptedMessage, error)
	SetMessageNumber(messageNumber int)
}

type BankMessage interface {
	Message
	Acknowledgements() []domain.Acknowledgement
}

type HBCIMessage interface {
	HBCIVersion() segment.HBCIVersion
	HBCISegments() []segment.ClientSegment
}

type SignedHBCIMessage interface {
	HBCIMessage
	SetNumbers()
	SetSignatureHeader(*segment.SignatureHeaderSegment)
	SetSignatureEnd(*segment.SignatureEndSegment)
}

func NewHBCIMessage(hbciVersion segment.HBCIVersion, segments ...segment.ClientSegment) HBCIMessage {
	return &hbciMessage{hbciSegments: segments, hbciVersion: hbciVersion}
}

type hbciMessage struct {
	hbciSegments []segment.ClientSegment
	hbciVersion  segment.HBCIVersion
}

func (h *hbciMessage) HBCIVersion() segment.HBCIVersion {
	return h.hbciVersion
}

func (h *hbciMessage) HBCISegments() []segment.ClientSegment {
	return h.hbciSegments
}

func NewHBCIClientMessage(segments ...segment.ClientSegment) *BasicClientMessage {
	return NewBasicClientMessage(hbciSegmentClientMessage(segments))
}

type hbciSegmentClientMessage []segment.ClientSegment

func (h hbciSegmentClientMessage) jobs() []segment.ClientSegment {
	return h
}

func NewBasicMessageWithHeaderAndEnd(header *segment.MessageHeaderSegment, end *segment.MessageEndSegment, message HBCIMessage) *BasicMessage {
	b := &BasicMessage{
		Header:      header,
		End:         end,
		HBCIMessage: message,
		hbciVersion: message.HBCIVersion(),
	}
	return b
}

func NewBasicMessage(message HBCIMessage) *BasicMessage {
	b := &BasicMessage{
		HBCIMessage: message,
		hbciVersion: message.HBCIVersion(),
	}
	return b
}

type BasicMessage struct {
	Header         *segment.MessageHeaderSegment
	End            *segment.MessageEndSegment
	SignatureBegin *segment.SignatureHeaderSegment
	SignatureEnd   *segment.SignatureEndSegment
	HBCIMessage
	hbciVersion      segment.HBCIVersion
	marshaledContent []byte
}

func (b *BasicMessage) SetNumbers() {
	if b.HBCIMessage == nil {
		panic(fmt.Errorf("HBCIMessage must be set"))
	}
	n := 0
	num := func() int {
		n += 1
		return n
	}
	b.Header.SetNumber(num)
	if b.SignatureBegin != nil {
		b.SignatureBegin.SetNumber(num)
	}
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if !reflect.ValueOf(segment).IsNil() {
			segment.SetNumber(num)
		}
	}
	if b.SignatureEnd != nil {
		b.SignatureEnd.SetNumber(num)
	}
	b.End.SetNumber(num)
}

func (b *BasicMessage) SetSize() error {
	if b.HBCIMessage == nil {
		return fmt.Errorf("HBCIMessage must be set")
	}
	var buffer bytes.Buffer
	headerBytes, err := b.Header.MarshalHBCI()
	if err != nil {
		return err
	}
	buffer.Write(headerBytes)
	if b.SignatureBegin != nil {
		sigBytes, err := b.SignatureBegin.MarshalHBCI()
		if err != nil {
			return err
		}
		buffer.Write(sigBytes)
	}
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if !reflect.ValueOf(segment).IsNil() {
			segBytes, err := segment.MarshalHBCI()
			if err != nil {
				return err
			}
			buffer.Write(segBytes)
		}
	}
	if b.SignatureEnd != nil {
		sigEndBytes, err := b.SignatureEnd.MarshalHBCI()
		if err != nil {
			return err
		}
		buffer.Write(sigEndBytes)
	}
	endBytes, err := b.End.MarshalHBCI()
	if err != nil {
		return err
	}
	buffer.Write(endBytes)
	b.Header.SetSize(buffer.Len())
	return nil
}

func (b *BasicMessage) SetMessageNumber(messageNumber int) {
	b.Header.SetMessageNumber(messageNumber)
}

func (b *BasicMessage) MarshalHBCI() ([]byte, error) {
	if b.HBCIMessage == nil {
		return nil, fmt.Errorf("HBCIMessage must be set")
	}
	err := b.SetSize()
	if err != nil {
		return nil, err
	}
	if len(b.marshaledContent) == 0 {
		var buffer bytes.Buffer
		headerBytes, err := b.Header.MarshalHBCI()
		if err != nil {
			return nil, err
		}
		buffer.Write(headerBytes)
		if b.SignatureBegin != nil {
			sigBytes, err := b.SignatureBegin.MarshalHBCI()
			if err != nil {
				return nil, err
			}
			buffer.Write(sigBytes)
		}
		for _, segment := range b.HBCIMessage.HBCISegments() {
			if !reflect.ValueOf(segment).IsNil() {
				segBytes, err := segment.MarshalHBCI()
				if err != nil {
					return nil, err
				}
				buffer.Write(segBytes)
			}
		}
		if b.SignatureEnd != nil {
			sigEndBytes, err := b.SignatureEnd.MarshalHBCI()
			if err != nil {
				return nil, err
			}
			buffer.Write(sigEndBytes)
		}
		endBytes, err := b.End.MarshalHBCI()
		if err != nil {
			return nil, err
		}
		buffer.Write(endBytes)
		b.marshaledContent = buffer.Bytes()
	}
	return b.marshaledContent, nil
}

func (b *BasicMessage) Sign(provider SignatureProvider) (*BasicSignedMessage, error) {
	if b.HBCIMessage == nil {
		panic(fmt.Errorf("HBCIMessage must be set"))
	}
	// TODO: fix only PinTan segments!!!
	b.SignatureBegin = b.hbciVersion.PinTanSignatureHeader("", "", domain.KeyName{})
	provider.WriteSignatureHeader(b.SignatureBegin)
	b.SignatureEnd = b.hbciVersion.SignatureEnd(-1, "")
	b.SetNumbers()
	var buffer bytes.Buffer
	buffer.WriteString(b.SignatureBegin.String())
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if !reflect.ValueOf(segment).IsNil() {
			buffer.WriteString(segment.String())
		}
	}
	sig, err := provider.Sign(buffer.Bytes())
	if err != nil {
		return nil, err
	}
	provider.WriteSignature(b.SignatureEnd, sig)
	signedMessage := NewBasicSignedMessage(b)
	return signedMessage, nil
}

func (b *BasicMessage) Encrypt(provider CryptoProvider) (*EncryptedMessage, error) {
	if b.HBCIMessage == nil {
		return nil, fmt.Errorf("HBCIMessage must be set")
	}
	var messageBytes []byte
	if b.SignatureBegin != nil {
		sigBytes, err := b.SignatureBegin.MarshalHBCI()
		if err != nil {
			return nil, err
		}
		messageBytes = append(messageBytes, sigBytes...)
	}
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if !reflect.ValueOf(segment).IsNil() {
			segBytes, err := segment.MarshalHBCI()
			if err != nil {
				return nil, err
			}
			messageBytes = append(messageBytes, segBytes...)
		}
	}
	if b.SignatureEnd != nil {
		sigEndBytes, err := b.SignatureEnd.MarshalHBCI()
		if err != nil {
			return nil, err
		}
		messageBytes = append(messageBytes, sigEndBytes...)
	}
	encryptedMessage, err := provider.Encrypt(messageBytes)
	if err != nil {
		return nil, err
	}
	encryptionMessage := NewEncryptedMessage(b.Header, b.End, b.hbciVersion)
	encryptionMessage.EncryptionHeader = b.hbciVersion.PinTanEncryptionHeader("", domain.KeyName{})
	provider.WriteEncryptionHeader(encryptionMessage.EncryptionHeader)
	encryptionMessage.EncryptedData = segment.NewEncryptedDataSegment(encryptedMessage)
	return encryptionMessage, nil
}

func (b *BasicMessage) MessageHeader() *segment.MessageHeaderSegment {
	return b.Header
}

func (b *BasicMessage) MessageEnd() *segment.MessageEndSegment {
	return b.End
}

func (b *BasicMessage) FindSegment(segmentID string) []byte {
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if segment.Header().ID.Val() == segmentID {
			return []byte(segment.String())
		}
	}
	return nil
}

func (b *BasicMessage) FindSegments(segmentID string) [][]byte {
	var segments [][]byte
	for _, segment := range b.HBCIMessage.HBCISegments() {
		if segment.Header().ID.Val() == segmentID {
			segments = append(segments, []byte(segment.String()))
		}
	}
	return segments
}

func (b *BasicMessage) SegmentNumber(segmentID string) int {
	idx := -1
	//for i, segment := range b.HBCIMessage.HBCISegments() {
	//if
	//}
	return idx
}

func NewBasicSignedMessage(message *BasicMessage) *BasicSignedMessage {
	b := &BasicSignedMessage{
		message: message,
	}
	return b
}

type BasicSignedMessage struct {
	message *BasicMessage
}

func (b *BasicSignedMessage) SetNumbers() {
	if b.message.SignatureBegin == nil || b.message.SignatureEnd == nil {
		panic(fmt.Errorf("Cannot call set Numbers when signature is not set"))
	}
	b.message.SetNumbers()
}

func (b *BasicSignedMessage) SetSignatureHeader(sigBegin *segment.SignatureHeaderSegment) {
	b.message.SignatureBegin = sigBegin
}

func (b *BasicSignedMessage) SetSignatureEnd(sigEnd *segment.SignatureEndSegment) {
	b.message.SignatureEnd = sigEnd
}

func (b *BasicSignedMessage) HBCIVersion() segment.HBCIVersion {
	return b.message.HBCIVersion()
}

func (b *BasicSignedMessage) HBCISegments() []segment.ClientSegment {
	return b.message.HBCISegments()
}

func (b *BasicSignedMessage) MarshalHBCI() ([]byte, error) {
	return b.message.MarshalHBCI()
}

func (b *BasicSignedMessage) Encrypt(provider CryptoProvider) (*EncryptedMessage, error) {
	return b.message.Encrypt(provider)
}

type bankMessage interface {
	dataSegments() []segment.Segment
}

type basicBankMessage struct {
	*BasicMessage
	bankMessage
	MessageAcknowledgements *segment.MessageAcknowledgement
	SegmentAcknowledgements *segment.SegmentAcknowledgement
}

type clientMessage interface {
	jobs() []segment.ClientSegment
}

func NewBasicClientMessage(clientMessage clientMessage) *BasicClientMessage {
	b := &BasicClientMessage{
		clientMessage: clientMessage,
	}
	b.BasicMessage = NewBasicMessage(b)
	return b
}

type BasicClientMessage struct {
	*BasicMessage
	clientMessage
}

func (b *BasicClientMessage) HBCISegments() []segment.ClientSegment {
	return b.clientMessage.jobs()
}
