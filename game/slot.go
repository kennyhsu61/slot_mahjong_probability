package game

type Slot struct {
}

func (slot *Slot) GetResult(gameInput *Input) (*Output, error) {

	output, err, done := checkInuptError(gameInput)
	if done {
		return output, err
	}

	gameOutput := &Output{
		RTPIndex: gameInput.RTPIndex,
		TotalBet: gameInput.BetAmount * gameInput.BetLevel * Config.GameLine * Config.FeatureBuyInfo[gameInput.FeatureBuyIndex].FeatureBuyLevel,
	}

	triggerFGInfo := checkNGWinResult(gameInput, gameOutput)
	checkFGWinResult(gameInput, gameOutput, triggerFGInfo)

	//風控最終檢查
	if gameInput.RoundMaxWinX != 0 && int64(gameOutput.TotalWin/gameOutput.TotalBet) > gameInput.RoundMaxWinX {
		return slot.GetResult(gameInput)
	}
	if gameInput.AgentMaxPayout != 0 && gameOutput.TotalWin > gameInput.AgentMaxPayout {
		return slot.GetResult(gameInput)
	}

	return gameOutput, nil
}

func checkNGWinResult(gameInput *Input, gameOutput *Output) *TriggerFGInfo {

	var cascadingRound = 1
	var useSpinState SpinState
	useSpinState = SpinState(getOption(Config.GameStripOption[gameInput.RTPIndex][ESS_NG_1+SpinState(gameInput.Selection)][0]))

	spinInfo := &SpinInfo{
		SpinState:      useSpinState,
		Round:          1,
		CascadingRound: cascadingRound,
	}

	ngWinResult := singleWinResult(gameInput, gameOutput, spinInfo)
	triggerInfo := &TriggerFGInfo{
		WinType:      ngWinResult.WinType,
		FGTotalRound: ngWinResult.FGTotalRound,
	}

	var winLineCount int
	var screenInput = ngWinResult.ScreenOutput
	var ngLastRNG = ngWinResult.LastRNG
	if ngWinResult.WinType&EWT_FREEGAME != 0 {
		winLineCount = 0
	} else {
		winLineCount = len(ngWinResult.WinLine)
	}
	for winLineCount > 0 {
		cascadingRound++
		spinInfo2 := &SpinInfo{
			SpinState:      useSpinState,
			Round:          1,
			CascadingRound: cascadingRound,
			ScreenInput:    screenInput,
			RNG:            ngLastRNG,
		}
		ngWinResult2 := singleWinResult(gameInput, gameOutput, spinInfo2)
		ngLastRNG = ngWinResult2.LastRNG
		screenInput = ngWinResult2.ScreenOutput
		triggerInfo = &TriggerFGInfo{
			WinType:      ngWinResult2.WinType,
			FGTotalRound: ngWinResult2.FGTotalRound,
		}

		if ngWinResult2.WinType&EWT_FREEGAME != 0 {
			winLineCount = 0
		} else {
			winLineCount = len(ngWinResult2.WinLine)
		}
	}

	return triggerInfo
}

func checkFGWinResult(gameInput *Input, gameOutput *Output, triggerFGInfo *TriggerFGInfo) {

	if triggerFGInfo.WinType&EWT_FREEGAME != 0 {
		for j := 1; j <= triggerFGInfo.FGTotalRound; j++ {
			var fgUseSpinState SpinState

			fgUseSpinState = SpinState(getOption(Config.GameStripOption[gameInput.RTPIndex][ESS_FG_1+SpinState(gameInput.Selection)][0]))

			//產生FreeGame SpinData
			cascadingRound := 1
			spinInfo := &SpinInfo{
				SpinState:      fgUseSpinState,
				Round:          j,
				FGTotalRound:   triggerFGInfo.FGTotalRound,
				CascadingRound: cascadingRound,
			}
			fgWinResult := singleWinResult(gameInput, gameOutput, spinInfo)

			triggerFGInfo.FGTotalRound = fgWinResult.FGTotalRound
			var fgLastRNG = fgWinResult.LastRNG
			var screenInput = fgWinResult.ScreenOutput
			var winLineCount = len(fgWinResult.WinLine)
			if fgWinResult.WinType&EWT_FREEGAME != 0 {
				winLineCount = 0
			} else {
				winLineCount = len(fgWinResult.WinLine)
			}
			for winLineCount > 0 {
				cascadingRound++
				spinInfo2 := &SpinInfo{
					SpinState:      fgUseSpinState,
					Round:          j,
					FGTotalRound:   triggerFGInfo.FGTotalRound,
					CascadingRound: cascadingRound,
					ScreenInput:    screenInput,
					RNG:            fgLastRNG,
				}
				fgWinResult2 := singleWinResult(gameInput, gameOutput, spinInfo2)

				triggerFGInfo.FGTotalRound = fgWinResult2.FGTotalRound
				fgLastRNG = fgWinResult2.LastRNG
				screenInput = fgWinResult2.ScreenOutput

				if fgWinResult2.WinType&EWT_FREEGAME != 0 {
					winLineCount = 0
				} else {
					winLineCount = len(fgWinResult2.WinLine)
				}
			}
		}
	}
}
