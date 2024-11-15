package service

import (
	"context"
	"time"

	pb "cobo-ucw-backend/api/ucw/v1"
	"cobo-ucw-backend/internal/biz"
	"cobo-ucw-backend/internal/biz/transaction"
	biz_vault "cobo-ucw-backend/internal/biz/vault"
	"cobo-ucw-backend/internal/biz/wallet"
	"cobo-ucw-backend/internal/conf"
	"cobo-ucw-backend/internal/data/model"
	"cobo-ucw-backend/internal/middleware/auth"

	CoboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UserControlWalletService struct {
	pb.UnimplementedUserControlWalletServer
	UserUsecase        biz.UserUsecase
	TransactionUsecase biz.TransactionUsecase
	WalletUsecase      biz.WalletUsecase
	VaultUsecase       biz.VaultUsecase
	ucw                *conf.UCW
	logger             *log.Helper
	am                 *auth.JwtMiddleware
}

func NewUserControlWalletService(
	ucw *conf.UCW,
	logger log.Logger,
	usecase biz.VaultUsecase,
	transactionUsecase biz.TransactionUsecase,
	walletUsecase biz.WalletUsecase,
	userUsecase biz.UserUsecase, am *auth.JwtMiddleware) *UserControlWalletService {
	return &UserControlWalletService{
		UnimplementedUserControlWalletServer: pb.UnimplementedUserControlWalletServer{},
		UserUsecase:                          userUsecase,
		TransactionUsecase:                   transactionUsecase,
		WalletUsecase:                        walletUsecase,
		VaultUsecase:                         usecase,
		ucw:                                  ucw,
		logger:                               log.NewHelper(logger),
		am:                                   am,
	}
}

func (s *UserControlWalletService) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	return &pb.PingReply{
		Timestamp: time.Now().String(),
	}, nil
}

func (s *UserControlWalletService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	user, err := s.UserUsecase.Login(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	token, err := s.am.GenRegisteredClaims(user.UserID)
	if err != nil {
		return nil, err
	}
	return &pb.LoginReply{
		Token: token,
	}, nil
}

func (s *UserControlWalletService) BindNode(ctx context.Context, req *pb.BindNodeRequest) (*pb.BindNodeReply, error) {
	userInfo, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("BindNode")
	}

	userNode, err := s.UserUsecase.BindUserNode(ctx, userInfo.GetUserId(), req.NodeId)
	return &pb.BindNodeReply{
		UserNode: userNode.ToProto(),
	}, err
}

func (s *UserControlWalletService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoReply, error) {
	userInfo, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("GetUserInfo")
	}

	user, err := s.UserUsecase.GetUserInfo(ctx, userInfo.GetUserId())
	if err != nil {
		return nil, err
	}
	var vault *biz_vault.Vault
	var wallet *wallet.Wallet
	var groups biz_vault.Groups
	if user.GetVaultID() != "" {
		vault, err = s.VaultUsecase.GetVault(ctx, user.GetVaultID())
		if err != nil {
			return nil, err
		}

		wallet, err = s.WalletUsecase.GetWalletByVaultID(ctx, vault.VaultID)
		if err != nil {
			return nil, err
		}

		groups, err = s.VaultUsecase.ListGroups(ctx, vault.VaultID, "", 0)
		if err != nil {
			return nil, err
		}
	}
	userNode, err := s.UserUsecase.GetUserNodes(ctx, userInfo.GetUserId())
	for _, each := range userNode {
		each.Role = groups.UserNodeRole(each.NodeID)
	}
	return &pb.GetUserInfoReply{
		User:      user.ToProto(),
		Vault:     vault.ToProto(),
		Wallet:    wallet.ToProto(),
		UserNodes: userNode.ToProto(),
	}, err
}

