package constants

const (
	// Error Messages
	ErrCustomerNotFound    = "customer not found"
	ErrTransactionNotFound = "transaction not found"
	ErrLimitNotFound       = "customer limit not found"
	ErrInvalidCredentials  = "invalid credentials"
	ErrInsufficientLimit   = "insufficient credit limit"
	ErrDuplicateNIK        = "customer with this NIK already exists"
	ErrInvalidTenor        = "invalid tenor, only 1, 2, 3, 4 months allowed"

	// Success Messages
	MsgCustomerCreated    = "customer created successfully"
	MsgTransactionCreated = "transaction created successfully"
	MsgLoginSuccessful    = "login successful"
	MsgDataRetrieved      = "data retrieved successfully"

	// Contract Number Prefix
	ContractPrefix = "XYZ"

	// JWT Defaults
	DefaultJWTSecret      = "xyz-secret-key"
	DefaultJWTExpiryHours = 24

	// Rate Limiting
	DefaultRateLimit       = 100
	DefaultRateLimitWindow = 60 // seconds

	// Pagination
	DefaultPageLimit = 10
	MaxPageLimit     = 100
)

// ValidTenors defines the allowed tenor months for PT XYZ Multifinance
var ValidTenors = []int{1, 2, 3, 4}

// IsValidTenor checks if the given tenor is valid (only 1, 2, 3, 4 allowed)
func IsValidTenor(tenor int) bool {
	for _, validTenor := range ValidTenors {
		if tenor == validTenor {
			return true
		}
	}
	return false
}

// GetValidTenors returns the list of valid tenor months
func GetValidTenors() []int {
	return ValidTenors
}

// GetValidTenorsString returns valid tenors as comma-separated string
func GetValidTenorsString() string {
	return "1, 2, 3, 4"
}
