package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	initialLevel      = 1
	initialSatiety    = 70
	initialHappiness  = 70
	maxParameterValue = 100
	petHappinessGain  = 5
)

type Pet struct {
	id              uuid.UUID
	userID          uuid.UUID
	name            string
	stage           Stage
	level           int
	xp              int
	nextLevelXP     int
	satiety         int
	happiness       int
	streakDays      int
	lastCheckInDate *time.Time
	lastDecayTime   time.Time
	updatedAt       time.Time
}

func New(userID uuid.UUID, now time.Time) *Pet {
	return &Pet{
		id:            uuid.New(),
		userID:        userID,
		stage:         StageEgg,
		level:         initialLevel,
		nextLevelXP:   nextThreshold(initialLevel),
		satiety:       initialSatiety,
		happiness:     initialHappiness,
		lastDecayTime: now,
		updatedAt:     now,
	}
}

type RestoreParams struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Name            string
	Stage           Stage
	Level           int
	XP              int
	NextLevelXP     int
	Satiety         int
	Happiness       int
	StreakDays      int
	LastCheckInDate *time.Time
	LastDecayTime   time.Time
	UpdatedAt       time.Time
}

func Restore(p RestoreParams) *Pet {
	return &Pet{
		id: p.ID, userID: p.UserID, name: p.Name, stage: p.Stage, level: p.Level, xp: p.XP,
		nextLevelXP: nextThreshold(p.Level), satiety: p.Satiety, happiness: p.Happiness,
		streakDays: p.StreakDays, lastCheckInDate: cloneTime(p.LastCheckInDate),
		lastDecayTime: p.LastDecayTime, updatedAt: p.UpdatedAt,
	}
}

func (p *Pet) Pet(now time.Time) {
	p.happiness = min(p.happiness+petHappinessGain, maxParameterValue)
	p.updatedAt = now
}

func (p *Pet) ID() uuid.UUID               { return p.id }
func (p *Pet) UserID() uuid.UUID           { return p.userID }
func (p *Pet) Name() string                { return p.name }
func (p *Pet) Stage() Stage                { return p.stage }
func (p *Pet) Level() int                  { return p.level }
func (p *Pet) XP() int                     { return p.xp }
func (p *Pet) NextLevelXP() int            { return p.nextLevelXP }
func (p *Pet) Satiety() int                { return p.satiety }
func (p *Pet) Happiness() int              { return p.happiness }
func (p *Pet) StreakDays() int             { return p.streakDays }
func (p *Pet) LastCheckInDate() *time.Time { return cloneTime(p.lastCheckInDate) }
func (p *Pet) LastDecayTime() time.Time    { return p.lastDecayTime }
func (p *Pet) UpdatedAt() time.Time        { return p.updatedAt }

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}