func (s *UserControlWalletService) InitVault(ctx context.Context, req *pb.InitVaultRequest) (*pb.InitVaultReply, error) {
	userInfo, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("InitVault")
	}
	user, err := s.UserUsecase.GetUserInfo(ctx, userInfo.GetUserId())
	if err != nil {
		return nil, err
	}

	var vault *biz_vault.Vault
	if user.GetVaultID() == "" {
		vault, err = s.VaultUsecase.CreateVault(ctx, s.ucw.ProjectId)
		if err != nil {
			return nil, err
		}
		if err := s.UserUsecase.BindUserVault(ctx, userInfo.GetUserId(), vault.VaultID); err != nil {
			return nil, err
		}
	} else {
		vault, err = s.VaultUsecase.GetVault(ctx, user.UserVault.VaultID)
		if err != nil {
			return nil, err
		}
	}
	return &pb.InitVaultReply{
		Vault: vault.ToProto(),
	}, nil
}

func (s *UserControlWalletService) ListGroups(ctx context.Context, req *pb.ListGroupsRequest) (*pb.ListGroupsReply, error) {
	groups, err := s.VaultUsecase.ListGroups(ctx, req.VaultId, "", req.GetGroupType())
	if err != nil {
		return nil, err
	}
	var res []*pb.Group
	for _, each := range groups {
		res = append(res, each.ToProto())
	}
	return &pb.ListGroupsReply{
		Groups: res,
	}, nil
}

func (s *UserControlWalletService) GenerateMainGroup(ctx context.Context, req *pb.GenerateMainGroupRequest) (*pb.GenerateMainGroupReply, error) {
	userInfo, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("GenerateMainGroup")
	}
	user, err := s.UserUsecase.GetUserInfo(ctx, userInfo.ID)
	if err != nil {
		return nil, err
	}
	if req.GetVaultId() != user.GetVaultID() {
		return nil, errors.Forbidden("", "GenerateMainGroup")
	}
	vault, err := s.VaultUsecase.GetVault(ctx, req.GetVaultId())
	if err != nil {
		return nil, err
	}
	if vault.VaultID == "" {
		return nil, pb.ErrorVaultNotFound("not found %v", req.GetVaultId())
	}
	if vault.MainGroupID == "" {
		groupID, err := s.VaultUsecase.CreateKeyGroup(ctx, vault.VaultID, []*biz_vault.GroupNode{{
			GroupNode: &model.GroupNode{
				NodeID:     req.GetNodeId(),
				GroupID:    "",
				HolderName: user.Email,
				UserID:     userInfo.GetUserId(),
			},
		}}, pb.Group_MAIN_GROUP)
		if err != nil {
			return nil, err
		}
		vault.MainGroupID = groupID
	}

	tssRequestID, err := s.VaultUsecase.KeyGen(ctx, userInfo.GetUserId(), req.GetVaultId(), "", vault.MainGroupID, pb.Group_MAIN_GROUP)
	if err != nil {
		return nil, err
	}
	return &pb.GenerateMainGroupReply{
		TssRequestId: tssRequestID,
	}, nil
}

func (s *UserControlWalletService) GenerateRecoveryGroup(ctx context.Context, req *pb.GenerateRecoveryGroupRequest) (*pb.GenerateRecoveryGroupReply, error) {
	userInfo, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("GenerateRecoveryGroup")
	}
	user, err := s.UserUsecase.GetUserInfo(ctx, userInfo.ID)
	if err != nil {
		return nil, err
	}
	if req.GetVaultId() != user.GetVaultID() {
		return nil, errors.Forbidden("", "GenerateRecoveryGroup")
	}
	var groupNodes []*biz_vault.GroupNode
	vault, err := s.VaultUsecase.GetVault(ctx, req.GetVaultId())
	if err != nil {
		return nil, errors.NotFound("", "GenerateRecoveryGroup")
	}

	userNodes, err := s.UserUsecase.GetUserNodeByNodeIDs(ctx, req.NodeIds)
	if err != nil {
		return nil, err
	}
	for _, each := range userNodes {
		groupNodes = append(groupNodes, &biz_vault.GroupNode{
			GroupNode: &model.GroupNode{
				NodeID:     each.NodeID,
				GroupID:    "",
				HolderName: user.Email,
				UserID:     user.UserID,
			},
		})
	}
	groupID, err := s.VaultUsecase.CreateKeyGroup(ctx, req.GetVaultId(), groupNodes, pb.Group_RECOVERY_GROUP)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	tssRequestID, err := s.VaultUsecase.KeyGen(ctx, userInfo.GetUserId(), req.GetVaultId(), vault.MainGroupID, groupID, pb.Group_RECOVERY_GROUP)
	if err != nil {
		return nil, err
	}
	return &pb.GenerateRecoveryGroupReply{
		TssRequestId: tssRequestID,
	}, nil
}

