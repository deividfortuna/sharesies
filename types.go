package sharesies

import "time"

const (
	OrderTypeDollarMarket = "dollar_market"
	OrderTypeShareMarket = "share_market"

	PaymentCurrency = "nzd"
	PaymentType     = "direct"
)

type ProfileResponse struct {
	Authenticated           bool             `json:"authenticated" validate:"required"`
	AutoinvestOrder         *AutoinvestOrder `json:"autoinvest_order" validate:"required"`
	CanEnterAddressToken    bool             `json:"can_enter_address_token" validate:"required"`
	CanWriteUntil           *CanWriteUntil   `json:"can_write_until" validate:"required"`
	DistillToken            string           `json:"distill_token" validate:"required"`
	Flags                   []string         `json:"flags" validate:"required"`
	GaID                    string           `json:"ga_id" validate:"required"`
	LiveData                *LiveData        `json:"live_data" validate:"required"`
	NzxIsOpen               bool             `json:"nzx_is_open" validate:"required"`
	NzxNextOpen             *NzxNextOpen     `json:"nzx_next_open" validate:"required"`
	Orders                  []interface{}    `json:"orders"`
	OutstandingSubscription interface{}      `json:"outstanding_subscription"`
	Participants            []string         `json:"participants" validate:"required"`
	Portfolio               []*Portfolio     `json:"portfolio" validate:"required"`
	RakaiaToken             interface{}      `json:"rakaia_token"`
	RakaiaTokenExpiry       interface{}      `json:"rakaia_token_expiry"`
	ReferralCode            string           `json:"referral_code" validate:"required"`
	Type                    string           `json:"type" validate:"required"`
	UpcomingDividends       []interface{}    `json:"upcoming_dividends"`
	User                    *User            `json:"user" validate:"required"`
	UserList                []*UserSimple    `json:"user_list" validate:"required"`
}
type Allocations struct {
	Allocation string `json:"allocation" validate:"required"`
	FundID     string `json:"fund_id" validate:"required"`
}
type AutoinvestOrder struct {
	Allocations    []Allocations `json:"allocations" validate:"required"`
	Amount         string        `json:"amount" validate:"required"`
	Interval       string        `json:"interval" validate:"required"`
	LastFailedDate interface{}   `json:"last_failed_date"`
	NextDate       string        `json:"next_date" validate:"required"`
	PremadeOrderID interface{}   `json:"premade_order_id"`
	State          string        `json:"state" validate:"required"`
}
type CanWriteUntil struct {
	Quantum int64 `json:"$quantum" validate:"required"`
}
type LiveData struct {
	EligibleForFreeMonth bool `json:"eligible_for_free_month" validate:"required"`
	IsActive             bool `json:"is_active" validate:"required"`
}
type NzxNextOpen struct {
	Quantum int64 `json:"$quantum" validate:"required"`
}
type Stats struct {
	CapitalReturn        string `json:"capital_return" validate:"required"`
	SharesBought         string `json:"shares_bought" validate:"required"`
	SharesSold           string `json:"shares_sold" validate:"required"`
	SharesTransferredIn  string `json:"shares_transferred_in" validate:"required"`
	SharesTransferredOut string `json:"shares_transferred_out" validate:"required"`
	ValueBought          string `json:"value_bought" validate:"required"`
	ValueSold            string `json:"value_sold" validate:"required"`
	ValueTransferredIn   string `json:"value_transferred_in" validate:"required"`
	ValueTransferredOut  string `json:"value_transferred_out" validate:"required"`
}
type Portfolio struct {
	Contribution        string `json:"contribution" validate:"required"`
	Currency            string `json:"currency" validate:"required"`
	CurrentTaxLiability string `json:"current_tax_liability" validate:"required"`
	Dividends           string `json:"dividends" validate:"required"`
	FundID              string `json:"fund_id" validate:"required"`
	GrossValue          string `json:"gross_value" validate:"required"`
	HoldingType         string `json:"holding_type" validate:"required"`
	ReturnDollars       string `json:"return_dollars" validate:"required"`
	ReturnPercent       string `json:"return_percent" validate:"required"`
	RiskRating          int    `json:"risk_rating" validate:"required"`
	Shares              string `json:"shares" validate:"required"`
	Stats               *Stats `json:"stats" validate:"required"`
	Value               string `json:"value" validate:"required"`
}
type Components struct {
	Locality     string `json:"locality" validate:"required"`
	PostalCode   string `json:"postal_code" validate:"required"`
	Route        string `json:"route" validate:"required"`
	StreetNumber string `json:"street_number" validate:"required"`
	Sublocality  string `json:"sublocality" validate:"required"`
}
type Address struct {
	Components Components  `json:"components" validate:"required"`
	Formatted  string      `json:"formatted" validate:"required"`
	Lat        interface{} `json:"lat"`
	Lng        interface{} `json:"lng"`
}
type Checks struct {
	AddressEntered       bool `json:"address_entered" validate:"required"`
	AddressVerified      bool `json:"address_verified" validate:"required"`
	DependentDeclaration bool `json:"dependent_declaration" validate:"required"`
	IDVerified           bool `json:"id_verified" validate:"required"`
	MadeDeposit          bool `json:"made_deposit" validate:"required"`
	PrescribedAnswered   bool `json:"prescribed_answered" validate:"required"`
	TaxQuestions         bool `json:"tax_questions" validate:"required"`
	TcAccepted           bool `json:"tc_accepted" validate:"required"`
}
type HasSeen struct {
	AuSharesIntro        bool `json:"au_shares_intro" validate:"required"`
	Autoinvest           bool `json:"autoinvest" validate:"required"`
	Companies            bool `json:"companies" validate:"required"`
	ExchangeInvestor     bool `json:"exchange_investor" validate:"required"`
	Funds                bool `json:"funds" validate:"required"`
	Investor             bool `json:"investor" validate:"required"`
	LimitOrders          bool `json:"limit_orders" validate:"required"`
	ManagedFundsInvestor bool `json:"managed_funds_investor" validate:"required"`
	ShowAuCurrency       bool `json:"show_au_currency" validate:"required"`
}
type TaxResidencies struct {
	Country     string `json:"country" validate:"required"`
	CountryName string `json:"country_name" validate:"required"`
	Tin         string `json:"tin" validate:"required"`
}
type WalletBalances struct {
	Aud string `json:"aud" validate:"required"`
	Nzd string `json:"nzd" validate:"required"`
	Usd string `json:"usd" validate:"required"`
}
type User struct {
	AccountFrozen           bool              `json:"account_frozen" validate:"required"`
	AccountReference        string            `json:"account_reference" validate:"required"`
	AccountRestricted       bool              `json:"account_restricted" validate:"required"`
	AccountRestrictedDate   interface{}       `json:"account_restricted_date"`
	Address                 *Address          `json:"address" validate:"required"`
	AddressRejectReason     interface{}       `json:"address_reject_reason"`
	AddressState            string            `json:"address_state" validate:"required"`
	Checks                  *Checks           `json:"checks" validate:"required"`
	Email                   string            `json:"email" validate:"required"`
	FirstTaxYear            int               `json:"first_tax_year" validate:"required"`
	HasSeen                 *HasSeen          `json:"has_seen" validate:"required"`
	HoldingBalance          string            `json:"holding_balance" validate:"required"`
	HomeCurrency            string            `json:"home_currency" validate:"required"`
	ID                      string            `json:"id" validate:"required"`
	Intercom                string            `json:"intercom" validate:"required"`
	IrdNumber               string            `json:"ird_number" validate:"required"`
	IsDependent             bool              `json:"is_dependent" validate:"required"`
	IsOwnerPrescribed       bool              `json:"is_owner_prescribed" validate:"required"`
	Jurisdiction            string            `json:"jurisdiction" validate:"required"`
	MaximumWithdrawalAmount string            `json:"maximum_withdrawal_amount" validate:"required"`
	MinimumWalletBalance    string            `json:"minimum_wallet_balance" validate:"required"`
	ParticipantEmails       []interface{}     `json:"participant_emails"`
	Phone                   string            `json:"phone" validate:"required"`
	Pir                     string            `json:"pir" validate:"required"`
	PortfolioID             string            `json:"portfolio_id" validate:"required"`
	PreferredName           string            `json:"preferred_name" validate:"required"`
	PrescribedApproved      bool              `json:"prescribed_approved" validate:"required"`
	PrescribedParticipant   interface{}       `json:"prescribed_participant"`
	RecentSearches          []string          `json:"recent_searches" validate:"required"`
	SeenFirstTimeAutoinvest bool              `json:"seen_first_time_autoinvest" validate:"required"`
	SeenFirstTimeInvestor   bool              `json:"seen_first_time_investor" validate:"required"`
	State                   string            `json:"state" validate:"required"`
	TaxResidencies          []*TaxResidencies `json:"tax_residencies" validate:"required"`
	TaxYear                 int               `json:"tax_year" validate:"required"`
	TfnNumber               interface{}       `json:"tfn_number"`
	TransferAge             interface{}       `json:"transfer_age"`
	TransferAgePassed       bool              `json:"transfer_age_passed" validate:"required"`
	UsEquitiesEnabled       bool              `json:"us_equities_enabled" validate:"required"`
	UsTaxTreatyStatus       string            `json:"us_tax_treaty_status" validate:"required"`
	WalletBalances          *WalletBalances   `json:"wallet_balances" validate:"required"`
}
type UserSimple struct {
	ID            string `json:"id" validate:"required"`
	PreferredName string `json:"preferred_name" validate:"required"`
	Primary       bool   `json:"primary" validate:"required"`
	State         string `json:"state" validate:"required"`
}

