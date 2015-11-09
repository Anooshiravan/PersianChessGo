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

import (
	"strconv"
	"unicode/utf8"
)

// ══════════════════════════
//  Engine system board
// ══════════════════════════

var brd_side = COLOURS_WHITE
var brd_pieces [BRD_SQ_NUM]int
var brd_enPas = SQUARES_NO_SQ
var brd_fiftyMove = 0
var brd_ply = 0
var brd_hisPly = 0
var brd_castlePerm = 0
var brd_posKey = 0
var brd_pceNum [21]int
var brd_material [2]int
var brd_pList [22 * 11]int

var brd_history_move []int
var brd_history_castlePerm []int
var brd_history_enPas []int
var brd_history_fiftyMove []int
var brd_history_posKey []int

var brd_moveList [MAXDEPTH * MAXPOSITIONMOVES]int
var brd_moveScores [MAXDEPTH * MAXPOSITIONMOVES]int
var brd_moveListStart [MAXDEPTH]int

var brd_PvTable_move []int
var brd_PvTable_posKey []int
var brd_PvArray [MAXDEPTH]int
var brd_searchHistory [22 * BRD_SQ_NUM]int
var brd_searchKillers [3 * MAXDEPTH]int

/*

// board functions
function BoardToFen() {
    var fenStr = '';
    var rank, file, sq, piece;
    var emptyCount = 0;

    for (rank = RANKS.RANK_11; rank >= RANKS.RANK_1; rank--) {
        for (file = FILES.FILE_A; file <= FILES.FILE_K; file++) {
            sq = FR2SQ(file, rank);
            piece = brd_pieces[sq];
            if (piece == PIECES.EMPTY || piece == SQUARES.OFFBOARD) {
                fenStr += '1';
            } else {
                fenStr += PceChar[piece];
            }
        }

        if (rank != RANKS.RANK_1) {
            fenStr += '/';
        } else {
            fenStr += ' ';
        }
    }

    fenStr += SideChar[brd_side] + ' ';

    if (brd_castlePerm == 0) {
        fenStr += '- ';
    } else {
        if (brd_castlePerm & CASTLEBIT.WKCA) fenStr += 'K';
        if (brd_castlePerm & CASTLEBIT.WQCA) fenStr += 'Q';
        if (brd_castlePerm & CASTLEBIT.BKCA) fenStr += 'k';
        if (brd_castlePerm & CASTLEBIT.BQCA) fenStr += 'q';
        fenStr += ' ';
    }

    if (brd_enPas == SQUARES.NO_SQ) {
        fenStr += '- ';
    } else {
        fenStr += PrSq(brd_enPas) + ' ';
    }
    fenStr += brd_fiftyMove;

    if (brd_hisPly > 2) {
        fenStr += ' ';
        var tempHalfMove = brd_hisPly;
        if (brd_side == COLOURS.BLACK) {
            tempHalfMove--;
        }
        var pLy = tempHalfMove / 2;
        if (pLy < 1) pLy = 1;
        fenStr += pLy;
    }
    return fenStr;
}

function CheckBoard() {

    var t_pceNum = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0];
    var t_material = [0, 0];

    var sq121, t_piece, t_pce_num, sq195, colour, pcount;

    // check piece lists
    for (t_piece = PIECES.wP; t_piece <= PIECES.bK; ++t_piece) {
        for (t_pce_num = 0; t_pce_num < brd_pceNum[t_piece]; ++t_pce_num) {
            sq195 = brd_pList[PCEINDEX(t_piece, t_pce_num)];
            if (brd_pieces[sq195] != t_piece) {
                printLine('Error Pce Lists');
                return BOOL.FALSE;
            }
        }
    }

    // check piece count and other counters
    for (sq121 = 0; sq121 < 121; ++sq121) {
        sq195 = SQ195(sq121);
        t_piece = brd_pieces[sq195];
        t_pceNum[t_piece]++;
        t_material[PieceCol[t_piece]] += PieceVal[t_piece];
    }

    for (t_piece = PIECES.wP; t_piece <= PIECES.bK; ++t_piece) {
        if (t_pceNum[t_piece] != brd_pceNum[t_piece]) {
            printLine('Error t_pceNum');
            return BOOL.FALSE;
        }
    }

    if (t_material[COLOURS.WHITE] != brd_material[COLOURS.WHITE] || t_material[COLOURS.BLACK] != brd_material[COLOURS.BLACK]) {
        printLine('Error t_material');
        return BOOL.FALSE;
    }
    if (brd_side != COLOURS.WHITE && brd_side != COLOURS.BLACK) {
        printLine('Error brd_side');
        return BOOL.FALSE;
    }
    if (GeneratePosKey() != brd_posKey) {
        printLine('Error brd_posKey');
        return BOOL.FALSE;
    }


    return BOOL.TRUE;
}

function printGameLine() {

    var moveNum = 0;
    var gameLine = "";
    for (moveNum = 0; moveNum < brd_hisPly; ++moveNum) {
        gameLine += PrMove(brd_history[moveNum].move) + " ";
    }
    //printLine('Game Line: ' + gameLine);
    return gameLine.trim();
}

function PrintPceLists() {
    var piece, pceNum;

    for (piece = PIECES.wP; piece <= PIECES.bK; ++piece) {
        for (pceNum = 0; pceNum < brd_pceNum[piece]; ++pceNum) {
            printLine("Piece " + PceChar[piece] + " on " + PrSq(brd_pList[PCEINDEX(piece, pceNum)]));
        }
    }

}

function UpdateListsMaterial() {

    var piece, sq, index, colour;

    for (index = 0; index < BRD_SQ_NUM; ++index) {
        sq = index;
        piece = brd_pieces[index];
        if (piece != PIECES.OFFBOARD && piece != PIECES.EMPTY) {
            colour = PieceCol[piece];

            brd_material[colour] += PieceVal[piece];

            brd_pList[PCEINDEX(piece, brd_pceNum[piece])] = sq;
            brd_pceNum[piece]++;
        }
    }
}
*/

