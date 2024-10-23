package vault

import (
	"context"
	"fmt"
	"time"

	v1 "cobo-ucw-backend/api/ucw/v1"
	"cobo-ucw-backend/integration/portal"
	"cobo-ucw-backend/internal/conf"
	"cobo-ucw-backend/internal/data/database"
	"cobo-ucw-backend/internal/data/model"

	CoboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	errors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type Usecase struct {
	client         *portal.Client
	vaultRepo      Repo
	groupRepo      GroupRepo
	groupNodeRepo  GroupNodeRepo
	tssRequestRepo TssRequestRepo
	ucw            *conf.UCW
	data           *database.Data
	logger         *log.Helper
}

func NewUsecase(
	ucw *conf.UCW,
	client *portal.Client,
	vaultRepo Repo,
	groupNodeRepo GroupNodeRepo,
	groupRepo GroupRepo,
	tssRequestRepo TssRequestRepo,
	data *database.Data,
	logger log.Logger,
) *Usecase {
	return &Usecase{
		client:         client,
		vaultRepo:      vaultRepo,
		groupRepo:      groupRepo,
		groupNodeRepo:  groupNodeRepo,
		tssRequestRepo: tssRequestRepo,
		ucw:            ucw,
		data:           data,
		logger:         log.NewHelper(logger),
	}
}

func (u Usecase) CreateVault(ctx context.Context, projectID string) (*Vault, error) {
	vault := &Vault{
		Vault: &model.Vault{
			VaultID:     "",
			Name:        fmt.Sprintf("vault%d", time.Now().UnixMilli()),
			MainGroupID: "",
			ProjectID:   projectID,
			Status:      0,
		},
	}
	mpcVault, _, err := u.client.WalletsMPCWalletsAPI.CreateMpcVault(u.client.WithContext(ctx)).CreateMpcVaultRequest(CoboWaas2.CreateMpcVaultRequest{
		ProjectId: CoboWaas2.PtrString(projectID),
		Name:      vault.Name,
		VaultType: CoboWaas2.MPCVAULTTYPE_USER_CONTROLLED,
	}).Execute()
	if err != nil {
		u.logger.Errorf("CreateVault create vault for project id %s, err %v", projectID, err)
		return nil, err
	}

	vault.VaultID = mpcVault.GetVaultId()
	vault.Status = int64(v1.Vault_CREATED)
	_, err = u.vaultRepo.Save(ctx, vault)
	if err != nil {
		u.logger.Errorf("CreateVault save project id %s, vault id %s, err %v", projectID, vault.VaultID, err)
		return nil, err
	}
	return vault, err
}

func (u Usecase) KeyGen(ctx context.Context, userID, vaultID, sourceGroupID, targetGroupID string, groupType v1.Group_GroupType) (string, error) {
	tssType := CoboWaas2.TSSREQUESTTYPE_KEY_GEN
	key := int64(v1.TssRequest_GENERATE_MAIN_KEY)
	if groupType == v1.Group_RECOVERY_GROUP {
		tssType = CoboWaas2.TSSREQUESTTYPE_KEY_GEN_FROM_KEY_GROUP
		key = int64(v1.TssRequest_GENERATE_RECOVERY_KEY)
	}
	tssRequest, _, err := u.client.WalletsMPCWalletsAPI.CreateTssRequest(u.client.WithContext(ctx), vaultID).CreateTssRequestRequest(CoboWaas2.CreateTssRequestRequest{
		Type:                        tssType,
		TargetKeyShareHolderGroupId: targetGroupID,
		SourceKeyShareHolderGroup: &CoboWaas2.SourceGroup{
			KeyShareHolderGroupId: sourceGroupID,
		},
	}).Execute()
	if err != nil {
		u.logger.Errorf("KeyGen CreateTssRequest userID %s, vaultID %s, sourceGroupID %s, targetGroupID %s, groupType %s, err %v",
			userID, vaultID, sourceGroupID, targetGroupID, groupType, err)
		return "", err
	}
	_, err = u.tssRequestRepo.Save(ctx, &TssRequest{&model.TssRequest{
		RequestID:     tssRequest.GetTssRequestId(),
		Type:          key,
		Status:        int64(TssStatusMap[tssRequest.GetStatus()]),
		TargetGroupID: targetGroupID,
		SourceGroupID: sourceGroupID,
		UserID:        userID,
		VaultID:       vaultID,
	}})
	if err != nil {
		u.logger.Errorf("KeyGen Save userID %s, vaultID %s, sourceGroupID %s, targetGroupID %s, groupType %s, tss request id %s, err %v",
			userID, vaultID, sourceGroupID, targetGroupID, groupType, tssRequest.TssRequestId, err)
		return "", err
	}
	return tssRequest.GetTssRequestId(), nil
}

func (u Usecase) KeyRecover(ctx context.Context, userID, vaultID, sourceGroupID, targetGroupID string) (string, error) {
	tssRequest, _, err := u.client.WalletsMPCWalletsAPI.CreateTssRequest(u.client.WithContext(ctx), vaultID).CreateTssRequestRequest(CoboWaas2.CreateTssRequestRequest{
		Type:                        CoboWaas2.TSSREQUESTTYPE_RECOVERY,
		TargetKeyShareHolderGroupId: targetGroupID,
		SourceKeyShareHolderGroup: &CoboWaas2.SourceGroup{
			KeyShareHolderGroupId: sourceGroupID,
			TssNodeIds:            nil,
		},
	}).Execute()
	if err != nil {
		u.logger.Errorf("KeyRecover CreateTssRequest userID %s, vaultID %s, sourceGroupID %s, targetGroupID %s, err %v",
			userID, vaultID, sourceGroupID, targetGroupID, err)
		return "", err
	}
	_, err = u.tssRequestRepo.Save(ctx, &TssRequest{&model.TssRequest{
		RequestID:     tssRequest.GetTssRequestId(),
		Type:          int64(v1.TssRequest_RECOVERY_MAIN_KEY),
		Status:        int64(TssStatusMap[tssRequest.GetStatus()]),
		TargetGroupID: targetGroupID,
		SourceGroupID: sourceGroupID,
		UserID:        userID,
		VaultID:       vaultID,
	}})
	if err != nil {
		u.logger.Errorf("KeyRecover Save userID %s, vaultID %s, sourceGroupID %s, targetGroupID %s, tss request id %s, err %v",
			userID, vaultID, sourceGroupID, targetGroupID, tssRequest.TssRequestId, err)
		return "", err
	}
	return tssRequest.GetTssRequestId(), err
}

func (u Usecase) CreateKeyGroup(ctx context.Context, vaultID string, groupNodes GroupNodes, groupType v1.Group_GroupType) (string, error) {
	keyGroupType := CoboWaas2.KEYSHAREHOLDERGROUPTYPE_MAIN_GROUP
	if groupType == v1.Group_RECOVERY_GROUP {
		keyGroupType = CoboWaas2.KEYSHAREHOLDERGROUPTYPE_RECOVERY_GROUP
	}

	keyGroup, _, err := u.client.WalletsMPCWalletsAPI.CreateKeyShareHolderGroup(u.client.WithContext(ctx), vaultID).CreateKeyShareHolderGroupRequest(CoboWaas2.CreateKeyShareHolderGroupRequest{
		KeyShareHolderGroupType: keyGroupType,
		Participants:            u.ucw.Participant,
		Threshold:               u.ucw.Threshold,
		KeyShareHolders:         groupNodes.ToCreateKeyShareHolders(),
	}).Execute()
	if err != nil {
		u.logger.Errorf("CreateKeyGroup CreateKeyShareHolderGroup vaultID %s, group type %s, err %v", vaultID, groupType, err)
		return "", err
	}

	for _, each := range groupNodes {
		each.SetGroupID(keyGroup.GetKeyShareHolderGroupId())
	}
	groupID, err := keyGroup.GetKeyShareHolderGroupId(), u.data.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if groupType == v1.Group_MAIN_GROUP {
			if err := u.vaultRepo.Tx(tx).UpdateByVaultID(ctx, Vault{
				Vault: &model.Vault{
					MainGroupID: keyGroup.GetKeyShareHolderGroupId(),
					Status:      int64(v1.Vault_MAIN_GROUP_CREATED),
				},
			}, vaultID); err != nil {
				return err
			}
		}
		if _, err := u.groupRepo.Tx(tx).Save(ctx, &Group{
			Group: &model.Group{
				VaultID:   vaultID,
				GroupID:   keyGroup.GetKeyShareHolderGroupId(),
				GroupType: int64(groupType),
			},
		}); err != nil {
			return err
		}
		return u.groupNodeRepo.Tx(tx).BatchSave(ctx, groupNodes)
	})
	if err != nil {
		u.logger.Errorf("CreateKeyGroup Transaction vaultID %s, group id %s, group type %s, err %v",
			vaultID, groupID, groupType, err)
	}
	return groupID, err
}

