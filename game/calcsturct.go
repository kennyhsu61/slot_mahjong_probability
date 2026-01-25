package game

type SpinInfo struct {
	SpinState        SpinState
	Round            int
	FGTotalRound     int
	CascadingRound   int
	BGSurviveRound   int
	Stage            int
	LineTableNo      int
	OddsTableNo      int
	ScreenMultiArray []uint64
	WildMultiArray   []uint64
	RNG              []int
	Screen           [][]SymbolID
	ScreenInput      [][]SymbolID
	MoneyScreen      [][]int
	MoneyScreenInput [][]int
	ExtraScreen      [][]SymbolID
	ExtraScreenInput [][]SymbolID
	ExtraSymbol      SymbolID
	OptionSelect     []int
}

type CalcData struct {
	Bet               uint64
	TotalBet          uint64
	SpinState         SpinState
	StripTable        [][]SymbolID
	ScreenMulti       uint64
	RoundBGSurvive    int
	LightingLinkCount int
	LightingLinkState SpinState
	IsExtraWin        bool
	ScreenMultiArray  []uint64
	WildMultiArray    []uint64
	RNG               []int
	Screen            [][]SymbolID
	ScreenInput       [][]SymbolID
	MoneyScreen       [][]int
	MoneyScreenInput  [][]int
	ExtraScreen       [][]SymbolID
	ExtraScreenInput  [][]SymbolID
	ExtraSymbol       SymbolID
	OptionSelect      []int
}