func GeneratePosKey() int {

	var sq = 0
	var finalKey = 0
	var piece = PIECES_EMPTY

	// pieces
	for sq = 0; sq < BRD_SQ_NUM; sq++ {
		piece = brd_pieces[sq]
		if piece != PIECES_EMPTY && piece != SQUARES_OFFBOARD {
			finalKey ^= PieceKeys[(piece*195)+sq]
		}
	}

	if brd_side == COLOURS_WHITE {
		finalKey ^= SideKey
	}

	if brd_enPas != SQUARES_NO_SQ {
		finalKey ^= PieceKeys[brd_enPas]
	}

	finalKey ^= CastleKeys[brd_castlePerm]

	return finalKey
}

func PrintBoard() {

	var sq, file, rank, piece int
	var line string

	printLine("\nGame Board:\n")

	for rank = RANK_11; rank >= RANK_1; rank-- {
		line = ""
		if rank+1 > 9 {
			line = (strconv.Itoa(rank+1) + "|")
		} else {
			line = (strconv.Itoa(rank+1) + " |")
		}

		for file = FILE_A; file <= FILE_K; file++ {
			sq = FR2SQ(file, rank)
			piece = brd_pieces[sq]
			if piece == SQUARES_OFFBOARD {
				line += " * "
			} else {
				line += " " + PceChar[piece] + " "
			}
		}
		printLine(line)
	}

	printLine("")
	line = "   "
	for file = FILE_A; file <= FILE_K; file++ {
		line += " " + FileChar[file] + " "
	}
	printLine(line)
	printLine("")
	printLine("side:" + SideChar[brd_side])
	printLine("enPas:" + strconv.Itoa(brd_enPas))
	line = ""

	if brd_castlePerm&CASTLEBIT_WKCA != 0 {
		line += "K"
	}
	if brd_castlePerm&CASTLEBIT_WQCA != 0 {
		line += "Q"
	}
	if brd_castlePerm&CASTLEBIT_BKCA != 0 {
		line += "k"
	}
	if brd_castlePerm&CASTLEBIT_BQCA != 0 {
		line += "q"
	}

	printLine("castle:" + line)
	printLine("key:" + strconv.Itoa(brd_posKey))
}

func ResetBoard() {

	var index = 0

	for index = 0; index < BRD_SQ_NUM; index++ {
		brd_pieces[index] = SQUARES_OFFBOARD
	}

	for index = 0; index < 121; index++ {
		if contains(FrameSQ, SQ195(index)) {
			brd_pieces[SQ195(index)] = SQUARES_OFFBOARD
		} else {
			brd_pieces[SQ195(index)] = PIECES_EMPTY
		}
	}

	for index = 0; index < 22*11; index++ {
		brd_pList[index] = PIECES_EMPTY
	}

	for index = 0; index < 2; index++ {
		brd_material[index] = 0
	}

	for index = 0; index < 21; index++ {
		brd_pceNum[index] = 0
	}

	brd_side = COLOURS_BOTH
	brd_enPas = SQUARES_NO_SQ
	brd_fiftyMove = 0
	brd_ply = 0
	brd_hisPly = 0
	brd_castlePerm = 0
	brd_posKey = 0
	brd_moveListStart[brd_ply] = 0

}

