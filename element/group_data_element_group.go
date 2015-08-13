package element

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/domain"
)

func NewAmount(value float64, currency string) *AmountDataElement {
	a := &AmountDataElement{
		Amount:   NewValue(value),
		Currency: NewCurrency(currency),
	}
	a.DataElement = NewGroupDataElementGroup(AmountGDEG, 2, a)
	return a
}

type AmountDataElement struct {
	DataElement
	Amount   *ValueDataElement
	Currency *CurrencyDataElement
}

func (a *AmountDataElement) Elements() []DataElement {
	return []DataElement{
		a.Amount,
		a.Currency,
	}
}

func (a *AmountDataElement) Val() domain.Amount {
	return domain.Amount{
		Amount:   a.Amount.Val(),
		Currency: a.Currency.Val(),
	}
}

func (a *AmountDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 2 {
		return fmt.Errorf("Malformed marshaled value")
	}
	a.Amount = &ValueDataElement{}
	err = a.Amount.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	a.Currency = &CurrencyDataElement{}
	err = a.Currency.UnmarshalHBCI(elements[1])
	if err != nil {
		return err
	}
	a.DataElement = NewGroupDataElementGroup(AmountGDEG, 2, a)
	return nil
}

func NewBankIndentification(bankId domain.BankId) *BankIdentificationDataElement {
	b := &BankIdentificationDataElement{
		CountryCode: NewCountryCode(bankId.CountryCode),
		BankID:      NewAlphaNumeric(bankId.ID, 30),
	}
	b.DataElement = NewGroupDataElementGroup(BankIdentificationGDEG, 2, b)
	return b
}

type BankIdentificationDataElement struct {
	DataElement
	CountryCode *CountryCodeDataElement
	BankID      *AlphaNumericDataElement
}

func (b *BankIdentificationDataElement) Val() domain.BankId {
	return domain.BankId{
		CountryCode: b.CountryCode.Val(),
		ID:          b.BankID.Val(),
	}
}

func (b *BankIdentificationDataElement) Elements() []DataElement {
	return []DataElement{
		b.CountryCode,
		b.BankID,
	}
}

func (b *BankIdentificationDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) < 2 {
		return fmt.Errorf("Malformed marshaled value")
	}
	countryCode := &CountryCodeDataElement{}
	err = countryCode.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	b.CountryCode = countryCode
	b.BankID = NewAlphaNumeric(charset.ToUtf8(elements[1]), 30)
	return nil
}

func NewAccountConnection(conn domain.AccountConnection) *AccountConnectionDataElement {
	a := &AccountConnectionDataElement{
		AccountId:                 NewIdentification(conn.AccountID),
		SubAccountCharacteristics: NewIdentification(conn.SubAccountCharacteristics),
		CountryCode:               NewCountryCode(conn.CountryCode),
		BankId:                    NewAlphaNumeric(conn.BankID, 30),
	}
	a.DataElement = NewGroupDataElementGroup(AccountConnectionGDEG, 4, a)
	return a
}

type AccountConnectionDataElement struct {
	DataElement
	AccountId                 *IdentificationDataElement
	SubAccountCharacteristics *IdentificationDataElement
	CountryCode               *CountryCodeDataElement
	BankId                    *AlphaNumericDataElement
}

func (a *AccountConnectionDataElement) Elements() []DataElement {
	return []DataElement{
		a.AccountId,
		a.SubAccountCharacteristics,
		a.CountryCode,
		a.BankId,
	}
}

func (a *AccountConnectionDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) < 4 {
		return fmt.Errorf("Malformed AccountConnection")
	}
	countryCode, err := strconv.Atoi(charset.ToUtf8(elements[2]))
	if err != nil {
		return fmt.Errorf("%T: Malformed CountryCode: %q", a, elements[2])
	}
	accountConnection := domain.AccountConnection{
		AccountID:                 charset.ToUtf8(elements[0]),
		SubAccountCharacteristics: charset.ToUtf8(elements[1]),
		CountryCode:               countryCode,
		BankID:                    charset.ToUtf8(elements[3]),
	}
	*a = *NewAccountConnection(accountConnection)
	return nil
}

func (a *AccountConnectionDataElement) Val() domain.AccountConnection {
	return domain.AccountConnection{
		AccountID:                 a.AccountId.Val(),
		SubAccountCharacteristics: a.SubAccountCharacteristics.Val(),
		CountryCode:               a.CountryCode.Val(),
		BankID:                    a.BankId.Val(),
	}
}

