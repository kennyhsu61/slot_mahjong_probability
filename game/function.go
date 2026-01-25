package game

import (
	"fmt"
	"math/rand"
)

func getOption(g *GameOption) int {
	selectOption := 0
	if len(g.Options) == 0 || len(g.Weights) == 0 {
		return -1
	}
	if len(g.Options) != len(g.Weights) {
		return -1
	}

	// 計算權重總和
	totalWeight := 0
	for _, w := range g.Weights {
		if w < 0 {
			return -1
		}
		totalWeight += w
	}

	if totalWeight == 0 {
		return -1
	}

	// 生成隨機數（0 ~ totalWeight-1）
	r := rand.Intn(totalWeight)

	// 依照累積權重挑選
	sum := 0
	for i, w := range g.Weights {
		sum += w
		if r < sum {
			selectOption = g.Options[i]
			break
		}
	}

	return selectOption
}

func GetSpinStateType(spinstate SpinState) SpinStateType {
	switch {
	case spinstate == ESS_XX:
		return ESST_XX
	case (spinstate >= ESS_NG_1 && spinstate <= ESS_NG_20) || (spinstate >= ESS_BUY_NG_1 && spinstate <= ESS_BUY_NG_20):
		return ESST_NG
	case (spinstate >= ESS_FG_1 && spinstate <= ESS_FG_20) || (spinstate >= ESS_BUY_FG_1 && spinstate <= ESS_BUY_FG_20):
		return ESST_FG
	case (spinstate >= ESS_BG_1 && spinstate <= ESS_BG_20) || (spinstate >= ESS_BUY_BG_1 && spinstate <= ESS_BUY_BG_20):
		return ESST_BG
	default:
		return ESST_OTHERS
	}
}

func GetSymboltype(symbolID SymbolID) SymbolType {
	switch {
	case symbolID == E_XX:
		return EST_NULL
	case symbolID >= 1 && symbolID <= 40:
		return EST_NORMAL
	case symbolID >= 41 && symbolID <= 60:
		return EST_WD
	case symbolID >= 61 && symbolID <= 80:
		return EST_FG
	case symbolID >= 81 && symbolID <= 100:
		return EST_SC
	case symbolID >= 101 && symbolID <= 120:
		return EST_MX
	case symbolID >= 121 && symbolID <= 160:
		return EST_NORMAL
	case symbolID >= 161 && symbolID <= 180:
		return EST_WD
	case symbolID >= 181 && symbolID <= 200:
		return EST_FG
	case symbolID >= 201 && symbolID <= 220:
		return EST_SC
	case symbolID >= 221 && symbolID <= 240:
		return EST_MX
	default:
		return EST_OTHERS
	}
}

func IsMixSymbol(symbolID SymbolID) bool {

	var isMixSymbol bool
	var symbolno = int(symbolID)

	//用於顯示特殊型態的symbol
	if symbolno >= 121 && symbolno <= 240 {
		isMixSymbol = true
	}
	return isMixSymbol
}

func AddWinLine(winResult *WinResult, winLine *WinLine) {
	winResult.WinLine = append(winResult.WinLine, winLine)
	winResult.TotalScoreOrg += winLine.Score
	winResult.WinType = winResult.WinType | winLine.WinType
}

func checkInuptError(gameInput *Input) (*Output, error, bool) {
	//防呆 RTPIndex, FeatureBuyIndex, Selection
	if _, ok := Config.FeatureBuyInfo[gameInput.FeatureBuyIndex]; ok {
	} else {
		return nil, ErrFeatureBuyIndex, true
	}
	if gameInput.RTPIndex < 0 || gameInput.RTPIndex > 2 {
		return nil, ErrRTPIndex, true
	}
	if gameInput.Selection != 0 {
		return nil, ErrSelection, true
	}
	return nil, nil, false
}

// 分隔線1
func printSeparatorShort() {
	fmt.Println("-------------------------------------------")
}

// 分隔線2
func printSeparatorShortAny() {
	fmt.Println("------------------------------------")
}
