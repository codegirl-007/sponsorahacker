package twitch

import (
	"fmt"
	"github.com/nicklaw5/helix"
)

type TwitchClient struct {
	client *helix.Client
}

type User struct {
	ID              string
	Username        string
	DisplayName     string
	Description     string
	ProfileImageURL string
	Following       []User
	Subscriptions   []Subscription
}

type Subscription struct {
	UserID      string
	UserName    string
	ChannelID   string
	ChannelName string
	Tier        string // Subscription tier (e.g., "1000", "2000", etc.)
}

type EventSubscription struct {
	UserID    string
	UserLogin string
	isGift    bool
	tier      int
}

func NewTwitchClient(clientID string, clientSecret string) (*TwitchClient, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})

	if err != nil {
		return nil, err
	}
	return &TwitchClient{client: client}, nil
}

func (t *TwitchClient) GetUser(username string) (*User, error) {
	// get the user
	user, err := t.client.GetUsers(&helix.UsersParams{
		Logins: []string{username},
	})

	if err != nil {
		fmt.Println("error getting user")
		return nil, err
	}

	if len(user.Data.Users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	helixUser := user.Data.Users[0]

	// get the people this user follows
	followResp, err := t.client.GetUsersFollows(&helix.UsersFollowsParams{
		FromID: helixUser.ID,
	})

	followedUsers := make([]User, 0)

	for i, follow := range followResp.Data.Follows {
		followedUsers[i] = User{
			ID:       follow.ToID,
			Username: follow.ToName,
		}
	}

	// get the user subscriptions
	subscriptionResp, err := t.client.GetSubscriptions(&helix.SubscriptionsParams{
		UserID: []string{helixUser.ID},
	})

	if err != nil {
		return nil, err
	}

	subscriptions := make([]Subscription, len(subscriptionResp.Data.Subscriptions))

	for i, sub := range subscriptionResp.Data.Subscriptions {
		subscriptions[i] = Subscription{
			UserID:      sub.UserID,
			UserName:    sub.UserName,
			ChannelID:   sub.BroadcasterID,
			ChannelName: sub.BroadcasterName,
			Tier:        sub.Tier,
		}
	}

	// build the custom user object that has all the things
	customUser := &User{
		ID:              helixUser.ID,
		Username:        helixUser.Login,
		DisplayName:     helixUser.DisplayName,
		Description:     helixUser.Description,
		ProfileImageURL: helixUser.ProfileImageURL,
		Following:       followedUsers,
		Subscriptions:   subscriptions,
	}

	return customUser, nil
}

func (t *TwitchClient) CreateEventSubscription(broadcasterID string) error {
	_, err := t.client.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type:      "channel.subscribe",
		Version:   "1",
		Condition: helix.EventSubCondition{BroadcasterUserID: broadcasterID},
		Transport: helix.EventSubTransport{
			Method:   "webhook",
			Callback: "http://localhost:8080/webhook",
			Secret:   "secret",
		},
	})

	if err != nil {
		fmt.Println("error creating event subscription")
		return err
	}

	return nil
}
