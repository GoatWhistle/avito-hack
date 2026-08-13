package domain

type Symbol string

const (
	SymbolBicycle    Symbol = "bicycle"
	SymbolSmartphone Symbol = "smartphone"
	SymbolSofa       Symbol = "sofa"
	SymbolSneakers   Symbol = "sneakers"
	SymbolDelivery   Symbol = "delivery"
	SymbolPromotion  Symbol = "promotion"
)

const (
	ScopeTypeCategory   = "category"
	CategoryElectronics = "electronics"
	CategoryFurniture   = "furniture"
	CategoryClothes     = "clothes"
	CategorySport       = "sport"
)

type Prize struct {
	ID           string
	Symbol       Symbol
	Title        string
	Description  string
	BenefitType  string
	BenefitValue int
	ScopeType    string
	ScopeValue   string
	Weight       int
}

var prizes = []Prize{
	{ID: "weekly_bicycle_5", Symbol: SymbolBicycle, Title: "Скидка 5% на спорт и отдых",
		Description: "Скидка на товар из категории спорта и отдыха", BenefitType: "percent_discount",
		BenefitValue: 5, ScopeType: ScopeTypeCategory, ScopeValue: CategorySport, Weight: 15},
	{ID: "weekly_smartphone_3", Symbol: SymbolSmartphone, Title: "Скидка 3% на электронику",
		Description: "Скидка на товар из категории электроники", BenefitType: "percent_discount",
		BenefitValue: 3, ScopeType: ScopeTypeCategory, ScopeValue: CategoryElectronics, Weight: 10},
	{ID: "weekly_furniture_5", Symbol: SymbolSofa, Title: "Скидка 5% на мебель",
		Description: "Скидка на товар из категории мебели", BenefitType: "percent_discount",
		BenefitValue: 5, ScopeType: ScopeTypeCategory, ScopeValue: CategoryFurniture, Weight: 10},
	{ID: "weekly_clothes_7", Symbol: SymbolSneakers, Title: "Скидка 7% на одежду и обувь",
		Description: "Скидка на товар из категории одежды и обуви", BenefitType: "percent_discount",
		BenefitValue: 7, ScopeType: ScopeTypeCategory, ScopeValue: CategoryClothes, Weight: 7},
	{ID: "weekly_delivery_500", Symbol: SymbolDelivery, Title: "Бесплатная доставка до 500 ₽",
		Description: "Компенсация доставки в пределах 500 ₽", BenefitType: "delivery_discount_kopeks",
		BenefitValue: 50000, ScopeType: "delivery", ScopeValue: "eligible", Weight: 5},
	{ID: "weekly_promotion_24h", Symbol: SymbolPromotion, Title: "Продвижение на 24 часа",
		Description: "Бесплатное продвижение одного объявления", BenefitType: "listing_promotion_hours",
		BenefitValue: 24, ScopeType: "own_listing", ScopeValue: "any", Weight: 3},
}

func Prizes() []Prize {
	out := make([]Prize, len(prizes))
	copy(out, prizes)

	return out
}

func PrizeByID(id string) (Prize, bool) {
	for _, prize := range prizes {
		if prize.ID == id {
			return prize, true
		}
	}

	return Prize{}, false
}