type InstrumentsRequest struct {
	Page            int      `json:"page" validate:"required"`
	Perpage         int      `json:"perPage" validate:"required"`
	Sort            string   `json:"sort" validate:"required"`
	Pricechangetime string   `json:"priceChangeTime" validate:"required"`
	Query           string   `json:"query"`
	Instruments     []string `json:"instruments"`
}

type InstrumentResponse struct {
	Total          int        `json:"total" validate:"required"`
	Currentpage    int        `json:"currentPage" validate:"required"`
	Resultsperpage int        `json:"resultsPerPage" validate:"required"`
	Numberofpages  int        `json:"numberOfPages" validate:"required"`
	Instruments    []*Company `json:"instruments" validate:"required"`
}
type Logos struct {
	Wide  string `json:"wide" validate:"required"`
	Thumb string `json:"thumb" validate:"required"`
	Micro string `json:"micro" validate:"required"`
}

type Company struct {
	ID                        string      `json:"id" validate:"required"`
	Urlslug                   string      `json:"urlSlug" validate:"required"`
	Instrumenttype            string      `json:"instrumentType" validate:"required"`
	Symbol                    string      `json:"symbol" validate:"required"`
	Kidsrecommended           bool        `json:"kidsRecommended" validate:"required"`
	Isvolatile                bool        `json:"isVolatile" validate:"required"`
	Name                      string      `json:"name" validate:"required"`
	Description               string      `json:"description" validate:"required"`
	Categories                []string    `json:"categories" validate:"required"`
	Logoidentifier            string      `json:"logoIdentifier" validate:"required"`
	Logos                     *Logos      `json:"logos" validate:"required"`
	Riskrating                int         `json:"riskRating" validate:"required"`
	Marketprice               string      `json:"marketPrice" validate:"required"`
	Marketlastcheck           time.Time   `json:"marketLastCheck" validate:"required"`
	Tradingstatus             string      `json:"tradingStatus" validate:"required"`
	Exchangecountry           string      `json:"exchangeCountry" validate:"required"`
	Peratio                   string      `json:"peRatio" validate:"required"`
	Marketcap                 int64       `json:"marketCap" validate:"required"`
	Websiteurl                string      `json:"websiteUrl" validate:"required"`
	Exchange                  string      `json:"exchange" validate:"required"`
	Legacyimageurl            interface{} `json:"legacyImageUrl"`
	Dominantcolour            string      `json:"dominantColour" validate:"required"`
	Pdsdriveid                interface{} `json:"pdsDriveId"`
	Assetmanager              interface{} `json:"assetManager"`
	Fixedfeespread            interface{} `json:"fixedFeeSpread"`
	Managementfeepercent      interface{} `json:"managementFeePercent"`
	Grossdividendyieldpercent string      `json:"grossDividendYieldPercent" validate:"required"`
	Annualisedreturnpercent   string      `json:"annualisedReturnPercent" validate:"required"`
	Ceo                       string      `json:"ceo" validate:"required"`
	Employees                 int         `json:"employees" validate:"required"`
}

