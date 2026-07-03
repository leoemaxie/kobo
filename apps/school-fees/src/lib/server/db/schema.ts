import {
  pgSchema,
  uuid,
  text,
  timestamp,
  boolean,
  uniqueIndex,
} from 'drizzle-orm/pg-core';

export const schoolFeesSchema = pgSchema('school_fees');

export const parents = schoolFeesSchema.table(
  'parents',
  {
    id: uuid('id').primaryKey().defaultRandom(),
    email: text('email').notNull().unique(),
    passwordHash: text('password_hash').notNull(),
    isAdmin: boolean('is_admin').notNull().default(false),
    createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
  },
  (table) => ({
    emailIdx: uniqueIndex('idx_parents_email').on(table.email),
  })
);

export const students = schoolFeesSchema.table('students', {
  id: uuid('id').primaryKey().defaultRandom(),
  name: text('name').notNull(),
  className: text('class_name').notNull(),
  koboIdentityId: text('kobo_identity_id').notNull().unique(),
  createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
});

export const parentStudents = schoolFeesSchema.table('parent_students', {
  id: uuid('id').primaryKey().defaultRandom(),
  parentId: uuid('parent_id')
    .notNull()
    .references(() => parents.id),
  studentId: uuid('student_id')
    .notNull()
    .references(() => students.id),
  createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
});

export const sessions = schoolFeesSchema.table('sessions', {
  id: text('id').primaryKey(),
  parentId: uuid('parent_id')
    .notNull()
    .references(() => parents.id),
  expiresAt: timestamp('expires_at', { withTimezone: true }).notNull(),
  createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
  revokedAt: timestamp('revoked_at', { withTimezone: true }),
});