func (s *UserControlWalletService) RecoverMainGroup(ctx context.Context, req *pb.RecoverMainGroupRequest) (*pb.RecoverMainGroupReply, error) {
	userInfo, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("RecoverMainGroup")
	}

	user, err := s.UserUsecase.GetUserInfo(ctx, userInfo.GetUserId())
	if err != nil {
		return nil, err
	}
	if req.GetVaultId() != user.GetVaultID() {
		return nil, errors.Forbidden("", "RecoverMainGroup")
	}
	groupID, err := s.VaultUsecase.CreateKeyGroup(ctx, req.GetVaultId(), []*biz_vault.GroupNode{
		{
			GroupNode: &model.GroupNode{
				NodeID:     req.GetNodeId(),
				GroupID:    "",
				HolderName: user.Email,
				UserID:     user.UserID,
			},
		},
	}, pb.Group_MAIN_GROUP)
	if err != nil {
		return nil, err
	}

	tssRequestID, err := s.VaultUsecase.KeyRecover(ctx, userInfo.GetUserId(), req.GetVaultId(), req.GetSourceGroupId(), groupID)
	if err != nil {
		return nil, err
	}
	return &pb.RecoverMainGroupReply{
		TssRequestId: tssRequestID,
	}, nil
}

func (s *UserControlWalletService) ListTssRequest(ctx context.Context, req *pb.ListTssRequestRequest) (*pb.ListTssRequestReply, error) {
	userInfo, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("ListTssRequest")
	}

	nodes, err := s.UserUsecase.GetUserNodeByNodeIDs(ctx, []string{req.GetNodeId()})
	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return nil, pb.ErrorInvalidRequestParams("ListTssRequest node need to bind")
	}

	if nodes[0].UserID != userInfo.GetUserId() {
		return nil, errors.Forbidden("", "ListTssRequest")
	}

	res, err := s.VaultUsecase.ListTssRequests(ctx, userInfo.GetUserId(), req.GetNodeId(), req.GetStatus())
	if err != nil {
		return nil, err
	}
	var list = make([]*pb.TssRequest, 0, len(res))
	for _, each := range res {
		list = append(list, each.ToProto())
	}
	return &pb.ListTssRequestReply{
		TssRequests: list,
	}, nil
}

func (s *UserControlWalletService) GetTssRequest(ctx context.Context, req *pb.GetTssRequestRequest) (*pb.GetTssRequestReply, error) {
	_, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("GetTssRequest")
	}

	res, err := s.VaultUsecase.GetTssRequest(ctx, req.GetTssRequestId())
	if err != nil {
		return nil, err
	}

	return &pb.GetTssRequestReply{
		TssRequest: res.ToProto(),
	}, nil
}

func (s *UserControlWalletService) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.GetGroupReply, error) {
	_, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("GetGroup")
	}

	res, err := s.VaultUsecase.GetGroup(ctx, req.GetVaultId(), req.GetGroupId())
	if err != nil {
		return nil, err
	}

	return &pb.GetGroupReply{
		Group: res.ToProtoGroupInfo(),
	}, nil
}

