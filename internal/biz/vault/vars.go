package vault

import (
	"time"

	v1 "cobo-ucw-backend/api/ucw/v1"

	CoboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

const (
	DefaultQueryPortalDuration = 20 * time.Second
)

var (
	TssStatusMap = map[CoboWaas2.TSSRequestStatus]v1.TssRequest_Status{
		CoboWaas2.TSSREQUESTSTATUS_PENDING_KEY_HOLDER_CONFIRMATION: v1.TssRequest_STATUS_PENDING_KEYHOLDER_CONFIRMATION,
		CoboWaas2.TSSREQUESTSTATUS_KEY_HOLDER_CONFIRMATION_FAILED:  v1.TssRequest_STATUS_KEYHOLDER_CONFIRMATION_FAILED,
		CoboWaas2.TSSREQUESTSTATUS_KEY_GENERATING:                  v1.TssRequest_STATUS_KEY_GENERATING,
		CoboWaas2.TSSREQUESTSTATUS_MPC_PROCESSING:                  v1.TssRequest_STATUS_MPC_PROCESSING,
		CoboWaas2.TSSREQUESTSTATUS_KEY_GENERATING_FAILED:           v1.TssRequest_STATUS_KEY_GENERATING_FAILED,
		CoboWaas2.TSSREQUESTSTATUS_SUCCESS:                         v1.TssRequest_STATUS_SUCCESS,
	}
)
