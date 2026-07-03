import { pgTable, text, timestamp, boolean } from 'drizzle-orm/pg-core';

export const parents = pgTable('parents', {
	id: text('id').primaryKey(),
	name: text('name').notNull(),
	email: text('email').notNull().unique(),
	passwordHash: text('password_hash').notNull(),
	isAdmin: boolean('is_admin').default(false).notNull(),
	createdAt: timestamp('created_at').notNull().defaultNow()
});

export const students = pgTable('students', {
	id: text('id').primaryKey(),
	name: text('name').notNull(),
	className: text('class_name').notNull(),
	koboIdentityId: text('kobo_identity_id').notNull(),
	createdAt: timestamp('created_at').notNull().defaultNow()
});

export const parentStudents = pgTable('parent_students', {
	parentId: text('parent_id').notNull().references(() => parents.id),
	studentId: text('student_id').notNull().references(() => students.id),
});

export const sessions = pgTable('sessions', {
	id: text('id').primaryKey(),
	userId: text('user_id').notNull().references(() => parents.id),
	expiresAt: timestamp('expires_at', { withTimezone: true, mode: 'date' }).notNull()
});