func (s *UserControlWalletService) DisasterRecovery(ctx context.Context, req *pb.DisasterRecoveryRequest) (*pb.DisasterRecoveryReply, error) {
	userInfo, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("DisasterRecovery")
	}

	user, err := s.UserUsecase.GetUserInfo(ctx, userInfo.GetUserId())
	if err != nil {
		return nil, err
	}
	if req.GetVaultId() != user.GetVaultID() {
		return nil, errors.Forbidden("", "DisasterRecovery")
	}
	vault, err := s.VaultUsecase.GetVault(ctx, req.GetVaultId())
	if err != nil {
		return nil, err
	}

	wallet, err := s.WalletUsecase.GetWalletByVaultID(ctx, req.GetVaultId())
	if err != nil {
		return nil, err
	}

	list, err := s.WalletUsecase.GetWalletAddress(ctx, wallet.WalletID, "")
	if err != nil {
		return nil, err
	}
	var addresses []*pb.Address
	for _, each := range list {
		addresses = append(addresses, each.ToProto())
	}
	return &pb.DisasterRecoveryReply{
		Vault:     vault.ToProto(),
		Wallet:    wallet.ToProto(),
		Addresses: addresses,
	}, nil
}

func (s *UserControlWalletService) CreateWallet(ctx context.Context, req *pb.CreateWalletRequest) (*pb.CreateWalletReply, error) {
	userInfo, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("CreateWallet")
	}
	user, err := s.UserUsecase.GetUserInfo(ctx, userInfo.GetUserId())
	if err != nil {
		return nil, err
	}
	if req.GetVaultId() != user.GetVaultID() {
		return nil, errors.Forbidden("", "CreateWallet")
	}

	userWallet, err := s.WalletUsecase.GetWalletByVaultID(ctx, user.GetVaultID())
	if err != nil {
		return nil, err
	}

	if userWallet.GetWalletID() != "" {
		return &pb.CreateWalletReply{
			WalletId: userWallet.GetWalletID(),
		}, nil
	}

	walletID, err := s.WalletUsecase.CreateWallet(ctx, &wallet.Wallet{
		Wallet: &model.Wallet{
			VaultID:  req.GetVaultId(),
			WalletID: uuid.New().String(),
			Name:     req.GetName(),
			UserID:   userInfo.GetUserId(),
		},
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateWalletReply{
		WalletId: walletID,
	}, nil
}

func (s *UserControlWalletService) GetWalletInfo(ctx context.Context, req *pb.GetWalletInfoRequest) (*pb.GetWalletInfoReply, error) {
	userInfo, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("GetWalletInfo")
	}
	walletInfo, err := s.WalletUsecase.GetWalletInfo(ctx, req.GetWalletId())
	if err != nil {
		return nil, err
	}
	if walletInfo.UserID != userInfo.GetUserId() {
		return nil, errors.Forbidden("", "InvalidWalletID")
	}
	return &pb.GetWalletInfoReply{
		WalletInfo: walletInfo.ToProtoWalletInfo(),
	}, nil
}

func (s *UserControlWalletService) AddWalletAddress(ctx context.Context, req *pb.AddWalletAddressRequest) (*pb.AddWalletAddressReply, error) {
	if err := s.CheckWalletResource(ctx, req.GetWalletId()); err != nil {
		return nil, err
	}
	address, err := s.WalletUsecase.AddWalletTokenAddress(ctx, req.GetWalletId(), req.GetChainId())
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return &pb.AddWalletAddressReply{
		Address: address.ToProto(),
	}, nil
}

func (s *UserControlWalletService) ListWalletToken(ctx context.Context, req *pb.ListWalletTokenRequest) (*pb.ListWalletTokenReply, error) {
	if err := s.CheckWalletResource(ctx, req.GetWalletId()); err != nil {
		return nil, err
	}
	walletTokens, err := s.WalletUsecase.ListWalletTokens(ctx, req.GetWalletId())
	return &pb.ListWalletTokenReply{
		List: walletTokens,
	}, err
}
func (s *UserControlWalletService) GetWalletToken(ctx context.Context, req *pb.GetWalletTokenRequest) (*pb.GetWalletTokenReply, error) {
	if err := s.CheckWalletResource(ctx, req.GetWalletId()); err != nil {
		return nil, err
	}
	walletToken, err := s.WalletUsecase.GetWalletToken(ctx, req.GetWalletId(), req.GetTokenId())
	if err != nil {
		return nil, err
	}

	return &pb.GetWalletTokenReply{
		Wallet:         walletToken.Wallet,
		TokenAddresses: walletToken.Token,
	}, nil
}
func (s *UserControlWalletService) GetTokenBalance(ctx context.Context, req *pb.GetTokenBalanceRequest) (*pb.GetTokenBalanceReply, error) {
	if err := s.CheckWalletResource(ctx, req.GetWalletId()); err != nil {
		return nil, err
	}
	balance, err := s.WalletUsecase.GetBalance(ctx, req.GetWalletId(), req.GetTokenId(), req.GetAddress())
	if err != nil {
		return nil, err
	}
	return &pb.GetTokenBalanceReply{
		TokenBalance: balance.ToProto(),
	}, err
}

