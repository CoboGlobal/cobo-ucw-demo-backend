package user

import (
	v1 "cobo-ucw-backend/api/ucw/v1"
	"cobo-ucw-backend/internal/data/model"

	CoboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

type User struct {
	*model.User

	UserVault *UserVault
}

func (u *User) GetVaultID() string {
	if u.UserVault != nil {
		return u.UserVault.VaultID
	}
	return ""
}

func (u *User) ToProto() *v1.User {
	if u == nil || u.User == nil {
		return nil
	}
	return &v1.User{
		UserId: u.UserID,
		Email:  u.Email,
	}
}

type UserNode struct {
	*model.UserNode

	Role v1.Role
	User User
}

type UserVault struct {
	*model.UserVault
}

func (n *UserNode) ToProto() *v1.UserNode {
	if n == nil || n.UserNode == nil {
		return nil
	}
	return &v1.UserNode{
		UserId: n.UserID,
		NodeId: n.NodeID,
		Role:   n.Role,
	}
}

type Nodes []*UserNode

func BuildNodes(nodes []*model.UserNode) Nodes {
	var res = make(Nodes, 0)
	for _, each := range nodes {
		res = append(res, &UserNode{UserNode: each})
	}
	return res
}

func (n Nodes) UniqueUser() bool {
	if n.Len() == 0 {
		return true
	}
	userID := n[0].UserID

	for _, node := range n {
		if node.UserID != userID {
			return false
		}
	}
	return true
}

func (n Nodes) ToProto() []*v1.UserNode {
	var res = make([]*v1.UserNode, 0)
	for _, each := range n {
		res = append(res, each.ToProto())
	}
	return res
}

func (n Nodes) Len() int {
	return len([]*UserNode(n))
}

func (n *UserNode) ToCreateKeyGroupRequestKeyHoldersInner() CoboWaas2.CreateKeyShareHolder {
	return CoboWaas2.CreateKeyShareHolder{
		Name:      CoboWaas2.PtrString(n.UserID),
		Type:      CoboWaas2.KEYSHAREHOLDERTYPE_API.Ptr(),
		TssNodeId: CoboWaas2.PtrString(n.NodeID),
	}
}

func (n Nodes) ToCreateKeyGroupRequestKeyHoldersInner() []CoboWaas2.CreateKeyShareHolder {
	var res = make([]CoboWaas2.CreateKeyShareHolder, 0)

	for _, each := range n {
		res = append(res, each.ToCreateKeyGroupRequestKeyHoldersInner())
	}
	return res
}