// Buy Transactions Types

type CostBuyRequest struct {
	FundID     string    `json:"fund_id" validate:"required"`
	ActingAsID string    `json:"acting_as_id" validate:"required"`
	Order      *OrderBuy `json:"order" validate:"required"`
}

type OrderBuy struct {
	Type           string `json:"type" validate:"required"`
	CurrencyAmount string `json:"currency_amount,omitempty"`
	ShareAmount    string `json:"share_amount,omitempty"`
}

type PaymentBreakdown struct {
	Currency     string `json:"currency" validate:"required"`
	TargetAmount string `json:"target_amount" validate:"required"`
	Type         string `json:"type" validate:"required"`
}

type CostBuyResponse struct {
	ExpectedFee      string              `json:"expected_fee" validate:"required"`
	FundID           string              `json:"fund_id" validate:"required"`
	PaymentBreakdown []*PaymentBreakdown `json:"payment_breakdown" validate:"required"`
	Request          *OrderBuy           `json:"request" validate:"required"`
	TotalCost        string              `json:"total_cost" validate:"required"`
	Type             string              `json:"type" validate:"required"`
}

type CreateBuyRequest struct {
	FundID           string              `json:"fund_id" validate:"required"`
	ActingAsID       string              `json:"acting_as_id" validate:"required"`
	Order            *OrderBuy           `json:"order" validate:"required"`
	IdempotencyKey   string              `json:"idempotency_key" validate:"required"`
	PaymentBreakdown []*PaymentBreakdown `json:"payment_breakdown" validate:"required"`
	ExpectedFee      string              `json:"expected_fee" validate:"required"`
}

// Sell Transactions Types

type CostSellRequest struct {
	FundID     string `json:"fund_id" validate:"required"`
	ActingAsID string `json:"acting_as_id" validate:"required"`
	Order      *OrderSell `json:"order" validate:"required"`
}

type OrderSell struct {
	Type           string `json:"type" validate:"required"`
	ShareAmount    string `json:"share_amount" validate:"required"`
}

type CostSellResponse struct {
	FundID  string `json:"fund_id" validate:"required"`
	Request *OrderSell `json:"request" validate:"required"`
	Type    string `json:"type" validate:"required"`
}

type CreateSellRequest struct {
	FundID           string              `json:"fund_id" validate:"required"`
	ActingAsID       string              `json:"acting_as_id" validate:"required"`
	Order            *OrderSell           `json:"order" validate:"required"`
	IdempotencyKey   string              `json:"idempotency_key" validate:"required"`
}
