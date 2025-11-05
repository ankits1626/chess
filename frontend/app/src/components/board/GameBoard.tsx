import { useGameStore } from '@/store/useGameStore';
import Square from './Square';
import type {
  SquareColor,
  Square as SquareType,
  ChessFile,
  ChessRank,
  ChessPiece,
  PieceType,
  PieceColor,
} from '@/types/chess';

const GameBoard = () => {
  const {
    game,
    selectedSquare,
    validMoves,
    lastMove,
    selectSquare: onSquareClick,
    mode,
    playerColor,
    isBoardFlipped,
  } = useGameStore();

  const board = game.board();

  const files: ChessFile[] = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'];
  const ranks: ChessRank[] = ['1', '2', '3', '4', '5', '6', '7', '8'];

  // Flip board based on player color or manual flip toggle
  // In replay/live mode: just use manual flip state
  // In computer mode: auto-flip for black, manual flip toggles (XOR logic)
  const autoFlip = playerColor === 'black';
  const isFlipped = playerColor ? (isBoardFlipped ? !autoFlip : autoFlip) : isBoardFlipped;

  const displayFiles = isFlipped ? [...files].reverse() : files;
  const displayRanks = isFlipped ? [...ranks] : [...ranks].reverse();

  const squares: Array<{
    name: SquareType;
    color: SquareColor;
    piece: ChessPiece | null;
  }> = [];

  // Iterate through ranks and files according to board orientation
  for (let rankIdx = 0; rankIdx < 8; rankIdx++) {
    const rank = displayRanks[rankIdx];
    const rankIndex = ranks.indexOf(rank);

    // Iterate through files
    for (let fileIdx = 0; fileIdx < 8; fileIdx++) {
      const file = displayFiles[fileIdx];
      const fileIndex = files.indexOf(file);
      const squareName = `${file}${rank}` as SquareType;

      // Calculate color
      const isLight = (rankIndex + fileIndex) % 2 !== 0;
      const squareColor: SquareColor = isLight ? 'light' : 'dark';

      // Get piece from board
      const chessJsPiece = board[7 - rankIndex][fileIndex];
      const pieceData = chessJsPiece
        ? {
            type: chessJsPiece.type as PieceType,
            color: chessJsPiece.color as PieceColor,
          }
        : null;

      squares.push({ name: squareName, color: squareColor, piece: pieceData });
    }
  }

  // Calculate board size based on mode
  // In replay mode: subtract space for player names (2 * 3rem), gaps (2 * 0.5rem), padding (4rem), file labels (1.5rem)
  // Total: ~14rem to subtract in replay mode
  const boardSize = mode === 'replay'
    ? 'w-[min(90vw,calc(100vh-16rem))] h-[min(90vw,calc(100vh-16rem))]'
    : 'w-[min(90vw,calc(100vh-8rem))] h-[min(90vw,calc(100vh-8rem))]';

  return (
    <div className={`relative ${boardSize} mx-auto transition-all duration-300`}>
      {/* Rank labels on the left */}
      <div className="absolute -left-6 top-0 h-full flex flex-col justify-around text-gray-400 text-sm transition-all duration-300">
        {displayRanks.map((rank) => (
          <span key={`rank-${rank}`} className="flex items-center justify-center h-full">
            {rank}
          </span>
        ))}
      </div>

      <div className="w-full h-full grid grid-cols-8 grid-rows-8 border-2 border-[#759656] transition-all duration-300">
        {squares.map((square) => (
          <Square
            key={square.name}
            squareColor={square.color}
            piece={square.piece}
            isSelected={selectedSquare === square.name}
            isValidMove={validMoves.includes(square.name)}
            isLastMoveFrom={lastMove?.from === square.name}
            isLastMoveTo={lastMove?.to === square.name}
            onClick={() => onSquareClick(square.name)}
          />
        ))}
      </div>

      {/* File labels at the bottom */}
      <div className="absolute -bottom-5 left-0 w-full flex justify-around text-gray-400 text-sm px-2 transition-all duration-300">
        {displayFiles.map((file) => (
          <span key={`file-${file}`} className="flex items-center justify-center w-full">
            {file}
          </span>
        ))}
      </div>
    </div>
  );
};

export default GameBoard;