func ParseFen(fen string) bool {

	var rank = RANK_11
	var file = FILE_A
	var piece = 0
	var count = 0
	var i = 0
	var sq121 = 0
	var sq195 = 0
	var fenCnt = 0

	ResetBoard()

	for rank >= RANK_1 && fenCnt < utf8.RuneCountInString(fen) {
		count = 1
		switch fen[fenCnt] {
		case 'p':
			piece = bP
			break
		case 'r':
			piece = bR
			break
		case 'n':
			piece = bN
			break
		case 'w':
			piece = bW
			break
		case 'c':
			piece = bC
			break
		case 'b':
			piece = bB
			break
		case 's':
			piece = bS
			break
		case 'f':
			piece = bF
			break
		case 'k':
			piece = bK
			break
		case 'q':
			piece = bQ
			break
		case 'P':
			piece = wP
			break
		case 'R':
			piece = wR
			break
		case 'N':
			piece = wN
			break
		case 'W':
			piece = wW
			break
		case 'C':
			piece = wC
			break
		case 'B':
			piece = wB
			break
		case 'S':
			piece = wS
			break
		case 'F':
			piece = wF
			break
		case 'K':
			piece = wK
			break
		case 'Q':
			piece = wQ
			break
		case '1':
			piece = PIECES_EMPTY
			break

		case '/':
		case ' ':
			rank--
			file = FILE_A
			fenCnt++
			continue

		default:
			printLine("FEN error \n")
			return false
		}

		for i = 0; i < count; i++ {
			sq121 = rank*11 + file
			sq195 = SQ195(sq121)
			if piece != PIECES_EMPTY {
				if brd_pieces[sq195] != SQUARES_OFFBOARD {
					brd_pieces[sq195] = piece
				}

			}
			file++
		}
		fenCnt++
	}

	/*

	   brd_side = (fen[fenCnt] == 'w') ? COLOURS.WHITE : COLOURS.BLACK;
	   fenCnt += 2;

	   for (i = 0; i < 4; i++) {
	       if (fen[fenCnt] == ' ') {
	           break;
	       }
	       switch (fen[fenCnt]) {

	           case 'K':
	               brd_castlePerm |= CASTLEBIT.WKCA;
	               break;
	           case 'Q':
	               brd_castlePerm |= CASTLEBIT.WQCA;
	               break;
	           case 'k':
	               brd_castlePerm |= CASTLEBIT.BKCA;
	               break;
	           case 'q':
	               brd_castlePerm |= CASTLEBIT.BQCA;
	               break;
	           default:
	               break;
	       }
	       fenCnt++;
	   }
	   fenCnt++;

	   if (fen[fenCnt] != '-' && fen[fenCnt] != undefined) {
	       file = fen[fenCnt].charCodeAt() - 'a'.charCodeAt();
	       rank = fen[fenCnt + 1].charCodeAt() - '1'.charCodeAt();
	       printLine("fen[fenCnt]:" + fen[fenCnt] + " File:" + file + " Rank:" + rank);
	       brd_enPas = FR2SQ(file, rank);
	   }

	   brd_posKey = GeneratePosKey();
	   UpdateListsMaterial();
	*/
	return true

}

