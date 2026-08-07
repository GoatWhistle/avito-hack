//go:build integration

package integration

import (
	"net/http"
	"testing"
	"time"
)

const (
	asyncTimeout  = 30 * time.Second
	asyncInterval = 500 * time.Millisecond
)

func awaitXPGain(t *testing.T, c *client, baselineXP int) pet {
	t.Helper()

	var last pet
	deadline := time.Now().Add(asyncTimeout)
	for time.Now().Before(deadline) {
		last = currentPet(t, c)
		if last.XP >= baselineXP+publishedItemXP {
			return last
		}
		time.Sleep(asyncInterval)
	}

	t.Fatalf("pet did not gain %d xp within %s: baseline %d, last %d",
		publishedItemXP, asyncTimeout, baselineXP, last.XP)

	return last
}

func awaitActivatableReward(t *testing.T, c *client) userReward {
	t.Helper()

	deadline := time.Now().Add(asyncTimeout)
	for time.Now().Before(deadline) {
		for _, granted := range myRewards(t, c) {
			if granted.Kind == "promo" && granted.Status == "granted" {
				return granted
			}
		}
		time.Sleep(asyncInterval)
	}

	t.Fatalf("no activatable promo reward was granted within %s", asyncTimeout)

	return userReward{}
}

func myRewards(t *testing.T, c *client) []userReward {
	t.Helper()

	c.mustDo(http.MethodGet, "/api/v1/pet", nil, http.StatusOK, nil)

	var mine userRewardList
	c.mustDo(http.MethodGet, "/api/v1/rewards/my", nil, http.StatusOK, &mine)

	return mine.Items
}

func activateReward(t *testing.T, c *client, rewardID string, want int) activation {
	t.Helper()

	var result activation
	c.mustDo(http.MethodPost, "/api/v1/rewards/"+rewardID+"/activate", nil, want, &result)

	return result
}

func assertReactivationRejected(t *testing.T, c *client, rewardID string) {
	t.Helper()

	status, payload := c.do(http.MethodPost, "/api/v1/rewards/"+rewardID+"/activate", nil)
	if status != http.StatusConflict {
		t.Fatalf("second activation: got %d want %d, body: %s", status, http.StatusConflict, truncate(payload))
	}
	if code := decodeError(t, payload).Error.Code; code != "conflict" {
		t.Fatalf("second activation code: got %q want %q", code, "conflict")
	}
}

func assertCodeIsStable(t *testing.T, c *client, rewardID, wantCode string) {
	t.Helper()

	for _, granted := range myRewards(t, c) {
		if granted.RewardID != rewardID {
			continue
		}
		if granted.Status != "activated" {
			t.Fatalf("reward %s status: got %q want %q", rewardID, granted.Status, "activated")
		}
		if granted.Code != wantCode {
			t.Fatalf("reward %s code: got %q want %q", rewardID, granted.Code, wantCode)
		}

		return
	}

	t.Fatalf("activated reward %s is missing from /rewards/my", rewardID)
}

func assertProgressMatchesPet(t *testing.T, c *client, state pet) {
	t.Helper()

	var bar progress
	c.mustDo(http.MethodGet, "/api/v1/progress", nil, http.StatusOK, &bar)

	if bar.Level != state.Level {
		t.Fatalf("progress level: got %d want %d", bar.Level, state.Level)
	}
	if bar.IsMaxLevel {
		return
	}
	if want := state.NextLevelXP - state.XP; bar.XPToNextLevel != want {
		t.Fatalf("xp_to_next_level: got %d want %d", bar.XPToNextLevel, want)
	}
}

func TestCheckInIsIdempotentWithinADay(t *testing.T) {
	c := requireStack(t)

	email := uniqueEmail("checkin")
	register(t, c, email)
	logIn(t, c, email)

	status, payload := c.do(http.MethodPost, "/api/v1/checkin", nil)
	if status != http.StatusOK && status != http.StatusConflict {
		t.Fatalf("first check-in: got %d want 200 or 409, body: %s", status, truncate(payload))
	}

	status, payload = c.do(http.MethodPost, "/api/v1/checkin", nil)
	if status != http.StatusConflict {
		t.Fatalf("repeat check-in: got %d want %d, body: %s", status, http.StatusConflict, truncate(payload))
	}
	if code := decodeError(t, payload).Error.Code; code != "conflict" {
		t.Fatalf("repeat check-in code: got %q want %q", code, "conflict")
	}
}

func TestUnauthorizedAccessIsRejected(t *testing.T) {
	c := requireStack(t)

	for _, path := range []string{"/api/v1/users/me", "/api/v1/pet", "/api/v1/progress", "/api/v1/rewards/my"} {
		status, payload := c.do(http.MethodGet, path, nil)
		if status != http.StatusUnauthorized {
			t.Fatalf("anonymous GET %s: got %d want %d, body: %s",
				path, status, http.StatusUnauthorized, truncate(payload))
		}
	}
}
