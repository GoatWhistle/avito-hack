//go:build integration

package integration

import (
	"net/http"
	"testing"
)

type user struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type session struct {
	Token string `json:"token"`
	User  user   `json:"user"`
}

type item struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Price  int64  `json:"price"`
	Status string `json:"status"`
}

type pet struct {
	Level       int    `json:"level"`
	XP          int    `json:"xp"`
	NextLevelXP int    `json:"next_level_xp"`
	Stage       string `json:"stage"`
}

type progress struct {
	Level         int  `json:"level"`
	XP            int  `json:"xp"`
	XPToNextLevel int  `json:"xp_to_next_level"`
	IsMaxLevel    bool `json:"is_max_level"`
}

type userReward struct {
	RewardID string `json:"reward_id"`
	Kind     string `json:"kind"`
	Status   string `json:"status"`
	Code     string `json:"code"`
}

type userRewardList struct {
	Items []userReward `json:"items"`
}

type activation struct {
	RewardID string `json:"reward_id"`
	Code     string `json:"code"`
	Status   string `json:"status"`
}

const (
	password        = "integration1234"
	publishedItemXP = 50
)

func TestJourneyPublishEarnLevelUpActivateReward(t *testing.T) {
	c := requireStack(t)

	email := uniqueEmail("journey")
	registered := register(t, c, email)

	logIn(t, c, email)
	assertProfile(t, c, registered.ID, email)

	before := currentPet(t, c)

	created := createItem(t, c, "Integration Bicycle")
	if created.Status != "draft" {
		t.Fatalf("new item status: got %q want %q", created.Status, "draft")
	}

	published := changeStatus(t, c, created.ID, "publish")
	if published.Status != "published" {
		t.Fatalf("published item status: got %q want %q", published.Status, "published")
	}

	after := awaitXPGain(t, c, before.XP)
	if after.Level <= before.Level {
		t.Fatalf("level did not increase: before %d after %d", before.Level, after.Level)
	}

	assertProgressMatchesPet(t, c, after)

	reward := awaitActivatableReward(t, c)

	activated := activateReward(t, c, reward.RewardID, http.StatusOK)
	if activated.Code == "" {
		t.Fatal("activation returned an empty promo code")
	}
	if activated.Status != "activated" {
		t.Fatalf("activation status: got %q want %q", activated.Status, "activated")
	}

	assertReactivationRejected(t, c, reward.RewardID)
	assertCodeIsStable(t, c, reward.RewardID, activated.Code)
}

func register(t *testing.T, c *client, email string) user {
	t.Helper()

	var opened session
	c.mustDo(http.MethodPost, "/api/v1/auth/register", map[string]string{
		"email":     email,
		"password":  password,
		"full_name": "Integration User",
	}, http.StatusCreated, &opened)

	created := opened.User

	if opened.Token == "" {
		t.Fatal("register returned an empty token")
	}
	if created.ID == "" {
		t.Fatal("register returned an empty user id")
	}
	if created.Role != "user" {
		t.Fatalf("register role: got %q want %q", created.Role, "user")
	}

	status, payload := c.do(http.MethodPost, "/api/v1/auth/register", map[string]string{
		"email":     email,
		"password":  password,
		"full_name": "Integration User",
	})
	if status != http.StatusConflict {
		t.Fatalf("duplicate register: got %d want %d, body: %s", status, http.StatusConflict, truncate(payload))
	}
	if code := decodeError(t, payload).Error.Code; code != "conflict" {
		t.Fatalf("duplicate register code: got %q want %q", code, "conflict")
	}

	return created
}

func logIn(t *testing.T, c *client, email string) {
	t.Helper()

	var issued session
	c.mustDo(http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    email,
		"password": password,
	}, http.StatusOK, &issued)

	if issued.Token == "" {
		t.Fatal("login returned an empty token")
	}
	c.token = issued.Token
}

func assertProfile(t *testing.T, c *client, wantID, wantEmail string) {
	t.Helper()

	var me user
	c.mustDo(http.MethodGet, "/api/v1/users/me", nil, http.StatusOK, &me)

	if me.ID != wantID {
		t.Fatalf("profile id: got %q want %q", me.ID, wantID)
	}
	if me.Email != wantEmail {
		t.Fatalf("profile email: got %q want %q", me.Email, wantEmail)
	}
}

func createItem(t *testing.T, c *client, title string) item {
	t.Helper()

	var created item
	c.mustDo(http.MethodPost, "/api/v1/items", map[string]any{
		"title":       title,
		"description": "Created by the integration journey test.",
		"price":       1500000,
	}, http.StatusCreated, &created)

	if created.ID == "" {
		t.Fatal("create item returned an empty id")
	}

	return created
}

func changeStatus(t *testing.T, c *client, itemID, action string) item {
	t.Helper()

	var updated item
	c.mustDo(http.MethodPost, "/api/v1/items/"+itemID+"/status",
		map[string]string{"action": action}, http.StatusOK, &updated)

	return updated
}

func currentPet(t *testing.T, c *client) pet {
	t.Helper()

	var state pet
	c.mustDo(http.MethodGet, "/api/v1/pet", nil, http.StatusOK, &state)

	return state
}
