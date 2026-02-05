package game

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

func TestFish_GetResult(t *testing.T) {
	slot := Slot{}
	var result *Output
	var err error
	var IsCalcVI = false        //是否開啟計算 VI (開啟影響速度）
	var totalbet uint64         //總下注
	var totalwin uint64         //總贏分
	var totalspin uint64        //總下注次數
	var totalhitcount uint64    //總中獎次數
	var ngspin uint64           //NG下注次數
	var ngwin uint64            //NG贏分
	var nghitcount uint64       //NG中獎次數
	var fgspin uint64           //FG旋轉次數
	var fgwin uint64            //FG贏分
	var fghitcount uint64       //FG中獎次數
	var fgtimes uint64          //觸發FG的次數
	var fgretriggertimes uint64 //觸發FG的次數
	var wintotalwin uint64      //贏分時總贏分
	var wintimes uint64         //贏分次數

	//VI計算參數
	var meanRTP, delta, sumSigma float64 //1.meanRTP 表示的是逐步更新的 RTP 平均值。 2.delta 是每次迭代中樣本與當前均值的差異。 3.sumSigma 用來計算平方差的累積值。
	multipleCount := uint64(0)           //統計所有累積倍數的總和
	round := 10000 * 10000

	gameInput := &Input{
		RTPIndex:        0,     //RTP設定->[0]:95 [1]:97.5 [2]:99
		RoundMaxWinX:    0,     //風控：最大贏分倍率 <設定0 = 關閉風控>
		AgentMaxPayout:  0,     //風控：最大贏分 <設定0 = 關閉風控>
		BetAmount:       10000, //玩家押注單位
		BetLevel:        1,     //玩家押注倍數
		FeatureBuyIndex: "",    //買特色參數，直接傳入字串 ex."":自然機率
	}
	fmt.Println("RTPIndex:", gameInput.RTPIndex)
	for i := 1; i <= round; i++ {

		result, err = slot.GetResult(gameInput)
		if err != nil {
			fmt.Println(err)
			break
		}
		if IsCalcVI {
			sampleRTP := float64(result.TotalWin) / float64(result.TotalBet)
			delta = sampleRTP - meanRTP
			meanRTP += delta / float64(i+1)
			sumSigma += delta * (sampleRTP - meanRTP)
		}

		//統計
		if result.TotalWin > result.TotalBet {
			wintimes++
		}
		totalspin++
		totalbet += result.TotalBet
		totalwin += result.TotalWin
		if result.TotalWin > 0 {
			totalhitcount++
		}

		ngComboCount := 0
		for j := 0; j < len(result.WinResult); j++ {
			if GetSpinStateType(result.WinResult[j].SpinState) == ESST_NG {
				ngComboCount += 1
				if result.WinResult[j].CascadingRound == 1 {
					ngspin++
				}
				if result.WinResult[j].WinType != EWT_NULL {
					if result.WinResult[j].CascadingRound == 1 {
						nghitcount++
					}
					if result.WinResult[j].TotalScore >= 0 {
						ngwin += result.WinResult[j].TotalScore
					}
				}
				if result.WinResult[j].WinType&EWT_FREEGAME != 0 {
					fgtimes++
					if result.WinResult[j].TotalScore >= 0 {
						for k := 0; k < len(result.WinResult[j].WinLine); k++ {
							if result.WinResult[j].WinLine[k].SymbolType == EST_FG {
								ngwin -= result.WinResult[j].WinLine[k].Score
								fgwin += result.WinResult[j].WinLine[k].Score
							}
						}
					}
				}
			} else if GetSpinStateType(result.WinResult[j].SpinState) == ESST_FG {
				if result.WinResult[j].CascadingRound == 1 {
					fgspin++
				}
				if result.WinResult[j].WinType != EWT_NULL {
					if result.WinResult[j].CascadingRound == 1 {
						fghitcount++
					}
					if result.WinResult[j].TotalScore >= 0 {
						fgwin += result.WinResult[j].TotalScore
					}
				}
				if result.WinResult[j].WinType&EWT_FREEGAME != 0 {
					fgretriggertimes++
				}
			}
		}

		if i%(round/10) == 0 {
			fmt.Println("==== ", i/(round/10)*10, "% ====")
			fmt.Println("TotalRTP:", cutFloat((float64(totalwin)/float64(totalbet))*100), "%")
			fmt.Println("NGRTP:", cutFloat((float64(ngwin)/float64(totalbet))*100), "%")
			fmt.Println("FGRTP:", cutFloat((float64(fgwin)/float64(totalbet))*100), "%")
			fmt.Println("TotalHitRate:", cutFloat((float64(totalhitcount)/float64(totalspin))*100), "%")
			fmt.Println("NGHitRate:", cutFloat((float64(nghitcount)/float64(ngspin))*100), "%")
			fmt.Println("FGHitRate:", cutFloat((float64(fghitcount)/float64(fgspin))*100), "%")
			fmt.Println("FGTimesRate:", cutFloat((float64(fgtimes)/float64(totalspin))*100), "%")
			fmt.Println("FGReTimesRate:", cutFloat((float64(fgretriggertimes)/float64(fgspin))*100), "%")
			fmt.Println("WinRTP:", cutFloat((float64(wintotalwin)/float64(totalbet))*100), "%")
			fmt.Println("WinRate:", cutFloat((float64(wintimes)/float64(totalspin))*100), "%")
		}
	}

	if IsCalcVI {
		var VI float64
		// 樣本方差
		variance := sumSigma / float64(round-1)
		// VI 計算
		VI = math.Sqrt(variance)
		fmt.Println("倍數觸發次數:", multipleCount)
		fmt.Println(".")
		printSeparatorShort()
		fmt.Println("VI:", VI)
		printSeparatorShort()
		fmt.Println(".")
		fmt.Println(".")
	}
}

func cutFloat(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", value), 64)
	return value
}
