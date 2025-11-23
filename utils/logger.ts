// Simple logger utility for SolVault CLI
export const logger = {
  info: (message: string, meta?: Record<string, any>) => {
    console.log(`[INFO] ${message}`, meta ? JSON.stringify(meta) : '');
  },
  error: (message: string, meta?: Record<string, any>) => {
    console.error(`[ERROR] ${message}`, meta ? JSON.stringify(meta) : '');
  },
  warn: (message: string, meta?: Record<string, any>) => {
    console.warn(`[WARN] ${message}`, meta ? JSON.stringify(meta) : '');
  },
  debug: (message: string, meta?: Record<string, any>) => {
    console.debug(`[DEBUG] ${message}`, meta ? JSON.stringify(meta) : '');
  }
};
