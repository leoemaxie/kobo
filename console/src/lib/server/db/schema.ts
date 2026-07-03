import {
  pgSchema,
  uuid,
  text,
  timestamp,
  boolean,
  bigint,
  jsonb,
  uniqueIndex,
  index,
} from 'drizzle-orm/pg-core';

export const consoleSchema = pgSchema('console');

export const webhookStatusEnum = consoleSchema.enum('webhook_status', ['active', 'disabled']);

export const environmentEnum = consoleSchema.enum('environment', ['sandbox', 'production']);
export const integratorStatusEnum = consoleSchema.enum('integrator_status', [
  'active',
  'suspended',
]);

export const apiIntegrators = consoleSchema.table('api_integrators', {
  id: uuid('id').primaryKey().defaultRandom(),
  name: text('name').notNull(),
  status: integratorStatusEnum('status').notNull().default('active'),
  productionAccessGranted: boolean('production_access_granted').notNull().default(false),
  productionAccessGrantedAt: timestamp('production_access_granted_at', { withTimezone: true }),
  productionAccessGrantedBy: uuid('production_access_granted_by'),
  createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
  updatedAt: timestamp('updated_at', { withTimezone: true }).notNull().defaultNow(),
});

export const apiCredentials = consoleSchema.table(
  'api_credentials',
  {
    id: uuid('id').primaryKey().defaultRandom(),
    integratorId: uuid('integrator_id').notNull().references(() => apiIntegrators.id),
    environment: environmentEnum('environment').notNull(),
    keyId: text('key_id').notNull().unique(),
    secretHash: text('secret_hash').notNull(),
    label: text('label'),
    createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
    createdBy: uuid('created_by').notNull(),
    rotatedAt: timestamp('rotated_at', { withTimezone: true }),
    revokedAt: timestamp('revoked_at', { withTimezone: true }),
    revokedBy: uuid('revoked_by'),
    revokedReason: text('revoked_reason'),
  },
  (table) => ({
    keyIdIdx: uniqueIndex('idx_api_credentials_key_id').on(table.keyId),
    integratorEnvIdx: index('idx_api_credentials_integrator_env').on(
      table.integratorId,
      table.environment
    ),
  })
);

export const userRoleEnum = consoleSchema.enum('user_role', ['owner', 'member', 'superadmin']);

export const users = consoleSchema.table(
  'users',
  {
    id: uuid('id').primaryKey().defaultRandom(),
    integratorId: uuid('integrator_id').references(() => apiIntegrators.id),
    email: text('email').notNull().unique(),
    passwordHash: text('password_hash').notNull(),
    role: userRoleEnum('role').notNull().default('member'),
    emailVerifiedAt: timestamp('email_verified_at', { withTimezone: true }),
    createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
    updatedAt: timestamp('updated_at', { withTimezone: true }).notNull().defaultNow(),
  },
  (table) => ({
    emailIdx: uniqueIndex('idx_users_email').on(table.email),
  })
);

export const sessions = consoleSchema.table(
  'sessions',
  {
    id: text('id').primaryKey(),
    userId: uuid('user_id').notNull().references(() => users.id),
    expiresAt: timestamp('expires_at', { withTimezone: true }).notNull(),
    createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
    revokedAt: timestamp('revoked_at', { withTimezone: true }),
  },
  (table) => ({
    userIdx: index('idx_sessions_user').on(table.userId),
  })
);

export const emailVerificationTokens = consoleSchema.table('email_verification_tokens', {
  id: text('id').primaryKey(),
  userId: uuid('user_id').notNull().references(() => users.id),
  expiresAt: timestamp('expires_at', { withTimezone: true }).notNull(),
  usedAt: timestamp('used_at', { withTimezone: true }),
  createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
});

export const billingRecords = consoleSchema.table(
  'billing_records',
  {
    id: uuid('id').primaryKey().defaultRandom(),
    integratorId: uuid('integrator_id').notNull().references(() => apiIntegrators.id),
    environment: environmentEnum('environment').notNull(),
    period: text('period').notNull(),
    accountsProvisioned: bigint('accounts_provisioned', { mode: 'number' }).notNull().default(0),
    transactionsProcessed: bigint('transactions_processed', { mode: 'number' }).notNull().default(0),
    amountDueKobo: bigint('amount_due_kobo', { mode: 'number' }).notNull().default(0),
    adjustmentKobo: bigint('adjustment_kobo', { mode: 'number' }).notNull().default(0),
    adjustmentReason: text('adjustment_reason'),
    syncedAt: timestamp('synced_at', { withTimezone: true }).notNull().defaultNow(),
  },
  (table) => ({
    integratorPeriodIdx: uniqueIndex('idx_billing_integrator_period_env').on(
      table.integratorId,
      table.period,
      table.environment
    ),
  })
);

export const webhooks = consoleSchema.table('webhooks', {
  id: uuid('id').primaryKey().defaultRandom(),
  integratorId: uuid('integrator_id').notNull().references(() => apiIntegrators.id),
  environment: environmentEnum('environment').notNull(),
  url: text('url').notNull(),
  secret: text('secret').notNull(),
  status: webhookStatusEnum('status').notNull().default('active'),
  events: jsonb('events').notNull().default([]),
  createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
});

export const adminActionEnum = consoleSchema.enum('admin_action', [
  'integrator_suspended',
  'integrator_reinstated',
  'production_access_granted',
  'credential_force_revoked',
  'billing_adjustment_applied',
  'user_session_revoked',
]);

export const adminAuditLog = consoleSchema.table('admin_audit_log', {
  id: uuid('id').primaryKey().defaultRandom(),
  actorUserId: uuid('actor_user_id').notNull().references(() => users.id),
  action: adminActionEnum('action').notNull(),
  targetIntegratorId: uuid('target_integrator_id').references(() => apiIntegrators.id),
  targetUserId: uuid('target_user_id').references(() => users.id),
  targetCredentialId: uuid('target_credential_id').references(() => apiCredentials.id),
  detail: jsonb('detail').notNull().default({}),
  createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
});
