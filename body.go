package fatturapa

import (
	"errors"
	"fmt"

	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
)

const (
	scontoMaggiorazioneTypeDiscount = "SC" // sconto
	scontoMaggiorazioneTypeCharge   = "MG" // maggiorazione
)

const (
	condizioniPagamentoInstallments = "TP01" // pagamenti in rate
	condizioniPagamentoFull         = "TP02" // pagamento completo
	condizioniPagamentoAdvance      = "TP03" // anticipo
)

const stampDutyCode = "SI"

// Body contains all invoice data apart from the parties involved, which are
// contained in Header.
type Body struct {
	GeneralData   *GeneralData   `xml:"DatiGenerali,omitempty"`
	GoodsServices *GoodsServices `xml:"DatiBeniServizi,omitempty"`
	PaymentData   *PaymentData   `xml:"DatiPagamento,omitempty"`
}

// GeneralData contains general data about the invoice such as retained taxes,
// invoice number, invoice date, document type, etc.
type GeneralData struct {
	Document  *GeneralDocumentData `xml:"DatiGeneraliDocumento"`
	Purchases []*DocumentRef       `xml:"DatiOrdineAcquisto,omitempty"`
	Contracts []*DocumentRef       `xml:"DatiContratto,omitempty"`
	Tender    []*DocumentRef       `xml:"DatiConvenzione,omitempty"`
	Receiving []*DocumentRef       `xml:"DatiRicezione,omitempty"`
	Preceding []*DocumentRef       `xml:"DatiFattureCollegate,omitempty"`
}

// DocumentRef contains data about a previous document.
type DocumentRef struct {
	Lines     []int  `xml:"RiferimentoNumeroLinea"`              // detail row of the invoice referred to (if the reference is to the entire invoice, this is not filled in)
	Code      string `xml:"IdDocumento"`                         // document number
	IssueDate string `xml:"Data,omitempty"`                      // document date (expressed according to the ISO 8601:2004 format)
	LineCode  string `xml:"NumItem,omitempty"`                   // identification of the single item on the document (e.g. in the case of a purchase order, this is the number of the row of the purchase order, or, in the case of a contract, it is the number of the row of the contract, etc. )
	OrderCode string `xml:"CodiceCommessaConvenzione,omitempty"` // order or agreement code
	CUPCode   string `xml:"CodiceCUP,omitempty"`                 // code managed by the CIPE (Interministerial Committee for Economic Planning) which characterises every public investment project (Individual Project Code).
	CIGCode   string `xml:"CodiceCIG,omitempty"`                 // Tender procedure identification code
}

type GeneralDocumentData struct {
	DocumentType     string             `xml:"TipoDocumento"`
	Divisa           string             `xml:"Divisa"`
	IssueDate        string             `xml:"Data"`
	Number           string             `xml:"Numero"`
	RetainedTaxes    []*RetainedTax     `xml:"DatiRitenuta,omitempty"`
	StampDuty        *StampDuty         `xml:"DatiBollo,omitempty"`
	PriceAdjustments []*PriceAdjustment `xml:"ScontoMaggiorazione,omitempty"`
	TotalAmount      string             `xml:"ImportoTotaleDocumento"`
	Causale          []string           `xml:"Causale,omitempty"`
}

// StampDuty contains data about the stamp duty
type StampDuty struct {
	VirtualStamp string `xml:"BolloVirtuale"`
	Amount       string `xml:"ImportoBollo,omitempty"`
}

// PriceAdjustment contains data about price adjustments like discounts and
// charges.
type PriceAdjustment struct {
	Type    string `xml:"Tipo"`
	Percent string `xml:"Percentuale,omitempty"`
	Amount  string `xml:"Importo,omitempty"`
}

func newBody(inv *bill.Invoice) (*Body, error) {
	dbs := newGoodsServices(inv)

	dp, err := newPaymentData(inv)
	if err != nil {
		return nil, err
	}

	dg, err := newGeneralData(inv)
	if err != nil {
		return nil, err
	}

	return &Body{
		GeneralData:   dg,
		GoodsServices: dbs,
		PaymentData:   dp,
	}, nil
}

