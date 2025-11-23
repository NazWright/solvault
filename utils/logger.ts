// Simple structured logger used across services/handlers.
// Keep this minimal and dependency-free so it's easy to replace with pino/winston later.

export type LogLevel = 'debug' | 'info' | 'warn' | 'error';

const levelOrder: Record<LogLevel, number> = {
  debug: 0,
  info: 1,
  warn: 2,
  error: 3,
};

const DEFAULT_LEVEL: LogLevel = (process.env.SOLVAULT_LOG_LEVEL as LogLevel) || 'info';

function timestamp() {
  return new Date().toISOString();
}

export const logger = {
  level: DEFAULT_LEVEL,

  log(level: LogLevel, message: string, meta?: Record<string, unknown>) {
    if (levelOrder[level] < levelOrder[this.level]) return;
    const out = {
      ts: timestamp(),
      level,
      message,
      ...meta,
    };
    if (level === 'error') {
      console.error(JSON.stringify(out));
    } else {
      console.log(JSON.stringify(out));
    }
  },

  debug(msg: string, meta?: Record<string, unknown>) {
    this.log('debug', msg, meta);
  },
  info(msg: string, meta?: Record<string, unknown>) {
    this.log('info', msg, meta);
  },
  warn(msg: string, meta?: Record<string, unknown>) {
    this.log('warn', msg, meta);
  },
  error(msg: string, meta?: Record<string, unknown>) {
    this.log('error', msg, meta);
  },
};

// logger added by Nazere Wright - atomic commit

Signed-off-by: Nazere Wright <76058043+NazWright@users.noreply.github.com>
