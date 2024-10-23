package vault

import (
	"strconv"

	v1 "cobo-ucw-backend/api/ucw/v1"
	"cobo-ucw-backend/internal/data/model"

	CoboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

type Group struct {
	*model.Group
	GroupNodes
}

func (g *Group) ToProto() *v1.Group {
	return &v1.Group{
		GroupId:   g.GroupID,
		GroupType: v1.Group_GroupType(g.GroupType),
	}
}

func (g *Group) ToProtoGroupInfo() *v1.GroupInfo {
	return &v1.GroupInfo{
		Group:      g.ToProto(),
		GroupNodes: g.GroupNodes.ToProto(),
	}
}

func (g *GroupNode) ToProto() *v1.GroupNode {
	if g == nil {
		return &v1.GroupNode{}
	}
	return &v1.GroupNode{
		GroupId:    g.GroupID,
		UserId:     g.UserID,
		NodeId:     g.NodeID,
		HolderName: g.HolderName,
	}
}
func (g *GroupNodes) ToProto() []*v1.GroupNode {
	var res []*v1.GroupNode
	for _, each := range *g {
		res = append(res, each.ToProto())
	}
	return res
}

type Groups []*Group

func (g Groups) UserNodeRole(nodeID string) v1.Role {
	var groupTypes []v1.Group_GroupType
	for _, each := range g {
		for _, groupNode := range each.GroupNodes {
			if groupNode.NodeID == nodeID {
				groupTypes = append(groupTypes, v1.Group_GroupType(each.GroupType))
			}
		}
	}

	if len(groupTypes) == 0 {
		return v1.Role_ROLE_UNSPECIFIED
	}

	if len(groupTypes) > 1 {
		return v1.Role_ROLE_ADMIN
	}

	if groupTypes[0] == v1.Group_RECOVERY_GROUP {
		return v1.Role_ROLE_RECOVERY
	}

	return v1.Role_ROLE_MAIN
}

type Vault struct {
	*model.Vault

	coboNodeID string
}

type TssRequest struct {
	*model.TssRequest
}

func (t *TssRequest) SetStatus(status v1.TssRequest_Status) {
	t.Status = int64(status)
}

func (t *TssRequest) ToProto() *v1.TssRequest {
	if t == nil || t.TssRequest == nil {
		return nil
	}
	return &v1.TssRequest{
		RequestId:       t.RequestID,
		Type:            v1.TssRequest_Type(t.Type),
		Status:          v1.TssRequest_Status(t.Status),
		SourceGroupId:   t.SourceGroupID,
		TargetGroupId:   t.TargetGroupID,
		CreateTimestamp: strconv.FormatInt(t.CreatedAt.Unix(), 10),
	}
}

type GroupNode struct {
	*model.GroupNode
}

type GroupNodes []*GroupNode

func (g *GroupNode) SetGroupID(groupID string) {
	g.GroupID = groupID
}

func (g *GroupNodes) ToCreateKeyShareHolders() []CoboWaas2.CreateKeyShareHolder {
	var res = make([]CoboWaas2.CreateKeyShareHolder, 0, len(*g))
	for _, each := range *g {
		res = append(res, CoboWaas2.CreateKeyShareHolder{
			Name:      CoboWaas2.PtrString(each.HolderName),
			Type:      CoboWaas2.KEYSHAREHOLDERTYPE_API.Ptr(),
			TssNodeId: CoboWaas2.PtrString(each.NodeID),
		})
	}
	return res
}

func (v *Vault) VaultStatus() v1.Vault_Status {
	if v == nil {
		return 0
	}
	return v1.Vault_Status(v.Status)
}

func (v *Vault) ToProto() *v1.Vault {
	if v == nil || v.Vault == nil {
		return nil
	}
	return &v1.Vault{
		VaultId:     v.VaultID,
		Name:        v.Name,
		MainGroupId: v.MainGroupID,
		ProjectId:   v.ProjectID,
		CoboNodeId:  v.coboNodeID,
		Status:      v.VaultStatus(),
	}
}
