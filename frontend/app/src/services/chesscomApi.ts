interface ChessComArchivesResponse {
  archives: string[];
}

interface ChessComGamesResponse {
  games: Array<{
    url: string;
    pgn: string;
    time_control: string;
    end_time: number;
    rated: boolean;
    white: { username: string; rating: number };
    black: { username: string; rating: number };
  }>;
}

class RateLimitError extends Error {
  retryAfter?: number;

  constructor(retryAfter?: number) {
    super('Rate limit exceeded');
    this.name = 'RateLimitError';
    this.retryAfter = retryAfter;
  }
}

const handleResponse = async (response: Response) => {
  if (response.status === 429) {
    const retryAfter = response.headers.get('Retry-After');
    throw new RateLimitError(retryAfter ? parseInt(retryAfter) : undefined);
  }

  if (!response.ok) {
    if (response.status === 404) {
      throw new Error('User not found');
    }
    throw new Error(`API error: ${response.statusText}`);
  }

  return response.json();
};

export const fetchUserArchives = async (
  username: string,
  signal?: AbortSignal
): Promise<ChessComArchivesResponse> => {
  const response = await fetch(
    `https://api.chess.com/pub/player/${username}/games/archives`,
    { signal }
  );

  return handleResponse(response);
};

export const fetchMonthGames = async (
  archiveUrl: string,
  signal?: AbortSignal
): Promise<ChessComGamesResponse> => {
  const response = await fetch(archiveUrl, { signal });
  return handleResponse(response);
};

export const fetchLatestGamePgn = async (
  username: string,
  signal?: AbortSignal
): Promise<string> => {
  // Validate username
  if (!username.trim()) {
    throw new Error('Please enter a username');
  }

  if (!/^[a-zA-Z0-9_-]{3,20}$/.test(username)) {
    throw new Error('Invalid username format');
  }

  try {
    const archivesResponse = await fetchUserArchives(username, signal);

    if (archivesResponse.archives.length === 0) {
      throw new Error(`User "${username}" has no game history`);
    }

    const latestArchiveUrl = archivesResponse.archives[archivesResponse.archives.length - 1];
    const monthGamesResponse = await fetchMonthGames(latestArchiveUrl, signal);

    if (monthGamesResponse.games.length === 0) {
      throw new Error('No games found in recent history');
    }

    const latestGame = monthGamesResponse.games[monthGamesResponse.games.length - 1];

    if (!latestGame.pgn || latestGame.pgn.trim() === '') {
      throw new Error('Game data is incomplete');
    }

    return latestGame.pgn;
  } catch (error) {
    if (error instanceof TypeError && error.message.includes('fetch')) {
      throw new Error('Network error. Check your connection.');
    }
    throw error;
  }
};

export { RateLimitError };
export type { ChessComArchivesResponse, ChessComGamesResponse };
