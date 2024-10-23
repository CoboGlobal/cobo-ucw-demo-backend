package transaction

import (
	"time"

	v1 "cobo-ucw-backend/api/ucw/v1"

	CoboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

const (
	DefaultQueryPortalDuration = 20 * time.Second
)

var (
	transactionTypeSet = map[CoboWaas2.TransactionType]v1.Transaction_Type{
		CoboWaas2.TRANSACTIONTYPE_WITHDRAWAL: v1.Transaction_WITHDRAW,
		CoboWaas2.TRANSACTIONTYPE_DEPOSIT:    v1.Transaction_DEPOSIT,
	}

	statusSet = map[CoboWaas2.TransactionStatus]v1.Transaction_Status{
		CoboWaas2.TRANSACTIONSTATUS_SUBMITTED:             v1.Transaction_STATUS_SUBMITTED,
		CoboWaas2.TRANSACTIONSTATUS_PENDING_SCREENING:     v1.Transaction_STATUS_PENDING_SCREENING,
		CoboWaas2.TRANSACTIONSTATUS_PENDING_AUTHORIZATION: v1.Transaction_STATUS_PENDING_AUTHORIZATION,
		CoboWaas2.TRANSACTIONSTATUS_PENDING_SIGNATURE:     v1.Transaction_STATUS_PENDING_SIGNATURE,
		CoboWaas2.TRANSACTIONSTATUS_BROADCASTING:          v1.Transaction_STATUS_BROADCASTING,
		CoboWaas2.TRANSACTIONSTATUS_CONFIRMING:            v1.Transaction_STATUS_CONFIRMING,
		CoboWaas2.TRANSACTIONSTATUS_COMPLETED:             v1.Transaction_STATUS_SUCCESS,
		CoboWaas2.TRANSACTIONSTATUS_REJECTED:              v1.Transaction_STATUS_REJECTED,
		CoboWaas2.TRANSACTIONSTATUS_FAILED:                v1.Transaction_STATUS_FAILED,
	}
)