func (s *UserControlWalletService) CheckWalletResource(ctx context.Context, walletID string) error {
	userInfo, ok := auth.FromContext(ctx)
	if !ok {
		return pb.ErrorUnauthorized("CheckWalletResource")
	}
	wallet, err := s.WalletUsecase.GetWalletInfo(ctx, walletID)
	if err != nil {
		return err
	}
	if wallet.UserID != userInfo.GetUserId() {
		return errors.Forbidden("", "InvalidWalletID")
	}
	return nil
}

func (s *UserControlWalletService) EstimateTransactionFee(ctx context.Context, req *pb.EstimateTransactionFeeRequest) (*pb.EstimateTransactionFeeReply, error) {
	amount, err := decimal.NewFromString(req.GetAmount())
	if err != nil {
		return nil, pb.ErrorInvalidRequestParams("invalid amount")
	}

	if err := s.CheckWalletResource(ctx, req.GetWalletId()); err != nil {
		return nil, err
	}
	fee, err := s.TransactionUsecase.PrepareTransaction(ctx, &transaction.Transaction{
		Transaction: &model.Transaction{
			WalletID: req.GetWalletId(),
			Type:     int64(req.GetType()),
			Amount:   amount,
			From:     req.GetFrom(),
			To:       req.GetTo(),
			TokenID:  req.TokenId,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.EstimateTransactionFeeReply{
		Slow:      fee.Slow.ToProto(),
		Recommend: fee.Recommend.ToProto(),
		Fast:      fee.Fast.ToProto(),
	}, nil
}
func (s *UserControlWalletService) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionReply, error) {
	userInfo, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("GetTransaction")
	}
	amount, err := decimal.NewFromString(req.GetAmount())
	if err != nil {
		return nil, pb.ErrorInvalidRequestParams("invalid amount")
	}
	fee, err := transaction.BuildFeeFromProto(req.GetFee())
	if err != nil {
		return nil, pb.ErrorInvalidRequestParams("invalid fee %v", err)
	}
	if err := s.CheckWalletResource(ctx, req.GetWalletId()); err != nil {
		return nil, err
	}
	id, err := s.TransactionUsecase.CreateTransaction(ctx, &transaction.Transaction{
		Transaction: &model.Transaction{
			WalletID: req.GetWalletId(),
			Type:     int64(req.GetType()),
			Chain:    req.GetChain(),
			Amount:   amount,
			From:     req.From,
			To:       req.To,
			Fee:      model.Fee(*fee),
			TokenID:  req.GetTokenId(),
			UserID:   userInfo.GetUserId(),
		},
	})
	return &pb.CreateTransactionReply{
		TransactionId: id,
	}, err
}
func (s *UserControlWalletService) ListTransaction(ctx context.Context, req *pb.ListTransactionRequest) (*pb.ListTransactionReply, error) {
	if err := s.CheckWalletResource(ctx, req.GetWalletId()); err != nil {
		return nil, err
	}

	res, err := s.TransactionUsecase.ListTransaction(ctx, transaction.ListTransactionParams{
		WalletID:        req.GetWalletId(),
		TokenID:         req.GetTokenId(),
		TransactionType: req.GetTransactionType(),
	})
	if err != nil {
		return nil, err
	}
	var list = make([]*pb.Transaction, 0, len(res))
	for _, each := range res {
		list = append(list, each.ToProto())
	}
	return &pb.ListTransactionReply{
		List: list,
	}, nil
}
func (s *UserControlWalletService) GetTransaction(ctx context.Context, req *pb.GetTransactionRequest) (*pb.GetTransactionReply, error) {
	userInfo, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("GetTransaction")
	}

	tx, err := s.TransactionUsecase.GetTransaction(ctx, req.GetTransactionId())

	if err != nil {
		return nil, err
	}

	if tx.UserID != userInfo.GetUserId() {
		return nil, errors.Forbidden("", "GetTransaction")
	}

	return &pb.GetTransactionReply{
		Transaction: tx.ToProto(),
	}, nil
}
func (s *UserControlWalletService) TransactionWebhook(ctx context.Context, req *pb.TransactionWebhookRequest) (*pb.TransactionWebhookReply, error) {

	switch CoboWaas2.TransactionType(req.Data.Type) {
	case CoboWaas2.TRANSACTIONTYPE_DEPOSIT:
		tx, err := transaction.BuildTransactionFromWebhook(req.GetData(), s.logger)
		if err != nil {
			return nil, err
		}
		_, err = s.TransactionUsecase.SyncDepositTransaction(ctx, tx)
		if err != nil {
			return nil, err
		}
	case CoboWaas2.TRANSACTIONTYPE_WITHDRAWAL, CoboWaas2.TRANSACTIONTYPE_CONTRACT_CALL:
		res, err := s.TransactionUsecase.ListTransaction(ctx, transaction.ListTransactionParams{
			ExternalID: req.Data.GetTransactionId(),
		})
		if err != nil {
			return nil, err
		}
		return &pb.TransactionWebhookReply{}, s.TransactionUsecase.SyncTransactions(ctx, res)
	default:
		log.Warnf("not handle tx type %s", req.Type)
	}

	return &pb.TransactionWebhookReply{}, nil
}