func NewBalance(amount domain.Amount, date time.Time, withTime bool) *BalanceDataElement {
	var debitCredit string
	if amount.Amount < 0 {
		debitCredit = "D"
	} else {
		debitCredit = "C"
	}
	b := &BalanceDataElement{
		DebitCreditIndicator: NewAlphaNumeric(debitCredit, 1),
		Amount:               NewValue(math.Abs(amount.Amount)),
		Currency:             NewCurrency(amount.Currency),
		TransmissionDate:     NewDate(date),
	}
	if withTime {
		b.TransmissionTime = NewTime(date)
	}
	b.DataElement = NewGroupDataElementGroup(BalanceGDEG, 5, b)
	return b
}

type BalanceDataElement struct {
	DataElement
	DebitCreditIndicator *AlphaNumericDataElement
	Amount               *ValueDataElement
	Currency             *CurrencyDataElement
	TransmissionDate     *DateDataElement
	TransmissionTime     *TimeDataElement
}

func (b *BalanceDataElement) Elements() []DataElement {
	return []DataElement{
		b.DebitCreditIndicator,
		b.Amount,
		b.Currency,
		b.TransmissionDate,
		b.TransmissionTime,
	}
}

func (b *BalanceDataElement) Balance() domain.Balance {
	sign := b.DebitCreditIndicator.Val()
	val := b.Amount.Val()
	if sign == "D" {
		val = -val
	}
	currency := b.Currency.Val()
	amount := domain.Amount{
		Amount:   val,
		Currency: currency,
	}
	balance := domain.Balance{
		Amount:           amount,
		TransmissionDate: b.TransmissionDate.Val(),
	}
	if transmissionTime := b.TransmissionTime; transmissionTime != nil {
		val := transmissionTime.Val()
		balance.TransmissionTime = &val
	}
	return balance
}

func (b *BalanceDataElement) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) < 4 {
		return fmt.Errorf("%T: Malformed marshaled value", b)
	}
	b.DebitCreditIndicator = &AlphaNumericDataElement{}
	err = b.DebitCreditIndicator.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	b.Amount = &ValueDataElement{}
	err = b.Amount.UnmarshalHBCI(elements[1])
	if err != nil {
		return err
	}
	b.Currency = &CurrencyDataElement{}
	err = b.Currency.UnmarshalHBCI(elements[2])
	if err != nil {
		return err
	}
	b.TransmissionDate = &DateDataElement{}
	err = b.TransmissionDate.UnmarshalHBCI(elements[3])
	if err != nil {
		return err
	}
	if len(elements) == 5 {
		b.TransmissionTime = &TimeDataElement{}
		err = b.TransmissionTime.UnmarshalHBCI(elements[4])
		if err != nil {
			return err
		}
	}
	b.DataElement = NewGroupDataElementGroup(BalanceGDEG, 5, b)
	return nil
}

func (b *BalanceDataElement) Date() time.Time {
	return b.TransmissionDate.Val()
}

func NewAddress(address domain.Address) *AddressDataElement {
	a := &AddressDataElement{
		Name1:       NewAlphaNumeric(address.Name1, 35),
		Name2:       NewAlphaNumeric(address.Name2, 35),
		Street:      NewAlphaNumeric(address.Street, 35),
		PLZ:         NewAlphaNumeric(address.PLZ, 10),
		City:        NewAlphaNumeric(address.City, 35),
		CountryCode: NewCountryCode(address.CountryCode),
		Phone:       NewAlphaNumeric(address.Phone, 35),
		Fax:         NewAlphaNumeric(address.Fax, 35),
		Email:       NewAlphaNumeric(address.Email, 35),
	}
	a.DataElement = NewGroupDataElementGroup(AddressGDEG, 9, a)
	return a
}

type AddressDataElement struct {
	DataElement
	Name1       *AlphaNumericDataElement
	Name2       *AlphaNumericDataElement
	Street      *AlphaNumericDataElement
	PLZ         *AlphaNumericDataElement
	City        *AlphaNumericDataElement
	CountryCode *CountryCodeDataElement
	Phone       *AlphaNumericDataElement
	Fax         *AlphaNumericDataElement
	Email       *AlphaNumericDataElement
}

func (a *AddressDataElement) Elements() []DataElement {
	return []DataElement{
		a.Name1,
		a.Name2,
		a.Street,
		a.PLZ,
		a.City,
		a.CountryCode,
		a.Phone,
		a.Fax,
		a.Email,
	}
}

func (a *AddressDataElement) Address() domain.Address {
	return domain.Address{
		Name1:       a.Name1.Val(),
		Name2:       a.Name2.Val(),
		Street:      a.Street.Val(),
		PLZ:         a.PLZ.Val(),
		City:        a.City.Val(),
		CountryCode: a.CountryCode.Val(),
		Phone:       a.Phone.Val(),
		Fax:         a.Fax.Val(),
		Email:       a.Email.Val(),
	}
}