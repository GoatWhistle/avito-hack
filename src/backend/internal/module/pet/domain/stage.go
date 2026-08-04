package domain

type Stage string

const (
	StageEgg    Stage = "egg"
	StageBaby   Stage = "baby"
	StageTeen   Stage = "teen"
	StageAdult  Stage = "adult"
	StageLegend Stage = "legend"
)

func (s Stage) Valid() bool {
	switch s {
	case StageEgg, StageBaby, StageTeen, StageAdult, StageLegend:
		return true
	default:
		return false
	}
}
