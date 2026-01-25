package game

func singleWinResult(gameInput *Input, gameOutput *Output, spinInfo *SpinInfo) *WinResult {

	calculator := getCalcData(gameInput, spinInfo)
	winResult := &WinResult{}
	resetScreenClear(winResult, calculator)
	calculator.calcWayWin(winResult)
	calculator.calcWinRecord(winResult)
	calculator.checkScreenOutput(winResult)
	if len(winResult.WinLine) == 0 {
		calculator.calcFGWin(winResult)
	}
	calculator.calcFGTriggerInfo(winResult, spinInfo)
	//填入其他資訊
	winResult.Screen = spinInfo.Screen
	winResult.SpinState = calculator.SpinState
	winResult.CascadingRound = spinInfo.CascadingRound
	winResult.ScreenMultiplier = calculator.ScreenMulti
	winResult.TotalScore = winResult.TotalScoreOrg * calculator.ScreenMulti
	winResult.RNG = calculator.RNG

	gameOutput.WinResult = append(gameOutput.WinResult, winResult)
	gameOutput.TotalWin += winResult.TotalScore

	return winResult
}

func resetScreenClear(winResult *WinResult, calculator *CalcData) {
	//重置ScreenClear
	winResult.ScreenClear = make([][]SymbolID, len(calculator.Screen))
	for i := 0; i < len(calculator.Screen); i++ {
		winResult.ScreenClear[i] = make([]SymbolID, len(calculator.Screen[i]))
		copy(winResult.ScreenClear[i], calculator.Screen[i])
	}
}

func (calculator *CalcData) calcWayWin(winResult *WinResult) {

	var screen = calculator.Screen
	var oddstable = Config.OddsTable

	//此段使用的變數
	var winlinelist = map[SymbolID]*WinLine{}

	for i := 0; i < len(screen[0]); i++ {
		var symbol = screen[0][i]
		if GetSymboltype(symbol) == EST_NORMAL {
			value, exist := winlinelist[symbol]
			if exist {
				value.Ways++
				value.Position = append(value.Position, []int{0, i})
			} else {
				var winunit = &WinLine{
					LineNo:         len(winlinelist) + 1,
					SymbolID:       screen[0][i],
					SymbolType:     GetSymboltype(screen[0][i]),
					Count:          1,
					Ways:           1,
					LineMultiplier: 1,
					Position:       [][]int{{0, i}},
					WinType:        EWT_NORMAL,
				}
				winlinelist[symbol] = winunit
			}
		}
	}

	for i := 1; i < len(screen); i++ {

		var calcWayList = map[SymbolID]uint64{}

		for j := 0; j < len(screen[i]); j++ {

			var symbol = screen[i][j]
			//特殊符號減120與原符號一起計算連線
			if IsMixSymbol(symbol) {
				symbol -= 120
			}

			if GetSymboltype(symbol) == EST_NORMAL {
				value, exist := winlinelist[symbol]
				if exist {
					if value.Count >= i {
						value.Count = i + 1
						calcWayList[symbol]++
						value.Position = append(value.Position, []int{i, j})
					}
				}
			} else if GetSymboltype(symbol) == EST_WD {
				for key, value := range winlinelist {
					if value.Count >= i {
						value.Count = i + 1
						calcWayList[key]++
						value.Position = append(value.Position, []int{i, j})
					}
				}
			}
		}
		for key, value := range winlinelist {
			if calcWayList[key] != 0 {
				value.Ways = value.Ways * calcWayList[key]
			}
		}
	}

	for key, value := range winlinelist {
		value.ScoreOrg = oddstable[key][value.Count-1] * calculator.Bet * value.Ways
		value.Score = value.ScoreOrg * value.LineMultiplier
		if value.Score > 0 {
			AddWinLine(winResult, value)
		}
	}
}

func (calculator *CalcData) calcWinRecord(winResult *WinResult) {

	//取得需要的資訊
	var screen [][]SymbolID
	screen = calculator.Screen //取得盤面

	//產生紀錄用的flagMap  (該位置是否已經被消除?)
	flagMap := make([][]bool, len(screen))
	for k := 0; k < len(screen); k++ {
		flagMap[k] = make([]bool, len(screen[k]))
	}

	//針對贏分線去做判斷
	for _, value := range winResult.WinLine {
		for _, pos := range value.Position {
			//該贏分線Symbol為WD類型不消除
			if flagMap[pos[0]][pos[1]] == false {

				//更新已消除盤面，M系列符號轉換成WD不消除
				if IsMixSymbol(winResult.ScreenClear[pos[0]][pos[1]]) {
					winResult.ScreenClear[pos[0]][pos[1]] = E_WD
				} else {
					winResult.ScreenClear[pos[0]][pos[1]] = E_XX
				}

				////將該位置先標記為已展演
				//flagMap[pos[0]][pos[1]] = true
				//
				//allSymbol, exist1 := calculator.WinRecord[E_XX]
				//if !exist1 {
				//	//全部Symbol未實作要new
				//	allSymbol = new(WinRecord)
				//	calculator.WinRecord[E_XX] = allSymbol
				//}
				//
				//selectSymbol, exist2 := calculator.WinRecord[screen[pos[0]][pos[1]]]
				//if !exist2 {
				//	//個別Symbol未實作要new
				//	selectSymbol = new(WinRecord)
				//	calculator.WinRecord[screen[pos[0]][pos[1]]] = selectSymbol
				//}
				////紀錄消除數量
				//allSymbol.TotalCount++
				//selectSymbol.TotalCount++
			}
		}
	}

	//winResult.WinRecord = calculator.WinRecord
}

