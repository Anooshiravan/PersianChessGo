/*
   _____              _                _____ _
  |  __ \            (_)              / ____| |
  | |__) |__ _ __ ___ _  __ _ _ __   | |    | |__   ___  ___ ___
  |  ___/ _ \ '__/ __| |/ _` | '_ \  | |    | '_ \ / _ \/ __/ __|
  | |  |  __/ |  \__ \ | (_| | | | | | |____| | | |  __/\__ \__ \
  |_|   \___|_|  |___/_|\__,_|_| |_|  \_____|_| |_|\___||___/___/

════════════════════════════════════════════════════════════════════
 Persian Chess (www.PersianChess.com)
 Copyright 2006 - 2015
 Anooshiravan Ahmadi (aahmadi@schubergphilis.com)
 http://www.PersianChess.com/About
 Licensed under GNU General Public License 3.0
 ════════════════════════════════════════════════════════════════════
*/

package main

// ══════════════════════════
//  Engine Init
// ══════════════════════════

func StartEngine() {
	init_engine()
	SendMessageToGui("init", "engine_started")
}

func InitBoardVars() {

	// Clear brd_history
	brd_history_move = nil
	brd_history_castlePerm = nil
	brd_history_enPas = nil
	brd_history_fiftyMove = nil
	brd_history_posKey = nil

	// Clear PvTable
	brd_PvTable_move = nil
	brd_PvTable_posKey = nil
}

func EvalInit() {
	var index = 0

	for index = 0; index < 10; index++ {
		PawnRanksWhite[index] = 0
		PawnRanksBlack[index] = 0
	}
}

func InitHashKeys() {
	var index = 0

	for index = 0; index < 21*195; index++ {
		PieceKeys[index] = RAND_32()
	}

	SideKey = RAND_32()

	for index = 0; index < 16; index++ {
		CastleKeys[index] = RAND_32()
	}
}

func InitSq195To121() {

	var index = 0
	var file = FILE_A
	var rank = RANK_1
	var sq = SQUARES_A1
	var sq121 = 0
	for index = 0; index < BRD_SQ_NUM; index++ {
		Sq195ToSq121[index] = 122
	}

	for index = 0; index < 121; index++ {
		Sq121ToSq195[index] = 195
	}

	for rank = RANK_1; rank <= RANK_11; rank++ {
		for file = FILE_A; file <= FILE_K; file++ {
			sq = FR2SQ(file, rank)
			Sq121ToSq195[sq121] = sq
			Sq195ToSq121[sq] = sq121
			sq121++
		}
	}
}

func InitFilesRanksBrd() {

	var index = 0
	var file = FILE_A
	var rank = RANK_1
	var sq = SQUARES_A1

	for index = 0; index < BRD_SQ_NUM; index++ {
		FilesBrd[index] = SQUARES_OFFBOARD
		RanksBrd[index] = SQUARES_OFFBOARD
	}

	for rank = RANK_1; rank <= RANK_11; rank++ {
		for file = FILE_A; file <= FILE_K; file++ {
			sq = FR2SQ(file, rank)
			FilesBrd[sq] = file
			RanksBrd[sq] = rank

			// setting frame squares to offboard
			if contains(FrameSQ, sq) {
				FilesBrd[sq] = SQUARES_OFFBOARD
				RanksBrd[sq] = SQUARES_OFFBOARD
			}
		}
	}
}

func init_engine() {
	InitFilesRanksBrd()
	InitSq195To121()
	InitHashKeys()
	InitBoardVars()
	InitMvvLva()
	EvalInit()
	srch_thinking = false
}

func NewGame() {
	init_engine()
	ParseFen(START_FEN)
	if debug {
		PrintBoard()
	}
	GameController_PlayerSide = brd_side
	GameController_GameSaved = false
	SendMessageToGui("init", "new_game_started")
	// SendPosition();
}