func newGeneralData(inv *bill.Invoice) (*GeneralData, error) {
	gd := new(GeneralData)
	var err error
	if gd.Document, err = newGeneralDocumentData(inv); err != nil {
		return nil, err
	}
	gd.Preceding = newDocumentRefs(inv.Preceding)
	if o := inv.Ordering; o != nil {
		gd.Purchases = newDocumentRefs(o.Purchases)
		gd.Contracts = newDocumentRefs(o.Contracts)
		gd.Tender = newDocumentRefs(o.Tender)
		gd.Receiving = newDocumentRefs(o.Receiving)
	}
	return gd, nil
}

func newDocumentRefs(refs []*org.DocumentRef) []*DocumentRef {
	out := make([]*DocumentRef, len(refs))
	for i, ref := range refs {
		out[i] = newDocumentRef(ref)
	}
	return out
}

func newDocumentRef(ref *org.DocumentRef) *DocumentRef {
	dr := &DocumentRef{
		Lines: ref.Lines,
		Code:  ref.Series.Join(ref.Code).String(),
	}
	if ref.IssueDate != nil {
		dr.IssueDate = ref.IssueDate.String()
	}
	for _, id := range ref.Identities {
		switch id.Key {
		case org.IdentityKeyOrder:
			dr.OrderCode = string(id.Code)
		case org.IdentityKeyItem:
			dr.LineCode = string(id.Code)
		}
		switch id.Type {
		case sdi.IdentityTypeCIG:
			dr.CIGCode = string(id.Code)
		case sdi.IdentityTypeCUP:
			dr.CUPCode = string(id.Code)
		}
	}

	return dr
}

func newGeneralDocumentData(inv *bill.Invoice) (*GeneralDocumentData, error) {
	dr, err := extractRetainedTaxes(inv)
	if err != nil {
		return nil, err
	}

	codeTipoDocumento, err := findCodeTipoDocumento(inv)
	if err != nil {
		return nil, err
	}

	switch codeTipoDocumento {
	case "TD07", "TD08", "TD09":
		return nil, errors.New("simplified invoices are not currently supported")
	}

	code := inv.Code
	if inv.Series != "" {
		code = cbc.Code(fmt.Sprintf("%s-%s", inv.Series, inv.Code))
	}

	doc := &GeneralDocumentData{
		DocumentType:     codeTipoDocumento,
		Divisa:           string(inv.Currency),
		IssueDate:        inv.IssueDate.String(),
		Number:           code.String(),
		RetainedTaxes:    dr,
		StampDuty:        newStampDuty(inv.Charges),
		TotalAmount:      formatAmount2(&inv.Totals.Payable),
		PriceAdjustments: extractPriceAdjustments(inv),
		Causale:          extractInvoiceReasons(inv),
	}

	return doc, nil
}

func findCodeTipoDocumento(inv *bill.Invoice) (string, error) {
	if inv.Tax == nil {
		return "", fmt.Errorf("missing tax")
	}

	val, ok := inv.Tax.Ext[sdi.ExtKeyDocumentType]
	if !ok || val == "" {
		return "", fmt.Errorf("missing %s", sdi.ExtKeyDocumentType)
	}

	return val.String(), nil
}

func newStampDuty(charges []*bill.Charge) *StampDuty {
	for _, charge := range charges {
		if charge.Key == bill.ChargeKeyStampDuty {
			return &StampDuty{
				VirtualStamp: stampDutyCode,
				Amount:       formatAmount2(&charge.Amount),
			}
		}
	}

	return nil
}

func extractPriceAdjustments(inv *bill.Invoice) []*PriceAdjustment {
	var priceAdjustments []*PriceAdjustment

	for _, discount := range inv.Discounts {
		priceAdjustments = append(priceAdjustments, &PriceAdjustment{
			Type:    scontoMaggiorazioneTypeDiscount,
			Percent: formatPercentage(discount.Percent),
			Amount:  formatAmount8(&discount.Amount),
		})
	}

	for _, charge := range inv.Charges {
		priceAdjustments = append(priceAdjustments, &PriceAdjustment{
			Type:    scontoMaggiorazioneTypeCharge,
			Percent: formatPercentage(charge.Percent),
			Amount:  formatAmount8(&charge.Amount),
		})
	}

	return priceAdjustments
}

func extractInvoiceReasons(inv *bill.Invoice) []string {
	// find inv.Notes with NoteKey as cbc.NoteKeyReason
	var reasons []string

	for _, note := range inv.Notes {
		if note.Key == org.NoteKeyReason {
			reasons = append(reasons, note.Text)
		}
	}

	return reasons
}