/*

function SqAttacked(sq, side) {
    var pce;
    var t_sq;
    var index;

    if (brd_pieces[sq] == SQUARES.OFFBOARD) return BOOL.FALSE;

    if (variant == "ASE" && ASEDIA.indexOf(sq) > -1) return BOOL.TRUE;

    if (side == COLOURS.WHITE) {
        if (brd_pieces[sq - 14] == PIECES.wP || brd_pieces[sq - 12] == PIECES.wP) {
            return BOOL.TRUE;
        }
    } else {
        if (brd_pieces[sq + 14] == PIECES.bP || brd_pieces[sq + 12] == PIECES.bP) {
            return BOOL.TRUE;
        }
    }

    // Knight, Princess and Fortress (non slide moves)

    for (index = 0; index < 8; ++index) {
        pce = brd_pieces[sq + KnDir[index]];
        if (pce != SQUARES.OFFBOARD && PieceKnightPrincessFortress[pce] == BOOL.TRUE && PieceCol[pce] == side) {
            return BOOL.TRUE;
        }
    }

    // Rook, Fortress and Queen (slide moves)

    for (index = 0; index < 4; ++index) {
        dir = RkDir[index];
        t_sq = sq + dir;
        pce = brd_pieces[t_sq];
        while (pce != SQUARES.OFFBOARD) {
            if (pce != PIECES.EMPTY) {
                if (PieceRookFortressQueen[pce] == BOOL.TRUE && PieceCol[pce] == side) {
                    return BOOL.TRUE;
                }
                break;
            }
            t_sq += dir;
            pce = brd_pieces[t_sq];
        }
    }

    // Bishop, Princess and Queen (slide moves)

    for (index = 0; index < 4; ++index) {
        dir = BiDir[index];
        t_sq = sq + dir;
        pce = brd_pieces[t_sq];
        while (pce != SQUARES.OFFBOARD) {
            if (pce != PIECES.EMPTY) {
                if (PieceBishopPrincessQueen[pce] == BOOL.TRUE && PieceCol[pce] == side) {
                    return BOOL.TRUE;
                }
                break;
            }
            t_sq += dir;
            pce = brd_pieces[t_sq];
        }
    }

    // Wizard and Champion

    if (variant == "Oriental") {
        for (index = 0; index < 12; ++index) {
            pce = brd_pieces[sq + WzDir[index]];
            if (pce != SQUARES.OFFBOARD && PieceWizard[pce] == BOOL.TRUE && PieceCol[pce] == side) {
                return BOOL.TRUE;
            }
        }

        for (index = 0; index < 12; ++index) {
            pce = brd_pieces[sq + ChDir[index]];
            if (pce != SQUARES.OFFBOARD && PieceChampion[pce] == BOOL.TRUE && PieceCol[pce] == side) {
                return BOOL.TRUE;
            }
        }
    }

    // King

    for (index = 0; index < 8; ++index) {
        pce = brd_pieces[sq + KiDir[index]];
        if (pce != SQUARES.OFFBOARD && PieceKing[pce] == BOOL.TRUE && PieceCol[pce] == side) {
            return BOOL.TRUE;
        }
    }

    return BOOL.FALSE;
}

function PrintSqAttacked() {

    var sq, file, rank, piece, line;

    printLine("\nAttacked by Black:\n");

    for (rank = RANKS.RANK_11; rank >= RANKS.RANK_1; rank--) {
        line = ((rank + 1) + "  ");
        for (file = FILES.FILE_A; file <= FILES.FILE_K; file++) {
            sq = FR2SQ(file, rank);
            if (SqAttacked(sq, COLOURS.BLACK) == BOOL.TRUE) piece = "X";
            else if (brd_pieces[sq] == SQUARES.OFFBOARD) piece = "*";
            else piece = "-";
            line += (" " + piece + " ");
        }
        printLine(line);
    }

    printLine("\nAttacked by White:\n");

    for (rank = RANKS.RANK_11; rank >= RANKS.RANK_1; rank--) {
        line = ((rank + 1) + "  ");
        for (file = FILES.FILE_A; file <= FILES.FILE_K; file++) {
            sq = FR2SQ(file, rank);
            if (SqAttacked(sq, COLOURS.WHITE) == BOOL.TRUE) piece = "X";
            else if (brd_pieces[sq] == SQUARES.OFFBOARD) piece = "*";
            else piece = "-";
            line += (" " + piece + " ");
        }
        printLine(line);
    }
}

function EvaluateSqAttacked() {
    // This function is not used in the evaluation yet, it is very slow
    var SqAttackedByWhite = 0;
    var SqAttackedByBlack = 0;
    var SqAttackedScore = 0;
    for (rank = RANKS.RANK_11; rank >= RANKS.RANK_1; rank--) {
        for (file = FILES.FILE_A; file <= FILES.FILE_K; file++) {
            sq = FR2SQ(file, rank);
            if (SqAttacked(sq, COLOURS.WHITE) == BOOL.TRUE) SqAttackedByWhite++;
            else if (SqAttacked(sq, COLOURS.BLACK) == BOOL.TRUE) SqAttackedByBlack++;
        }
    }
    SqAttackedScore = SqAttackedByWhite - SqAttackedByBlack;
    return SqAttackedScore;
}

*/
