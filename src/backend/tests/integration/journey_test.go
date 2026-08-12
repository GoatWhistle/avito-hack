//go:build integration

package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	itemapp "github.com/avito-hack/backend/internal/module/item/app"
	itemdomain "github.com/avito-hack/backend/internal/module/item/domain"
	petdomain "github.com/avito-hack/backend/internal/module/pet/domain"
	userapp "github.com/avito-hack/backend/internal/module/user/app"
	"github.com/avito-hack/backend/internal/shared/auth"
)

const (
	oneDay          = 24 * time.Hour
	growthDays      = 8
	favoritesPerDay = 5
)

const qualityDescription = "Полностью рабочая консоль с двумя джойстиками и набором фирменных картриджей. " +
	"Приставка бережно хранилась в коробке, следов вскрытия и ремонта нет, все разъёмы чистые. " +
	"В комплекте оригинальный блок питания, кабель для телевизора и инструкция на русском языке. " +
	"Продаю в связи с переездом, готов показать работу устройства при встрече."

func TestFullUserJourney(t *testing.T) {
	ctx := context.Background()
	e := newEnv(t)

	email := uniqueEmail()

	registered, err := e.register.Handle(ctx, userapp.RegisterUserCommand{
		Email:    email,
		Password: "Str0ngPassw0rd!",
		FullName: "Integration User",
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, registered.User.ID())
	assert.Equal(t, email, registered.User.Email().String())

	userID := registered.User.ID()

	loggedIn, err := e.login.Handle(ctx, userapp.LoginUserCommand{
		Email:    email,
		Password: "Str0ngPassw0rd!",
	})
	require.NoError(t, err)
	require.NotEmpty(t, loggedIn.Token)

	actor, err := e.tokens.Parse(loggedIn.Token)
	require.NoError(t, err)
	assert.Equal(t, userID, actor.ID)

	item, err := e.create.Handle(ctx, itemapp.CreateItemCommand{
		OwnerID:     userID,
		Title:       "Ретро-приставка в отличном состоянии",
		Description: qualityDescription,
		PriceKopeks: 1_500_000,
	})
	require.NoError(t, err)
	assert.Equal(t, itemdomain.StatusDraft, item.Status())

	submitted, err := e.status.Handle(ctx, itemapp.ChangeStatusCommand{
		ItemID: item.ID(),
		Actor:  auth.Actor{ID: userID, Role: auth.RoleUser},
		Action: itemapp.ActionSubmit,
	})
	require.NoError(t, err)
	require.Equal(t, itemdomain.StatusModeration, submitted.Status())

	require.NoError(t, submitted.Publish(e.clock.Now()))
	require.NoError(t, e.items.Save(ctx, submitted))
	published := submitted

	awarded, err := e.pets.RewardQualityListing(ctx, userID, item.ID(), petdomain.QualityListing{
		HasPhoto:    true,
		HasPrice:    true,
		Description: published.Description(),
	})
	require.NoError(t, err)
	require.NotNil(t, awarded.Pet)
	require.Positive(t, awarded.Progress.XPGranted, "publishing a quality listing must grant XP")

	xpAfterListing := awarded.Pet.XP()
	assert.Equal(t, awarded.Progress.XPGranted, xpAfterListing)

	levelAfterListing := awarded.Pet.Level()

	grown := growToRewardLevel(t, e, userID)
	require.Greater(t, grown.XP(), xpAfterListing, "repeated actions must accumulate XP")
	require.Greater(t, grown.Level(), levelAfterListing, "accumulated XP must raise the level")

	issued, err := e.rewards.GrantEligible(ctx, userID)
	require.NoError(t, err)
	require.NotEmpty(t, issued, "a pet at level %d must be eligible for at least one reward", grown.Level())

	rewardID := issued[0].ID()

	code, err := e.rewards.Activate(ctx, userID, rewardID)
	require.NoError(t, err)
	require.NotEmpty(t, code)

	_, err = e.rewards.Activate(ctx, userID, rewardID)
	require.ErrorIs(t, err, petdomain.ErrRewardAlreadyActivated,
		"activating an already activated reward must be rejected")

	mine, err := e.rewards.Mine(ctx, userID)
	require.NoError(t, err)
	require.NotEmpty(t, mine)
}

func growToRewardLevel(t *testing.T, e *env, userID uuid.UUID) *petdomain.Pet {
	t.Helper()

	ctx := context.Background()

	var pet *petdomain.Pet

	for day := range growthDays {
		if day > 0 {
			e.clock.Advance(oneDay)
		}

		if result, err := e.pets.CheckIn(ctx, userID); err == nil {
			pet = result.Pet
		} else {
			require.ErrorIs(t, err, petdomain.ErrLimitReached)
		}

		for range favoritesPerDay {
			if result, err := e.pets.AddFavorite(ctx, userID, uuid.New()); err == nil {
				pet = result.Pet

				continue
			} else {
				require.ErrorIs(t, err, petdomain.ErrLimitReached)
			}

			break
		}
	}

	require.NotNil(t, pet, "growth loop must have awarded XP at least once")

	return pet
}
