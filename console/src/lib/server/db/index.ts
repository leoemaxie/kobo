import { drizzle } from 'drizzle-orm/postgres-js';
import postgres from 'postgres';
import { env } from '$env/dynamic/private';
import * as schema from './schema';

const connectionString = env.DATABASE_URL;

if (!connectionString) {
  throw new Error('DATABASE_URL is not set in the environment variables');
}

const queryClient = postgres(connectionString, {
  ssl: 'require',
  max: 1,
  connection: {
    search_path: 'public,console',
  },
});
export const db = drizzle(queryClient, { schema });

export type User = typeof schema.users.$inferSelect;
export type Session = typeof schema.sessions.$inferSelect;
