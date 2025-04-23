package fee_engine

import "time"

type TransactionType string

const (
	BillPayment TransactionType = "BILL_PAYMENT"
	WalletTopUp TransactionType = "WALLET_TOPUP"
	Transfer    TransactionType = "TRANSFER"
)

type UserTier string

const (
	Basic      UserTier = "Basic"
	Premium    UserTier = "Premium"
	Enterprise UserTier = "Enterprise"
)

type TransactionDetails struct {
	Amount          float64
	Type            TransactionType
	UserID          string // To fetch user tier
	Timestamp       time.Time
	ServiceProvider string // Optional: For provider-specific fees
}

// FeeConfig holds the dynamic fee configuration, fetched from DB.
// This is simplified; a real system might has a more complex rule structures.
type FeeConfig struct {
	TransactionType TransactionType
	Tier            UserTier // Can be empty for base rules
	BasePercentage  float64
	MinFee          float64 // Floor limit
	MaxFee          float64 // Cap limit
	PeakStartTime   string  // e.g., "18:00"
	PeakEndTime     string  // e.g., "22:00"
	PeakSurcharge   float64 // Additional percentage during peak hours
}

// FeeBreakdown provides details about how the fee was calculated.
type FeeBreakdown struct {
	BaseFee         float64 `json:"base_fee"`
	TierAdjustment  float64 `json:"tier_adjustment"` // Amount added/subtracted due to tier
	TimeSurcharge   float64 `json:"time_surcharge"`
	CalculatedFee   float64 `json:"calculated_fee"` // Before caps/floors
	AppliedMinMax   bool    `json:"applied_min_max"`
	FinalFee        float64 `json:"final_fee"`
	AppliedRuleInfo string  `json:"applied_rule_info"` // Description of the rule used
}
