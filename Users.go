package itswizard_m_msgraph05

import (
	"errors"
	"fmt"
	msgraph "github.com/yaegashi/msgraph.go/beta"
	P "github.com/yaegashi/msgraph.go/ptr"
)

/////////////////////  PERSONS   /////////////////////  PERSONS   /////////////////////    PERSONS   /////////////////////    PERSONS   /////////////////////    PERSONS   /////////////////////
/*
Response all Users
*/
func (p *AADAction) GetAllUsers() ([]msgraph.User, error) {
	r := p.graphClient.Users().Request()
	users, err := r.Get(p.ctx)
	return users, err
}

/*
Get one user with username
*/
func (p *AADAction) GetUserWithUsername(username string) (msgraph.User, error) {
	r := p.graphClient.Users().Request()
	r.Filter(fmt.Sprintf("userPrincipalName eq '%s'", username))
	users, err := r.Get(p.ctx)
	if err != nil {
		return msgraph.User{}, err
	}
	if len(users) < 1 {
		return msgraph.User{}, errors.New("User does not exist")
	}
	return users[0], nil
}

/*
Get one User with AzureUserID
*/
func (p *AADAction) GetUserWithID(azureUserID string) (user msgraph.User, err error) {
	r := p.graphClient.Users().ID(azureUserID).Request()
	tmp, err := r.Get(p.ctx)
	if err != nil {
		return user, err
	}
	return *tmp, nil
}

/*
Create a new user
*/
func (p *AADAction) CreateUser(firstname, lastname, profile, password, username, mailNick string) (azureUserId string, err error) {

	newUser := msgraph.User{
		DisplayName:  P.String(fmt.Sprint(firstname, " ", lastname)),
		GivenName:    P.String(firstname),
		Surname:      P.String(lastname),
		MailNickname: P.String(mailNick),
		PasswordProfile: &msgraph.PasswordProfile{
			ForceChangePasswordNextSignIn: P.Bool(true),
			Password:                      P.String(password),
		},
		UserPrincipalName: P.String(fmt.Sprint(username)),
		AccountEnabled:    P.Bool(true),
		UsageLocation:     P.String("DE"),
		JobTitle:          P.String(profile),
	}

	u, err := p.graphClient.Users().Request().Add(p.ctx, &newUser)
	if err != nil {
		return azureUserId, err
	}
	return *u.ID, nil
}

/*
Add a License to the user
*/
func (p *AADAction) AddLicense(azureUserID string, sKUId string) error {

	guid := msgraph.UUID(sKUId)
	reqObj := map[string]interface{}{
		"addLicenses":    []msgraph.AssignedLicense{msgraph.AssignedLicense{SKUID: &guid}},
		"removeLicenses": []msgraph.UUID{},
	}
	err := p.graphClient.Users().ID(azureUserID).Request().JSONRequest(p.ctx, "POST", "/assignLicense", reqObj, nil)
	if err != nil {
		return err
	}
	return nil
}

/*
Delete a User with username
*/
func (p *AADAction) DeleteUser(azureUserId string) error {
	return p.graphClient.Users().ID(azureUserId).Request().Delete(p.ctx)
}

/*
Update a User
*/
func (p *AADAction) UpdateUser(azureUserId, firstname, lastname, profile, syncID, username, domain string) error {
	//PATCH /users/{id | userPrincipalName}
	reqObj := map[string]interface{}{
		"DisplayName":       fmt.Sprint(firstname, " ", lastname),
		"GivenName":         firstname,
		"Surname":           lastname,
		"JobTitle":          profile,
		"EmployeeID":        syncID,
		"userPrincipalName": fmt.Sprint(username, "@", domain),
	}

	err := p.graphClient.Users().ID(azureUserId).Request().JSONRequest(p.ctx, "PATCH", "", reqObj, nil)
	if err != nil {
		return err
	}
	return nil
}

/*
Update a User
*/
func (p *AADAction) UpdateUserPassword(username, newPassword string) error {
	//PATCH /users/{id | userPrincipalName}
	pw := msgraph.PasswordProfile{
		ForceChangePasswordNextSignIn: P.Bool(false),
		Password:                      P.String(newPassword),
	}

	reqObj := map[string]interface{}{
		"passwordProfile": pw,
	}

	err := p.graphClient.Users().ID(username).Request().JSONRequest(p.ctx, "PATCH", "", reqObj, nil)
	if err != nil {
		return err
	}
	return nil
}
