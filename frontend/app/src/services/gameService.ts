/**
 * WebSocket Game Service
 * Handles real-time communication with the chess backend for computer games
 */

import type {
  GameInfo,
  MoveResult,
  MoveEvent,
  GameResult,
  WebSocketMessage,
  Difficulty,
} from '../types/game';

type MessageHandler = (data: any) => void;
type EventHandler = (data: any) => void;

export class GameService {
  private ws: WebSocket | null = null;
  private userId: string = '';
  private requestId = 0;
  private pendingRequests = new Map<string, {
    resolve: (data: any) => void;
    reject: (error: Error) => void;
    timeout: NodeJS.Timeout;
  }>();
  private eventHandlers = new Map<string, EventHandler[]>();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000; // Start with 1 second
  private isIntentionalDisconnect = false;

  constructor(private wsUrl: string = 'ws://localhost:8080/ws') {}

  /**
   * Connect to WebSocket server
   */
  async connect(userId: string): Promise<void> {
    this.userId = userId;
    this.isIntentionalDisconnect = false;

    return new Promise((resolve, reject) => {
      try {
        const url = `${this.wsUrl}?user_id=${userId}`;
        this.ws = new WebSocket(url);

        this.ws.onopen = () => {
          console.log('[GameService] Connected to WebSocket');
          this.reconnectAttempts = 0;
          this.reconnectDelay = 1000;
          resolve();
        };

        this.ws.onmessage = (event) => {
          this.handleMessage(event.data);
        };

        this.ws.onerror = (error) => {
          console.error('[GameService] WebSocket error:', error);
          reject(new Error('WebSocket connection failed'));
        };

        this.ws.onclose = () => {
          console.log('[GameService] WebSocket closed');
          this.handleDisconnect();
        };
      } catch (error) {
        reject(error);
      }
    });
  }

  /**
   * Handle incoming WebSocket message
   */
  private handleMessage(data: string): void {
    try {
      const message: WebSocketMessage = JSON.parse(data);
      console.log('[GameService] Received:', message);

      if (message.type === 'response') {
        // Handle response to our request
        const pending = this.pendingRequests.get(message.id);
        if (pending) {
          clearTimeout(pending.timeout);
          this.pendingRequests.delete(message.id);

          if (message.success && !message.error) {
            pending.resolve(message.data);
          } else {
            pending.reject(new Error(message.error || 'Request failed'));
          }
        }
      } else if (message.type === 'event') {
        // Handle server event
        const handlers = this.eventHandlers.get(message.event || '');
        if (handlers) {
          handlers.forEach(handler => handler(message.data));
        }
      }
    } catch (error) {
      console.error('[GameService] Failed to parse message:', error);
    }
  }

  /**
   * Send request and wait for response
   */
  private sendRequest<T>(action: string, data: any, timeoutMs = 10000): Promise<T> {
    return new Promise((resolve, reject) => {
      if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
        reject(new Error('WebSocket not connected'));
        return;
      }

      const id = `req-${++this.requestId}`;
      const message: WebSocketMessage = {
        id,
        type: 'request',
        action,
        data,
      };

      // Set timeout for request
      const timeout = setTimeout(() => {
        this.pendingRequests.delete(id);
        reject(new Error(`Request timeout: ${action}`));
      }, timeoutMs);

      // Store pending request
      this.pendingRequests.set(id, {
        resolve,
        reject,
        timeout,
      });

      // Send message
      console.log('[GameService] Sending:', message);
      this.ws.send(JSON.stringify(message));
    });
  }

  /**
   * Register event handler
   */
  private on(event: string, handler: EventHandler): void {
    if (!this.eventHandlers.has(event)) {
      this.eventHandlers.set(event, []);
    }
    this.eventHandlers.get(event)!.push(handler);
  }

  /**
   * Unregister event handler
   */
  private off(event: string, handler: EventHandler): void {
    const handlers = this.eventHandlers.get(event);
    if (handlers) {
      const index = handlers.indexOf(handler);
      if (index > -1) {
        handlers.splice(index, 1);
      }
    }
  }

  /**
   * Handle WebSocket disconnect
   */
  private handleDisconnect(): void {
    // Clear all pending requests
    this.pendingRequests.forEach(({ reject, timeout }) => {
      clearTimeout(timeout);
      reject(new Error('WebSocket disconnected'));
    });
    this.pendingRequests.clear();

    // Attempt reconnection if not intentional
    if (!this.isIntentionalDisconnect && this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1); // Exponential backoff

      console.log(`[GameService] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);

      setTimeout(() => {
        this.connect(this.userId).catch(error => {
          console.error('[GameService] Reconnection failed:', error);
        });
      }, delay);
    }
  }

  /**
   * Create a new game
   */
  async createGame(
    mode: 'human_vs_human' | 'human_vs_computer',
    difficulty: Difficulty,
    timeControl: string,
    playerColor: 'white' | 'black'
  ): Promise<GameInfo> {
    return this.sendRequest<GameInfo>('createGame', {
      mode,
      difficulty,
      timeControl,
      playerColor,
    });
  }

  /**
   * Join an existing game
   */
  async joinGame(gameId: string): Promise<GameInfo> {
    return this.sendRequest<GameInfo>('joinGame', {
      gameId,
    });
  }

  /**
   * Make a move
   */
  async makeMove(gameId: string, move: string): Promise<MoveResult> {
    return this.sendRequest<MoveResult>('makeMove', {
      gameId,
      move,
    });
  }

  /**
   * Resign from game
   */
  async resign(gameId: string): Promise<void> {
    return this.sendRequest<void>('resign', {
      gameId,
    });
  }

  /**
   * Leave game
   */
  async leaveGame(gameId: string): Promise<void> {
    return this.sendRequest<void>('leaveGame', {
      gameId,
    });
  }

  /**
   * Listen for move events (computer moves, opponent moves)
   */
  onMove(callback: (move: MoveEvent) => void): () => void {
    this.on('move', callback);
    return () => this.off('move', callback);
  }

  /**
   * Listen for game end events
   */
  onGameEnd(callback: (result: GameResult) => void): () => void {
    this.on('gameEnd', callback);
    return () => this.off('gameEnd', callback);
  }

  /**
   * Listen for game start events (for multiplayer)
   */
  onGameStart(callback: (data: any) => void): () => void {
    this.on('gameStart', callback);
    return () => this.off('gameStart', callback);
  }

  /**
   * Listen for errors
   */
  onError(callback: (error: string) => void): () => void {
    this.on('error', callback);
    return () => this.off('error', callback);
  }

  /**
   * Disconnect from WebSocket
   */
  disconnect(): void {
    this.isIntentionalDisconnect = true;

    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }

    // Clear all handlers and pending requests
    this.eventHandlers.clear();
    this.pendingRequests.forEach(({ reject, timeout }) => {
      clearTimeout(timeout);
      reject(new Error('Service disconnected'));
    });
    this.pendingRequests.clear();
  }

  /**
   * Check if connected
   */
  isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
  }

  /**
   * Get connection state
   */
  getConnectionState(): 'disconnected' | 'connecting' | 'connected' | 'error' {
    if (!this.ws) return 'disconnected';

    switch (this.ws.readyState) {
      case WebSocket.CONNECTING:
        return 'connecting';
      case WebSocket.OPEN:
        return 'connected';
      case WebSocket.CLOSING:
      case WebSocket.CLOSED:
        return 'disconnected';
      default:
        return 'error';
    }
  }
}

// Export singleton instance
export const gameService = new GameService();

export default gameService;
