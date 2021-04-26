package sharesies

import "time"

const (
	OrderTypeDollarMarket = "dollar_market"

	PaymentCurrency = "nzd"
	PaymentType     = "direct"
)

type ProfileResponse struct {
	Authenticated           bool            `json:"authenticated"`
	AutoinvestOrder         AutoinvestOrder `json:"autoinvest_order"`
	CanEnterAddressToken    bool            `json:"can_enter_address_token"`
	CanWriteUntil           CanWriteUntil   `json:"can_write_until"`
	DistillToken            string          `json:"distill_token"`
	Flags                   []string        `json:"flags"`
	GaID                    string          `json:"ga_id"`
	LiveData                LiveData        `json:"live_data"`
	NzxIsOpen               bool            `json:"nzx_is_open"`
	NzxNextOpen             NzxNextOpen     `json:"nzx_next_open"`
	Orders                  []interface{}   `json:"orders"`
	OutstandingSubscription interface{}     `json:"outstanding_subscription"`
	Participants            []string        `json:"participants"`
	Portfolio               []Portfolio     `json:"portfolio"`
	RakaiaToken             interface{}     `json:"rakaia_token"`
	RakaiaTokenExpiry       interface{}     `json:"rakaia_token_expiry"`
	ReferralCode            string          `json:"referral_code"`
	Type                    string          `json:"type"`
	UpcomingDividends       []interface{}   `json:"upcoming_dividends"`
	User                    User            `json:"user"`
	UserList                []UserList      `json:"user_list"`
}
type Allocations struct {
	Allocation string `json:"allocation"`
	FundID     string `json:"fund_id"`
}
type AutoinvestOrder struct {
	Allocations    []Allocations `json:"allocations"`
	Amount         string        `json:"amount"`
	Interval       string        `json:"interval"`
	LastFailedDate interface{}   `json:"last_failed_date"`
	NextDate       string        `json:"next_date"`
	PremadeOrderID interface{}   `json:"premade_order_id"`
	State          string        `json:"state"`
}
type CanWriteUntil struct {
	Quantum int64 `json:"$quantum"`
}
type LiveData struct {
	EligibleForFreeMonth bool `json:"eligible_for_free_month"`
	IsActive             bool `json:"is_active"`
}
type NzxNextOpen struct {
	Quantum int64 `json:"$quantum"`
}
type Stats struct {
	CapitalReturn        string `json:"capital_return"`
	SharesBought         string `json:"shares_bought"`
	SharesSold           string `json:"shares_sold"`
	SharesTransferredIn  string `json:"shares_transferred_in"`
	SharesTransferredOut string `json:"shares_transferred_out"`
	ValueBought          string `json:"value_bought"`
	ValueSold            string `json:"value_sold"`
	ValueTransferredIn   string `json:"value_transferred_in"`
	ValueTransferredOut  string `json:"value_transferred_out"`
}
type Portfolio struct {
	Contribution        string `json:"contribution"`
	Currency            string `json:"currency"`
	CurrentTaxLiability string `json:"current_tax_liability"`
	Dividends           string `json:"dividends"`
	FundID              string `json:"fund_id"`
	GrossValue          string `json:"gross_value"`
	HoldingType         string `json:"holding_type"`
	ReturnDollars       string `json:"return_dollars"`
	ReturnPercent       string `json:"return_percent"`
	RiskRating          int    `json:"risk_rating"`
	Shares              string `json:"shares"`
	Stats               Stats  `json:"stats"`
	Value               string `json:"value"`
}
type Components struct {
	Locality     string `json:"locality"`
	PostalCode   string `json:"postal_code"`
	Route        string `json:"route"`
	StreetNumber string `json:"street_number"`
	Sublocality  string `json:"sublocality"`
}
type Address struct {
	Components Components  `json:"components"`
	Formatted  string      `json:"formatted"`
	Lat        interface{} `json:"lat"`
	Lng        interface{} `json:"lng"`
}
type Checks struct {
	AddressEntered       bool `json:"address_entered"`
	AddressVerified      bool `json:"address_verified"`
	DependentDeclaration bool `json:"dependent_declaration"`
	IDVerified           bool `json:"id_verified"`
	MadeDeposit          bool `json:"made_deposit"`
	PrescribedAnswered   bool `json:"prescribed_answered"`
	TaxQuestions         bool `json:"tax_questions"`
	TcAccepted           bool `json:"tc_accepted"`
}
type HasSeen struct {
	AuSharesIntro        bool `json:"au_shares_intro"`
	Autoinvest           bool `json:"autoinvest"`
	Companies            bool `json:"companies"`
	ExchangeInvestor     bool `json:"exchange_investor"`
	Funds                bool `json:"funds"`
	Investor             bool `json:"investor"`
	LimitOrders          bool `json:"limit_orders"`
	ManagedFundsInvestor bool `json:"managed_funds_investor"`
	ShowAuCurrency       bool `json:"show_au_currency"`
}
type TaxResidencies struct {
	Country     string `json:"country"`
	CountryName string `json:"country_name"`
	Tin         string `json:"tin"`
}
type WalletBalances struct {
	Aud string `json:"aud"`
	Nzd string `json:"nzd"`
	Usd string `json:"usd"`
}
type User struct {
	AccountFrozen           bool             `json:"account_frozen"`
	AccountReference        string           `json:"account_reference"`
	AccountRestricted       bool             `json:"account_restricted"`
	AccountRestrictedDate   interface{}      `json:"account_restricted_date"`
	Address                 Address          `json:"address"`
	AddressRejectReason     interface{}      `json:"address_reject_reason"`
	AddressState            string           `json:"address_state"`
	Checks                  Checks           `json:"checks"`
	Email                   string           `json:"email"`
	FirstTaxYear            int              `json:"first_tax_year"`
	HasSeen                 HasSeen          `json:"has_seen"`
	HoldingBalance          string           `json:"holding_balance"`
	HomeCurrency            string           `json:"home_currency"`
	ID                      string           `json:"id"`
	Intercom                string           `json:"intercom"`
	IrdNumber               string           `json:"ird_number"`
	IsDependent             bool             `json:"is_dependent"`
	IsOwnerPrescribed       bool             `json:"is_owner_prescribed"`
	Jurisdiction            string           `json:"jurisdiction"`
	MaximumWithdrawalAmount string           `json:"maximum_withdrawal_amount"`
	MinimumWalletBalance    string           `json:"minimum_wallet_balance"`
	ParticipantEmails       []interface{}    `json:"participant_emails"`
	Phone                   string           `json:"phone"`
	Pir                     string           `json:"pir"`
	PortfolioID             string           `json:"portfolio_id"`
	PreferredName           string           `json:"preferred_name"`
	PrescribedApproved      bool             `json:"prescribed_approved"`
	PrescribedParticipant   interface{}      `json:"prescribed_participant"`
	RecentSearches          []string         `json:"recent_searches"`
	SeenFirstTimeAutoinvest bool             `json:"seen_first_time_autoinvest"`
	SeenFirstTimeInvestor   bool             `json:"seen_first_time_investor"`
	State                   string           `json:"state"`
	TaxResidencies          []TaxResidencies `json:"tax_residencies"`
	TaxYear                 int              `json:"tax_year"`
	TfnNumber               interface{}      `json:"tfn_number"`
	TransferAge             interface{}      `json:"transfer_age"`
	TransferAgePassed       bool             `json:"transfer_age_passed"`
	UsEquitiesEnabled       bool             `json:"us_equities_enabled"`
	UsTaxTreatyStatus       string           `json:"us_tax_treaty_status"`
	WalletBalances          WalletBalances   `json:"wallet_balances"`
}
type UserList struct {
	ID            string `json:"id"`
	PreferredName string `json:"preferred_name"`
	Primary       bool   `json:"primary"`
	State         string `json:"state"`
}