func (calculator *CalcData) checkScreenOutput(winResult *WinResult) {

	//取得需要的資訊
	var screen = calculator.Screen          //取得盤面
	var screenClear = winResult.ScreenClear //取得消除盤面

	//產生暫存掉落盤面
	tempScreen := make([][]SymbolID, len(screenClear))
	for k := 0; k < len(screen); k++ {
		tempScreen[k] = make([]SymbolID, len(screenClear[k]))
	}

	//產生新盤面流程＆步驟： 剩下盤面跌落
	for i := 0; i < len(screenClear); i++ {
		tempY := len(screenClear[i]) - 1
		for j := len(screenClear[i]) - 1; j >= 0; j-- {
			if screenClear[i][j] != E_XX {
				tempScreen[i][tempY] = screenClear[i][j]
				tempY--
			}
		}
	}

	var lastRndNum []int
	lastRndNum = make([]int, len(Config.ScreenSize))
	for i := 0; i < len(calculator.RNG); i++ {
		lastRndNum[i] = calculator.RNG[i]
	}
	//產生新盤面流程＆步驟： 新盤面補充(目前產盤面為簡易版)
	for i := 0; i < len(tempScreen); i++ {
		reel := calculator.StripTable[i]
		stripLen := len(reel)
		for j := len(screenClear[i]) - 1; j >= 0; j-- {
			if tempScreen[i][j] == E_XX {
				if lastRndNum[i] == 0 {
					lastRndNum[i] = stripLen - 1
				} else {
					lastRndNum[i]--
				}
				pos := lastRndNum[i]
				tempScreen[i][j] = reel[pos]
			}
		}
	}

	winResult.LastRNG = lastRndNum
	winResult.ScreenOutput = tempScreen
}

func (calculator *CalcData) calcFGWin(winResult *WinResult) {

	//取得需要的資訊
	var screen = calculator.Screen     //取得盤面
	var hitcount = 3                   //幾顆FG進FreeGame
	var oddstable = Config.OddsTable   //使用的賠率表
	var totalbet = calculator.TotalBet //總押注
	var unitmulti uint64 = 1           //單元倍率
	var count int
	var pos [][]int

	for i := 0; i < len(screen); i++ {
		for j := 0; j < len(screen[i]); j++ {
			if GetSymboltype(screen[i][j]) == EST_FG {
				count++
				pos = append(pos, []int{i, j})
			}
		}
	}

	//判斷有無FG
	var wintype WinType
	if count >= hitcount {
		wintype = EWT_FREEGAME
	} else {
		wintype = EWT_NULL
	}

	//準備WinUnit, 寫入WinInfo
	if wintype != EWT_NULL {
		odds := oddstable[E_FG][0]
		scoreorg := odds * totalbet
		score := scoreorg * unitmulti

		winUnit := &WinLine{
			LineNo:         0,
			SymbolID:       E_FG,
			SymbolType:     EST_FG,
			Count:          count,
			Score:          score,
			ScoreOrg:       scoreorg,
			LineMultiplier: unitmulti,
			Position:       pos,
			WinType:        wintype,
		}

		AddWinLine(winResult, winUnit)
	}

}

func (calculator *CalcData) calcFGTriggerInfo(winResult *WinResult, spinInfo *SpinInfo) {
	if winResult.WinType&EWT_FREEGAME != 0 {
		if GetSpinStateType(calculator.SpinState) == ESST_NG {
			fgHitNum := winResult.WinLine[0].Count
			extraHits := fgHitNum - 3
			if extraHits < 0 {
				extraHits = 0
			}
			winResult.FGTotalRound = Config.GameRule.RoundFG[0] + extraHits*Config.GameRule.RoundFG[1]
			winResult.FGRound = 0
		} else if GetSpinStateType(calculator.SpinState) == ESST_FG {

			rertiggerHitNum := winResult.WinLine[0].Count
			extraHits := rertiggerHitNum - 3
			if extraHits < 0 {
				extraHits = 0
			}
			winResult.FGTotalRound = spinInfo.FGTotalRound + Config.GameRule.RoundFGRetrigger[0] + extraHits*Config.GameRule.RoundFGRetrigger[1]
			if winResult.FGTotalRound > Config.GameRule.RoundFGMax {
				winResult.FGTotalRound = Config.GameRule.RoundFGMax
			}
			winResult.FGRound = spinInfo.Round
		}
	} else {
		if GetSpinStateType(calculator.SpinState) == ESST_NG {
		} else if GetSpinStateType(calculator.SpinState) == ESST_FG {
			winResult.FGTotalRound = spinInfo.FGTotalRound
			winResult.FGRound = spinInfo.Round
		}
	}
}
