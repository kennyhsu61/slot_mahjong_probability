package game

import (
	"math/rand"
)

func getCalcData(gameInput *Input, spinInfo *SpinInfo) *CalcData {
	calculator := &CalcData{
		Bet:       gameInput.BetAmount * gameInput.BetLevel,
		TotalBet:  gameInput.BetAmount * gameInput.BetLevel * Config.GameLine * Config.FeatureBuyInfo[gameInput.FeatureBuyIndex].FeatureBuyLevel,
		SpinState: spinInfo.SpinState,
	}
	calculator.StripTable = Config.StripTable[calculator.SpinState]
	calculator.getScreenMultiplier(gameInput, spinInfo)
	calculator.getCalcRNG(spinInfo)
	calculator.getCalcScreen(gameInput, spinInfo)

	return calculator
}

func (calculator *CalcData) getScreenMultiplier(gameInput *Input, spinInfo *SpinInfo) {
	//整場倍數 - 預設邏輯:隨場次增加
	spinInfo.ScreenMultiArray = Config.GameRule.ScreenMultiArray[GetSpinStateType(spinInfo.SpinState)][gameInput.Selection]
	if spinInfo.CascadingRound > 3 {
		calculator.ScreenMulti = spinInfo.ScreenMultiArray[len(spinInfo.ScreenMultiArray)-1]
	} else {
		calculator.ScreenMulti = spinInfo.ScreenMultiArray[spinInfo.CascadingRound-1]
	}
}

func (calculator *CalcData) getCalcRNG(spinInfo *SpinInfo) {
	if spinInfo.ScreenInput == nil {
		for i := 0; i < len(Config.ScreenSize); i++ {
			rndlen := len(calculator.StripTable[i])
			rndnum := rand.Intn(rndlen)
			spinInfo.RNG = append(spinInfo.RNG, rndnum)
		}
		//Debug 用，指定滾輪位置
		//spinInfo.RNG = []int{27, 54, 64, 100, 11}
	}
	calculator.RNG = spinInfo.RNG
}

func (calculator *CalcData) getCalcScreen(gameInput *Input, spinInfo *SpinInfo) {
	if spinInfo.ScreenInput == nil {
		spinInfo.Screen = calculator.genScreen(gameInput)
		calculator.Screen = make([][]SymbolID, len(spinInfo.Screen))
		for i := 0; i < len(spinInfo.Screen); i++ {
			calculator.Screen[i] = make([]SymbolID, len(spinInfo.Screen[i]))
			copy(calculator.Screen[i], spinInfo.Screen[i])
		}
	} else {
		spinInfo.Screen = make([][]SymbolID, len(spinInfo.ScreenInput))
		calculator.Screen = make([][]SymbolID, len(spinInfo.ScreenInput))
		for i := 0; i < len(spinInfo.ScreenInput); i++ {
			//原始盤面 - 由外部指定
			spinInfo.Screen[i] = make([]SymbolID, len(spinInfo.ScreenInput[i]))
			copy(spinInfo.Screen[i], spinInfo.ScreenInput[i])

			//計算盤面(原始) - 由外部指定
			calculator.Screen[i] = make([]SymbolID, len(spinInfo.ScreenInput[i]))
			copy(calculator.Screen[i], spinInfo.ScreenInput[i])
		}
	}
}

func (calculator *CalcData) genScreen(gameInput *Input) [][]SymbolID {
	var screen [][]SymbolID

	for i := 0; i < len(Config.ScreenSize); i++ {
		//開始每列取值
		var iarray []SymbolID
		reel := calculator.StripTable[i]
		iarray = copyReel(reel, calculator.RNG[i], Config.ScreenSize[i])
		screen = append(screen, iarray)
	}

	//買特色強塞FG進盤面<預留功能尚未實作>
	if gameInput.FeatureBuyIndex != "" && GetSpinStateType(calculator.SpinState) == ESST_NG {
		//骰買特色觸發FG顆數
		featureBuyHitFGNum := getOption(Config.GameOption[calculator.SpinState][0])
		for currentReel := 0; currentReel < featureBuyHitFGNum; currentReel++ {
			striplen := len(screen[currentReel])
			rndFGnum := rand.Intn(striplen)
			screen[currentReel][rndFGnum] = E_FG
		}
	}

	return screen
}

// 複製滾輪
func copyReel(source []SymbolID, index, length int) (target []SymbolID) {
	sourceSize := len(source)
	var temp []SymbolID
	if index+length >= len(source) {
		temp = append(source[index:sourceSize], source[0:index+length-sourceSize]...)
	} else {
		temp = source[index : index+length]
	}
	target = make([]SymbolID, len(temp))
	copy(target, temp)
	return
}
