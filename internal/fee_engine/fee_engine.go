package fee_engine

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/storage"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func CalculateFee(details TransactionDetails) (float64, FeeBreakdown, error) {
	breakdown := FeeBreakdown{}

	userTier, err := getUserTier(details.UserID)
	if err != nil {
		return 0, breakdown, fmt.Errorf("could not determine user tier: %w", err)
	}

	config, err := fetchFeeConfig(details.Type, userTier)
	if err != nil {
		return 0, breakdown, fmt.Errorf("could not retrieve fee configuration: %w", err)
	}

	breakdown.AppliedRuleInfo = fmt.Sprintf("Rule for Type: %s, Tier: %s (Base: %.2f%%, Min: %.2f, Max: %.2f)",
		config.TransactionType, config.Tier, config.BasePercentage*100, config.MinFee, config.MaxFee)

	baseFee := details.Amount * config.BasePercentage
	breakdown.BaseFee = baseFee
	calculatedFee := baseFee

	
	isPeak, surchargeAmount := applyTimeSurcharge(details.Timestamp, details.Amount, config)
	if isPeak {
		calculatedFee += surchargeAmount
		breakdown.TimeSurcharge = surchargeAmount
		breakdown.AppliedRuleInfo += fmt.Sprintf(" + Peak Surcharge: %.2f (%.2f%%)", surchargeAmount, config.PeakSurcharge*100)
	}
	
	breakdown.CalculatedFee = calculatedFee
	
	finalFee := calculatedFee
	appliedMinMax := false
	if config.MaxFee > 0 && finalFee > config.MaxFee {
		finalFee = config.MaxFee
		appliedMinMax = true
		breakdown.AppliedRuleInfo += fmt.Sprintf(" | Applied Max Cap: %.2f", config.MaxFee)
	}
	if finalFee < config.MinFee {
		finalFee = config.MinFee
		appliedMinMax = true
		breakdown.AppliedRuleInfo += fmt.Sprintf(" | Applied Min Floor: %.2f", config.MinFee)

	}
	
	breakdown.AppliedMinMax = appliedMinMax
	breakdown.FinalFee = math.Round(finalFee*100) / 100 // Round to 2 decimal places
	
	log.Printf("Fee calculated for TxType %s, User %s, Tier %s, Amount %.2f -> Fee: %.2f",
	details.Type, details.UserID, userTier, details.Amount, breakdown.FinalFee)
	
	breakdown.TierAdjustment = breakdown.FinalFee // Tier adjustment is the final fee because this is a simple fee engine
	return breakdown.FinalFee, breakdown, nil 
}

func getUserTier(userID string) (UserTier, error) {
	var tier string
	// check inside the cache first
	// if its not inside the cache, query the db
	err := storage.DB.QueryRow("SELECT tier FROM users WHERE user_id = $1", userID).Scan(&tier)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found: %s", userID)
		}
		return "", fmt.Errorf("database error fetching user tier: %w", err)
	}
	
	if tier == string(Basic) || tier == string(Premium) || tier == string(Enterprise) {
		return UserTier(tier), nil
	}

	return "", fmt.Errorf("unrecognized user tier: %s", tier)
}

// fetchFeeConfig retrieves the most specific fee configuration for the transaction.
// Tier-specific rules override general transaction type rules.
func fetchFeeConfig(txType TransactionType, tier UserTier) (*FeeConfig, error) {
	var config FeeConfig
	// Prioritize tier-specific rule
	query := `
        SELECT transaction_type, tier, base_percentage, min_fee, max_fee, peak_start_time, peak_end_time, peak_surcharge
        FROM fee_configs
        WHERE transaction_type = $1 AND tier = $2
        ORDER BY tier DESC NULLS LAST -- Prioritize specific tier rule
        LIMIT 1`

	row := storage.DB.QueryRow(query, txType, tier) // using placeholders
	err := row.Scan(
		&config.TransactionType, &config.Tier, &config.BasePercentage,
		&config.MinFee, &config.MaxFee, &config.PeakStartTime, &config.PeakEndTime, &config.PeakSurcharge,
	)

	// no tier specific rule
	if err == sql.ErrNoRows {
		queryBase := `
            SELECT transaction_type, tier, base_percentage, min_fee, max_fee, peak_start_time, peak_end_time, peak_surcharge
            FROM fee_configs
            WHERE transaction_type = $1 AND tier IS NULL
            LIMIT 1`
		rowBase := storage.DB.QueryRow(queryBase, txType)
		errBase := rowBase.Scan(
			&config.TransactionType, &config.Tier, &config.BasePercentage,
			&config.MinFee, &config.MaxFee, &config.PeakStartTime, &config.PeakEndTime, &config.PeakSurcharge,
		)
		if errBase != nil {
			if errBase == sql.ErrNoRows {
				return nil, fmt.Errorf("no fee configuration found for transaction type %s and tier %s", txType, tier)
			}
			return nil, fmt.Errorf("database error fetching base fee config: %w", errBase)
		}
		log.Printf("Using base fee config for type %s", txType)
		return &config, nil
	} else if err != nil {
		return nil, fmt.Errorf("database error fetching tier-specific fee config: %w", err)
	}

	log.Printf("Using tier-specific fee config for type %s, tier %s", txType, tier)
	return &config, nil
}

// to check if the transaction occurs during peak hours and adds surcharge.
func applyTimeSurcharge(txTime time.Time, amount float64, config *FeeConfig) (bool, float64) {
	if config.PeakSurcharge <= 0 || config.PeakStartTime == "" || config.PeakEndTime == "" {
		return false, 0 
	}

	format := "15:04"
	peakStart, _ := time.Parse(format, config.PeakStartTime)
	peakEnd, _ := time.Parse(format, config.PeakEndTime)

	// Get current time in HH:MM format for comparison
	currentTime := txTime.Format(format)
	currentParsed, _ := time.Parse(format, currentTime)

	// PeakStartTime < PeakEndTime for simple logoc
	isPeak := false
	if !currentParsed.Before(peakStart) && currentParsed.Before(peakEnd) {
		isPeak = true
	}

	if isPeak {
		surcharge := amount * config.PeakSurcharge
		return true, surcharge
	}

	return false, 0
}