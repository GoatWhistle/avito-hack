package domain

import (
	"math"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	dailyLoginXP          = 1
	searchSubscriptionXP  = 3
	favoriteXP            = 1
	newDialogueXP         = 2
	quickReplyXP          = 3
	qualityListingXP      = 2
	textReviewXP          = 1
	photoReviewXP         = 2
	videoReviewXP         = 4
	maxSearchesPerWeek    = 3
	maxFavoritesPerDay    = 5
	qualityDescriptionLen = 200
	quickReplyLimit       = time.Hour
	videoReviewMinLength  = 30 * time.Second
	daytimeStartHour      = 8
	daytimeEndHour        = 23
)

type LimitedAction struct {
	SubjectID     uuid.UUID
	Unique        bool
	RewardedCount int
	Now           time.Time
}

type QualityListing struct {
	HasVideo    bool
	Description string
	Now         time.Time
}

type Review struct {
	ConfirmedDeal bool
	HasAttachment bool
	VideoDuration time.Duration
	Now           time.Time
}

func (p *Pet) DailyCheckIn(now time.Time) (Progress, error) {
	today := day(now)
	if p.lastCheckInDate != nil && day(*p.lastCheckInDate).Equal(today) {
		return Progress{}, ErrDuplicateAction
	}

	if p.lastCheckInDate != nil && day(*p.lastCheckInDate).Equal(today.AddDate(0, 0, -1)) {
		p.streakDays++
	} else {
		p.streakDays = 1
	}
	p.lastCheckInDate = &today

	amount := dailyLoginXP
	if p.streakDays >= 7 {
		amount = int(math.Round(float64(dailyLoginXP) * 1.5))
	}

	return p.applyAward(amount, now), nil
}

func (p *Pet) RewardSearchSubscription(action LimitedAction) (Progress, error) {
	return p.rewardLimited(action, maxSearchesPerWeek, searchSubscriptionXP)
}

func (p *Pet) RewardFavorite(action LimitedAction) (Progress, error) {
	return p.rewardLimited(action, maxFavoritesPerDay, favoriteXP)
}

func (p *Pet) RewardNewDialogue(action LimitedAction) (Progress, error) {
	if action.SubjectID == uuid.Nil {
		return Progress{}, ErrInvalidAction
	}
	if !action.Unique {
		return Progress{}, ErrDuplicateAction
	}

	return p.applyAward(newDialogueXP, action.Now), nil
}

func (p *Pet) RewardQuickReply(delay time.Duration, now time.Time) (Progress, error) {
	if delay < 0 || delay > quickReplyLimit || !isDaytime(now) {
		return Progress{}, ErrConditionNotMet
	}

	return p.applyAward(quickReplyXP, now), nil
}

func (p *Pet) RewardQualityListing(listing QualityListing) (Progress, error) {
	if !listing.HasVideo || utf8.RuneCountInString(listing.Description) <= qualityDescriptionLen {
		return Progress{}, ErrConditionNotMet
	}

	return p.applyAward(qualityListingXP, listing.Now), nil
}

func (p *Pet) RewardTextReview(review Review) (Progress, error) {
	if !review.ConfirmedDeal {
		return Progress{}, ErrConditionNotMet
	}

	return p.applyAward(textReviewXP, review.Now), nil
}

func (p *Pet) RewardPhotoReview(review Review) (Progress, error) {
	if !review.ConfirmedDeal || !review.HasAttachment {
		return Progress{}, ErrConditionNotMet
	}

	return p.applyAward(photoReviewXP, review.Now), nil
}

func (p *Pet) RewardVideoReview(review Review) (Progress, error) {
	if !review.ConfirmedDeal || review.VideoDuration <= videoReviewMinLength {
		return Progress{}, ErrConditionNotMet
	}

	return p.applyAward(videoReviewXP, review.Now), nil
}

func (p *Pet) rewardLimited(action LimitedAction, limit, amount int) (Progress, error) {
	if action.SubjectID == uuid.Nil {
		return Progress{}, ErrInvalidAction
	}
	if action.RewardedCount < 0 {
		return Progress{}, ErrInvalidAction
	}
	if !action.Unique {
		return Progress{}, ErrDuplicateAction
	}
	if action.RewardedCount >= limit {
		return Progress{}, ErrLimitReached
	}

	return p.applyAward(amount, action.Now), nil
}

func (p *Pet) applyAward(amount int, now time.Time) Progress {
	p.updatedAt = now
	return p.awardXP(amount)
}

func day(value time.Time) time.Time {
	year, month, date := value.Date()
	return time.Date(year, month, date, 0, 0, 0, 0, value.Location())
}

func isDaytime(now time.Time) bool {
	return now.Hour() >= daytimeStartHour && now.Hour() < daytimeEndHour
}
