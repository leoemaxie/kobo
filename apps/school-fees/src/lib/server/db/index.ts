import { drizzle } from 'drizzle-orm/postgres-js';
import postgres from 'postgres';
import * as schema from './schema';

const connectionString = process.env.DATABASE_URL || 'postgres://kobo_school_fees:pass@localhost:5432/school_fees';
const queryClient = postgres(connectionString);
export const db = drizzle(queryClient, { schema });

export type Parent = typeof schema.parents.$inferSelect;
export type Session = typeof schema.sessions.$inferSelect;
