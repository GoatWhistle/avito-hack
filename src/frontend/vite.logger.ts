import { createLogger, type Logger, type Plugin } from 'vite';

const RESET = '\x1b[0m';
const BOLD = '\x1b[1m';
const DIM = '\x1b[2m';
const RED = '\x1b[31m';
const GREEN = '\x1b[32m';
const YELLOW = '\x1b[33m';
const BLUE = '\x1b[34m';
const MAGENTA = '\x1b[35m';
const CYAN = '\x1b[36m';
const WHITE = '\x1b[37m';

const METHOD_COLORS: Record<string, string> = {
  GET: BLUE,
  POST: GREEN,
  PUT: YELLOW,
  PATCH: YELLOW,
  DELETE: RED,
};

const statusColor = (status: number): string => {
  if (status >= 500) return RED;
  if (status >= 400) return YELLOW;
  if (status >= 300) return CYAN;
  return GREEN;
};

const durationColor = (ms: number): string => {
  if (ms >= 500) return RED;
  if (ms >= 100) return YELLOW;
  return DIM;
};

const formatDuration = (ms: number): string => {
  if (ms >= 1000) return `${(ms / 1000).toFixed(2)}s`;
  if (ms >= 1) return `${ms.toFixed(1)}ms`;
  return `${Math.round(ms * 1000)}µs`;
};

const formatBytes = (bytes: number): string => {
  if (bytes >= 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)}MB`;
  if (bytes >= 1024) return `${(bytes / 1024).toFixed(1)}KB`;
  return `${bytes}B`;
};

const timestamp = (): string => {
  const now = new Date();
  const pad = (value: number, size = 2): string => String(value).padStart(size, '0');
  return `${pad(now.getHours())}:${pad(now.getMinutes())}:${pad(now.getSeconds())}.${pad(now.getMilliseconds(), 3)}`;
};

const LEVEL_COLORS: Record<string, string> = {
  INFO: GREEN,
  WARNING: YELLOW,
  ERROR: RED,
};

const stripAnsi = (text: string): string => text.replace(/\[[0-9;]*m/g, '');

const cleanMessage = (text: string): string =>
  stripAnsi(text)
    .replace(/^\d{1,2}:\d{2}:\d{2}(\s?[AP]M)?\s*/, '')
    .replace(/^\[vite\]\s*/, '')
    .trimEnd();

export const emitPretty = (level: keyof typeof LEVEL_COLORS, message: string): void => {
  for (const line of cleanMessage(message).split('\n')) {
    if (line.trim() === '') continue;

    process.stdout.write(
      `${DIM}${timestamp()}${RESET} ${LEVEL_COLORS[level] ?? GREEN}${BOLD}${level.padEnd(7)}${RESET} ` +
        `${BOLD}${WHITE}${line.trim()}${RESET}\n`,
    );
  }
};

export const createPrettyLogger = (): Logger => {
  const base = createLogger();

  return {
    ...base,
    info: (message) => emitPretty('INFO', message),
    warn: (message) => {
      base.hasWarned = true;
      emitPretty('WARNING', message);
    },
    warnOnce: (message) => {
      base.hasWarned = true;
      emitPretty('WARNING', message);
    },
    error: (message) => emitPretty('ERROR', message),
  };
};

export const prettyAccessLog = (): Plugin => ({
  name: 'pretty-access-log',
  apply: 'serve',
  configureServer(server) {
    server.middlewares.use((req, res, next) => {
      const started = process.hrtime.bigint();

      res.on('finish', () => {
        const elapsed = Number(process.hrtime.bigint() - started) / 1_000_000;
        const method = (req.method ?? 'GET').padEnd(6);
        const path = (req.url ?? '/').split('?')[0];
        const status = res.statusCode;
        const size = Number(res.getHeader('content-length') ?? 0);
        const address = req.socket.remoteAddress ?? '-';

        process.stdout.write(
          `${DIM}${timestamp()}${RESET} ${GREEN}${BOLD}INFO   ${RESET} ` +
            `${BOLD}${WHITE}http request${RESET}           ` +
            `${METHOD_COLORS[req.method ?? ''] ?? MAGENTA}${BOLD}${method}${RESET} ` +
            `${CYAN}${path}${RESET} ` +
            `${statusColor(status)}${BOLD}${status}${RESET} ` +
            `${durationColor(elapsed)}${formatDuration(elapsed)}${RESET} ` +
            `${DIM}${formatBytes(size)}${RESET} ` +
            `${DIM}from ${address}${RESET}\n`,
        );
      });

      next();
    });
  },
});