func (s *UserControlWalletService) Callback(ctx context.Context, req *pb.CoboCallbackRequest) (*pb.CoboCallbackReply, error) {
	return &pb.CoboCallbackReply{}, nil
}

func (s *UserControlWalletService) TssRequestWebhook(ctx context.Context, req *pb.TssRequestWebhookRequest) (*pb.TssRequestWebhookReply, error) {
	tssRequest, err := s.VaultUsecase.GetTssRequest(ctx, req.GetData().GetTssRequestId())
	if err != nil {
		return nil, err
	}
	return &pb.TssRequestWebhookReply{}, s.VaultUsecase.SyncTssRequests(ctx, []*biz_vault.TssRequest{tssRequest})
}

func (s *UserControlWalletService) TransactionReport(ctx context.Context, req *pb.TransactionReportRequest) (*pb.TransactionReportReply, error) {
	_, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("TransactionReport")
	}
	if req.Action == pb.TransactionReportRequest_ACTION_APPROVED {
		return &pb.TransactionReportReply{}, s.TransactionUsecase.ApproveTransaction(ctx, req.GetTransactionId())
	} else if req.Action == pb.TransactionReportRequest_ACTION_REJECTED {
		return &pb.TransactionReportReply{}, s.TransactionUsecase.RejectTransaction(ctx, req.GetTransactionId())
	}
	s.logger.Infof("TransactionReport %+v", req)
	return &pb.TransactionReportReply{}, nil
}

func (s *UserControlWalletService) TssRequestReport(ctx context.Context, req *pb.TssRequestReportRequest) (*pb.TssRequestReportReply, error) {
	_, ok := auth.FromContext(ctx)
	if !ok {
		return nil, pb.ErrorUnauthorized("TssRequestReport")
	}
	if req.Action == pb.TssRequestReportRequest_ACTION_APPROVED {
		return &pb.TssRequestReportReply{}, s.VaultUsecase.ApproveTssRequest(ctx, req.GetTssRequestId())
	} else if req.Action == pb.TssRequestReportRequest_ACTION_REJECTED {
		return &pb.TssRequestReportReply{}, s.VaultUsecase.RejectTssRequest(ctx, req.GetTssRequestId())
	}
	s.logger.Infof("TssRequestRepost %+v", req)
	return &pb.TssRequestReportReply{}, nil
}