func (u Usecase) GetVault(ctx context.Context, vaultID string) (*Vault, error) {
	vault, err := u.vaultRepo.GetByVaultID(ctx, vaultID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	vault.coboNodeID, err = u.client.GetCoboNodeID(ctx)
	if err != nil {
		u.logger.Errorf("GetVault GetCoboNodeID vaultID %s, err %v", vaultID, err)
		return nil, err
	}
	return vault, nil
}

func (u Usecase) ListTssRequests(ctx context.Context, userID, nodeID string, status v1.TssRequest_Status) ([]*TssRequest, error) {
	groupNodes, err := u.groupNodeRepo.ListNodeGroups(ctx, userID, nodeID)
	if err != nil {
		u.logger.Errorf("ListTssRequests ListNodeGroups userID %s, nodeID %s, status %s, err %v", userID, nodeID, status, err)
		return nil, err
	}
	groupIDs := make([]string, 0, len(groupNodes))
	for _, each := range groupNodes {
		groupIDs = append(groupIDs, each.GroupID)
	}
	return u.tssRequestRepo.ListGroupRelatedTssRequest(ctx, userID, groupIDs, int64(status))
}

func (u Usecase) ListGroups(ctx context.Context, vaultID string, nodeID string, groupType v1.Group_GroupType) (Groups, error) {
	groups, err := u.groupRepo.ListGroups(ctx, vaultID, groupType)
	if err != nil {
		u.logger.Errorf("ListGroups ListGroups vaultID %s, nodeID %s, groupType %s, err %v", vaultID, nodeID, groupType, err)
		return nil, err
	}
	if len(groups) == 0 {
		return []*Group{}, nil
	}
	groupIDs := make([]string, 0, len(groups))
	groupSet := make(map[string]*Group)
	for _, each := range groups {
		groupIDs = append(groupIDs, each.GroupID)
		groupSet[each.GroupID] = each
	}
	groupNodes, err := u.groupNodeRepo.ListByGroupIDs(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	for _, each := range groupNodes {
		if nodeID != "" && each.NodeID != nodeID {
			continue
		}
		groupSet[each.GroupID].GroupNodes = append(groupSet[each.GroupID].GroupNodes, each)
	}
	return groups, nil
}

func (u Usecase) GetGroup(ctx context.Context, vaultID, groupID string) (*Group, error) {
	group, err := u.groupRepo.GetByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	groupNodes, err := u.groupNodeRepo.ListByGroupIDs(ctx, []string{groupID})
	if err != nil {
		return nil, err
	}

	group.GroupNodes = groupNodes
	return group, nil
}

func (u Usecase) GetTssRequest(ctx context.Context, tssRequestID string) (*TssRequest, error) {
	return u.tssRequestRepo.GetByTssRequestID(ctx, tssRequestID)
}

func (u Usecase) SyncTssRequestTask() error {
	var lastID = int64(-1)
	ctx := context.Background()
	limit := 20
	for lastID != 0 {
		list, err := u.tssRequestRepo.ListTssRequest(ctx, ListTssRequestParams{
			VaultID:       "",
			TargetGroupID: "",
			SourceGroupID: "",
			Status: []v1.TssRequest_Status{
				v1.TssRequest_STATUS_PENDING_KEYHOLDER_CONFIRMATION,
				v1.TssRequest_STATUS_KEY_GENERATING,
				v1.TssRequest_STATUS_MPC_PROCESSING,
				v1.TssRequest_STATUS_UNSPECIFIED,
			},
			Limit:  limit,
			LastID: lastID,
		})
		if err != nil {
			return err
		}

		if len(list) == limit {
			lastID = int64(list[len(list)-1].ID)
		} else {
			lastID = 0
		}

		if err := u.SyncTssRequests(ctx, list); err != nil {
			u.logger.Errorf("SyncTssRequestTask err %v", err)
			return err
		}
	}
	return nil
}

func (u Usecase) SyncTssRequests(ctx context.Context, tssRequests []*TssRequest) error {
	for _, each := range tssRequests {
		if each.RequestID == "" {
			continue
		}
		u.logger.Infof("SyncTssRequests start sync tss request id %s, vaultID %s", each.RequestID, each.VaultID)
		response, _, err := u.client.WalletsMPCWalletsAPI.GetTssRequestById(u.client.WithContext(ctx), each.VaultID, each.RequestID).Execute()
		if err != nil {
			u.logger.Errorf("SyncTssRequests GetTssRequestById request id %s, vaultID %s, err %v", each.RequestID, each.VaultID, err)
			continue
		}
		each.SetStatus(TssStatusMap[response.GetStatus()])

		if each.Status == int64(v1.TssRequest_STATUS_SUCCESS) {
			if err := u.data.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				vault := Vault{
					Vault: &model.Vault{Status: int64(v1.Vault_MAIN_GENERATED)},
				}

				if each.Type == int64(v1.TssRequest_GENERATE_MAIN_KEY) || each.Type == int64(v1.TssRequest_RECOVERY_MAIN_KEY) {
					vault.MainGroupID = each.TargetGroupID
				}
				if err := u.vaultRepo.Tx(tx).UpdateByVaultID(ctx, vault, each.VaultID); err != nil {
					return err
				}
				return u.tssRequestRepo.Tx(tx).Update(ctx, *each)
			}); err != nil {
				u.logger.Errorf("SyncTssRequests Update request id %s, vaultID %s, err %v", each.RequestID, each.VaultID, err)
			}
		}
		if err := u.tssRequestRepo.Update(ctx, *each); err != nil {
			u.logger.Errorf("SyncTssRequests Update request id %s, vaultID %s, err %v", each.RequestID, each.VaultID, err)
		}
	}
	return nil
}
