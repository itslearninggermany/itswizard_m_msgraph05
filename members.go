package itswizard_m_msgraph05

import (
	msgraph "github.com/yaegashi/msgraph.go/beta"
)

/////////////////////   Membership   /////////////////////   Membership   /////////////////////   Membership   /////////////////////   Membership   /////////////////////   Membership
/*
TODO: Return all members from a group
*/

type Tmp struct {
	OdataContext string         `json:"@odata.context"`
	Value        []msgraph.User `json:"value"`
}

func (p *AADAction) GetAllMembersOfAGroup(azureGroupID string) ([]msgraph.User, error) {
	r := p.graphClient.Groups().ID(azureGroupID).Members().Request()
	var out Tmp
	err := r.JSONRequest(p.ctx, "GET", "", nil, &out)
	if err != nil {
		return nil, err
	}
	return out.Value, nil
}

/*
Add a User to a Group
*/
func (p *AADAction) AddMemberToAGroup(azureGroupId string, azureUserID string) error {
	reqObj := map[string]interface{}{
		"@odata.id": p.graphClient.DirectoryObjects().ID(azureUserID).Request().URL(),
	}
	r := p.graphClient.Groups().ID(azureGroupId).Members().Request()
	err := r.JSONRequest(p.ctx, "POST", "/$ref", reqObj, nil)
	if err != nil {
		return err
	}
	return nil
}

/*
Delete a Member from a Group
*/
func (p *AADAction) DeleteMemberFromAGroup(azureGroupID string, azureUserID string) error {
	r := p.graphClient.Groups().ID(azureGroupID).Members().ID(azureUserID).Request()
	err := r.JSONRequest(p.ctx, "DELETE", "/$ref", nil, nil)
	if err != nil {
		return err
	}
	return nil
}
