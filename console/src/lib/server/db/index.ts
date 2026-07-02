import { drizzle } from 'drizzle-orm/postgres-js';
import postgres from 'postgres';
import * as schema from './schema';

// This is a placeholder for the actual env variable
const connectionString = process.env.DATABASE_URL || 'postgres://kobo_console_app:pass@localhost:5432/kobo';

const queryClient = postgres(connectionString);
export const db = drizzle(queryClient, { schema });

// Export inferred types for convenience in app.d.ts
export type User = typeof schema.users.$inferSelect;
export type Session = typeof schema.sessions.$inferSelect;
