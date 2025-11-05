/**
 * Game-related types for computer player and WebSocket integration
 */

export type GameMode = 'live' | 'replay' | 'computer';

export type OpponentType = 'human' | 'computer';

export type Difficulty = 'easy' | 'medium' | 'hard';

export type PlayerColor = 'white' | 'black';

export type ConnectionStatus = 'disconnected' | 'connecting' | 'connected' | 'error';

export interface GameInfo {
  gameId: string;
  status: 'waiting' | 'active' | 'finished';
  side: PlayerColor;
  mode: string;
  difficulty: Difficulty;
  timeControl: string;
  fen: string;
}

export interface MoveResult {
  moveNumber: number;
  san: string;
  uci: string;
  fen: string;
  gameOver: boolean;
  result?: GameResult;
}

export interface MoveEvent {
  gameId: string;
  moveSAN: string;
  moveUCI: string;
  fen: string;
  moveNumber: number;
  isGameOver: boolean;
  result: GameResult | null;
}

export interface GameResult {
  winner: 'white' | 'black' | 'draw' | '';
  method: 'checkmate' | 'resignation' | 'stalemate' | 'timeout' | 'draw';
  pgn: string;
}

export interface WebSocketMessage {
  id: string;
  type: 'request' | 'response' | 'event';
  action?: string;
  event?: string;
  data?: any;
  error?: string;
  success?: boolean;
}
