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
	qualityDescriptionLen = 200
	qualityVideoBonusXP   = 4
	quickReplyLimit       = time.Hour
	videoReviewMinLength  = 30 * time.Second
	daytimeStartHour      = 8
	daytimeEndHour        = 23
	moscowOffsetSeconds   = 3 * 60 * 60
	streakBonusDays       = 7
	streakBonusFactor     = 1.5
)

const (
	MaxSearchesPerWeek = 3
	MaxFavoritesPerDay = 5
)

var moscowZone = time.FixedZone("MSK", moscowOffsetSeconds)

type LimitedAction struct {
	SubjectID     uuid.UUID
	Unique        bool
	RewardedCount int
	Now           time.Time
}

type QualityListing struct {
	HasPhoto    bool
	HasPrice    bool
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
	outcome, err := p.CheckIn(now)
	if err != nil {
		return Progress{}, err
	}

	return outcome.Progress, nil
}

type CheckInResult struct {
	Progress Progress
	Streak   StreakOutcome
}

func (p *Pet) CheckIn(now time.Time) (CheckInResult, error) {
	today := day(now)
	if p.lastCheckInDate != nil && day(*p.lastCheckInDate).Equal(today) {
		return CheckInResult{}, ErrDuplicateAction
	}

	p.ApplyDecay(now)
	streak := p.advanceStreak(now)

	amount := dailyLoginXP
	if p.streakDays >= streakBonusDays {
		amount = int(math.Round(float64(dailyLoginXP) * streakBonusFactor))
	}
	amount += streak.MilestoneBonus

	p.satiety = shiftParameter(p.satiety, checkInSatietyGain)

	return CheckInResult{Progress: p.applyAward(amount, now), Streak: streak}, nil
}

func (p *Pet) RewardSearchSubscription(action LimitedAction) (Progress, error) {
	return p.rewardLimited(action, MaxSearchesPerWeek, searchSubscriptionXP)
}

func (p *Pet) RewardFavorite(action LimitedAction) (Progress, error) {
	return p.rewardLimited(action, MaxFavoritesPerDay, favoriteXP)
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
	if !listing.HasPhoto || !listing.HasPrice || !hasQualityDescription(listing.Description) {
		return Progress{}, ErrConditionNotMet
	}

	amount := qualityListingXP
	if listing.HasVideo {
		amount = qualityVideoBonusXP
	}

	return p.applyAward(amount, listing.Now), nil
}

func hasQualityDescription(description string) bool {
	return utf8.RuneCountInString(description) > qualityDescriptionLen
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
	p.ApplyDecay(now)
	p.updatedAt = now

	return p.awardXP(p.effectiveAward(amount))
}

func day(value time.Time) time.Time {
	local := value.In(moscowZone)
	year, month, date := local.Date()

	return time.Date(year, month, date, 0, 0, 0, 0, moscowZone)
}

func isDaytime(now time.Time) bool {
	hour := now.In(moscowZone).Hour()

	return hour >= daytimeStartHour && hour < daytimeEndHour
}
