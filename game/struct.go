package game

type Input struct {
	RTPIndex        int         //RTP設定->[0]:95 [1]:97.5 [2]:99
	BetAmount       uint64      //玩家押注單位
	BetLevel        uint64      //玩家押注倍數
	Selection       int         //玩家選擇項（預留欄位）
	Collection      []int       //收集品（預留欄位）
	FeatureBuyIndex string      //買特色參數，直接傳入字串 ex."":自然機率、"FG"
	RoundMaxWinX    int64       //風控：最大贏分倍率->會影響RTP，請斟酌使用
	AgentMaxPayout  uint64      //風控：最大贏分->會影響RTP，請斟酌使用
	ExtraData       interface{} //特殊需要存取的資訊
}

type Output struct {
	TotalBet   uint64       //總下注
	TotalWin   uint64       //總贏分
	Collection []int        //收集品（預留欄位）
	RTPIndex   int          //RTP設定->[0]:95 [1]:97.5 [2]:99
	WinResult  []*WinResult //贏分資訊
	ExtraData  interface{}  //特殊需要存取的資訊
}

type WinResult struct {
	SpinState        SpinState    //滾輪狀態
	TotalScore       uint64       //單次 Spin 總贏分＝原始單次 Spin 總贏分 x 盤面倍率
	TotalScoreOrg    uint64       //原始單次 Spin 總贏分
	WinType          WinType      //贏分類型（一般贏分、觸發FG、觸發Bonus）
	WinLine          []*WinLine   //贏線資訊
	Screen           [][]SymbolID //盤面
	ScreenOrg        [][]SymbolID //原始盤面
	MoneyScreen      [][]int      //分數盤面
	MoneyScreenOrg   [][]int      //原始分數盤面
	ExtraScreen      [][]SymbolID //特殊盤面
	ExtraScreenOrg   [][]SymbolID //原始特殊盤面
	ScreenClear      [][]SymbolID //消除盤面
	ScreenOutput     [][]SymbolID //輸出盤面
	RNG              []int        //RNG 位置
	LastRNG          []int        //最後 RNG 位置（消除類）
	FGTotalRound     int          //Free Game 總場次
	FGRound          int          //Free Game 當前場次
	BGSurviveRound   int          //Bonus 剩餘場次
	CascadingRound   int          //消除層數
	ScreenMultiplier uint64       //盤面倍率
}

type WinLine struct {
	LineNo         int        //連線編號
	WinType        WinType    //贏分類型（一般贏分、觸發FG、觸發Bonus）
	SymbolID       SymbolID   //圖標
	SymbolType     SymbolType //圖標類型
	Count          int        //連線長度
	Ways           uint64     //堆疊 Way 數
	Score          uint64     //分數 = 原始分數 x 單線倍數
	ScoreOrg       uint64     //原始分數
	LineMultiplier uint64     //單線倍數
	Position       [][]int    //連線位置
}

type GameConfig struct {
	//遊戲設定
	GameLine        uint64                              //遊戲基注or線數
	ScreenSize      []int                               //行列數 EX: [3,3,3,3,3]
	GameRule        GameRule                            //遊戲規則
	LineTable       map[int][][]int                     //賠付線
	OddsTable       map[SymbolID][]uint64               //賠付表
	StripTable      map[SpinState][][]SymbolID          //滾輪表
	ShowStripInfo   map[SpinState]ShowReel              //各狀態對應哪個表演滾輪
	ShowStripTable  map[ShowReel][][]SymbolID           //表演用滾輪表
	GameStripOption map[int]map[SpinState][]*GameOption //[RTPIndex]遊戲權重
	GameOption      map[SpinState][]*GameOption         //遊戲權重
	FeatureBuyInfo  map[string]FeatureBuyInfo           //買特色資訊
}

type GameRule struct {
	RoundFG          []int
	RoundFGRetrigger []int
	RoundFGMax       int
	RoundBGSurvive   []int
	ScreenMultiArray map[SpinStateType][][]uint64
	WDMultiArray     map[SpinState][][]uint64
	SCMultiArray     map[SpinState][][]uint64
}

type GameOption struct {
	Options []int //選項
	Weights []int //權重
}

type FeatureBuyInfo struct {
	FeatureBuyLevel uint64
}

type TriggerFGInfo struct {
	WinType      WinType
	FGTotalRound int
}