type InstrumentsRequest struct {
	Page            int      `json:"page"`
	Perpage         int      `json:"perPage"`
	Sort            string   `json:"sort"`
	Pricechangetime string   `json:"priceChangeTime"`
	Query           string   `json:"query"`
	Instruments     []string `json:"instruments"`
}

type InstrumentResponse struct {
	Total          int       `json:"total"`
	Currentpage    int       `json:"currentPage"`
	Resultsperpage int       `json:"resultsPerPage"`
	Numberofpages  int       `json:"numberOfPages"`
	Instruments    []Company `json:"instruments"`
}
type Logos struct {
	Wide  string `json:"wide"`
	Thumb string `json:"thumb"`
	Micro string `json:"micro"`
}

type Company struct {
	ID                        string      `json:"id"`
	Urlslug                   string      `json:"urlSlug"`
	Instrumenttype            string      `json:"instrumentType"`
	Symbol                    string      `json:"symbol"`
	Kidsrecommended           bool        `json:"kidsRecommended"`
	Isvolatile                bool        `json:"isVolatile"`
	Name                      string      `json:"name"`
	Description               string      `json:"description"`
	Categories                []string    `json:"categories"`
	Logoidentifier            string      `json:"logoIdentifier"`
	Logos                     Logos       `json:"logos"`
	Riskrating                int         `json:"riskRating"`
	Marketprice               string      `json:"marketPrice"`
	Marketlastcheck           time.Time   `json:"marketLastCheck"`
	Tradingstatus             string      `json:"tradingStatus"`
	Exchangecountry           string      `json:"exchangeCountry"`
	Peratio                   string      `json:"peRatio"`
	Marketcap                 int64       `json:"marketCap"`
	Websiteurl                string      `json:"websiteUrl"`
	Exchange                  string      `json:"exchange"`
	Legacyimageurl            interface{} `json:"legacyImageUrl"`
	Dominantcolour            string      `json:"dominantColour"`
	Pdsdriveid                interface{} `json:"pdsDriveId"`
	Assetmanager              interface{} `json:"assetManager"`
	Fixedfeespread            interface{} `json:"fixedFeeSpread"`
	Managementfeepercent      interface{} `json:"managementFeePercent"`
	Grossdividendyieldpercent string      `json:"grossDividendYieldPercent"`
	Annualisedreturnpercent   string      `json:"annualisedReturnPercent"`
	Ceo                       string      `json:"ceo"`
	Employees                 int         `json:"employees"`
}

type CostBuyRequest struct {
	FundID     string `json:"fund_id"`
	ActingAsID string `json:"acting_as_id"`
	Order      *Order `json:"order"`
}

type Order struct {
	Type           string `json:"type"`
	CurrencyAmount string `json:"currency_amount"`
}

type CostBuyResponse struct {
	ExpectedFee      string             `json:"expected_fee"`
	FundID           string             `json:"fund_id"`
	PaymentBreakdown []PaymentBreakdown `json:"payment_breakdown"`
	Request          *Order             `json:"request"`
	TotalCost        string             `json:"total_cost"`
	Type             string             `json:"type"`
}

type PaymentBreakdown struct {
	Currency     string `json:"currency"`
	TargetAmount string `json:"target_amount"`
	Type         string `json:"type"`
}

type CreateBuyRequest struct {
	FundID           string              `json:"fund_id"`
	ActingAsID       string              `json:"acting_as_id"`
	Order            *Order              `json:"order"`
	IdempotencyKey   string              `json:"idempotency_key"`
	PaymentBreakdown *[]PaymentBreakdown `json:"payment_breakdown"`
	ExpectedFee      string              `json:"expected_fee"`
}
