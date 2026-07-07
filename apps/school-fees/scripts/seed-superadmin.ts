import { drizzle } from 'drizzle-orm/postgres-js';
import postgres from 'postgres';
import * as argon2 from 'argon2';
import 'dotenv/config';
import { pgTable, text, timestamp, boolean } from 'drizzle-orm/pg-core';
import { eq } from 'drizzle-orm';
import * as readline from 'node:readline/promises';
import { stdin as input, stdout as output } from 'node:process';

// Basic schema definition for seeding without relying on relative path aliases
const parents = pgTable('parents', {
	id: text('id').primaryKey(),
	name: text('name').notNull(),
	email: text('email').notNull().unique(),
	passwordHash: text('password_hash').notNull(),
	role: text('role', { enum: ['parent', 'admin', 'superadmin'] }).default('parent').notNull(),
	status: text('status', { enum: ['pending', 'active', 'revoked'] }).default('active').notNull(),
	scope: text('scope').default('read'),
	createdAt: timestamp('created_at').notNull().defaultNow()
});

// ANSI Color Codes
const c = {
  reset: "\x1b[0m",
  bold: "\x1b[1m",
  dim: "\x1b[2m",
  cyan: "\x1b[36m",
  green: "\x1b[32m",
  red: "\x1b[31m",
  yellow: "\x1b[33m",
  magenta: "\x1b[35m"
};

const logo = `
${c.cyan}${c.bold}
  _______   _                       _     
 |__   __| (_)                     | |    
    | |_ __ _ _   _ _ __ ___  _ __ | |__  
    | | '__| | | | | '_ \` _ \\| '_ \\| '_ \\ 
    | | |  | | |_| | | | | | | |_) | | | |
    |_|_|  |_|\\__,_|_| |_| |_| .__/|_| |_|
                             | |          
                             |_|          
  School Fees Superadmin Seeder
${c.reset}
`;

const askPassword = (query: string): Promise<string> => {
  return new Promise((resolve) => {
    output.write(query);
    if (input.setRawMode) {
      input.setRawMode(true);
    }
    input.resume();
    input.setEncoding('utf8');

    let password = '';

    const onData = (char: string) => {
      if (char === '\n' || char === '\r' || char === '\u0004') {
        input.removeListener('data', onData);
        if (input.setRawMode) {
            input.setRawMode(false);
        }
        input.pause();
        output.write('\n');
        resolve(password);
        return;
      }
      // Handle Ctrl+C
      if (char === '\u0003') {
        input.removeListener('data', onData);
        if (input.setRawMode) {
            input.setRawMode(false);
        }
        input.pause();
        output.write('\n');
        process.exit(1);
      }
      // Handle backspace
      if (char === '\b' || char === '\x7f') {
        if (password.length > 0) {
          password = password.slice(0, -1);
          output.write('\b \b');
        }
        return;
      }
      password += char;
      output.write('*');
    };

    input.on('data', onData);
  });
};

async function main() {
    console.clear();
    console.log(logo);

    const databaseUrl = process.env.DATABASE_URL;
    if (!databaseUrl) {
        console.error(`${c.red}✖ ERROR: DATABASE_URL environment variable is missing.${c.reset}`);
        process.exit(1);
    }

    const client = postgres(databaseUrl);
    const db = drizzle(client);

    const rl = readline.createInterface({ input, output });

    try {
        console.log(`${c.cyan}Initialize your school-fees superadmin account.${c.reset}\n`);

        const email = await rl.question(`${c.bold}❯ Enter superadmin email: ${c.reset}`);
        if (!email || !email.includes('@')) {
            console.log(`\n${c.red}✖ Invalid email format. Aborting.${c.reset}\n`);
            process.exit(1);
        }

        const password = await askPassword(`${c.bold}❯ Enter strong password: ${c.reset}`);
        
        if (!password || password.length < 8) {
            console.log(`\n${c.red}✖ Password must be at least 8 characters long. Aborting.${c.reset}\n`);
            process.exit(1);
        }

        console.log(`\n${c.dim}◌ Checking database...${c.reset}`);
        
        const existing = await db.select().from(parents).where(eq(parents.email, email)).limit(1);
        
        if (existing.length > 0) {
            const user = existing[0];
            if (user.role === 'superadmin') {
                console.log(`${c.green}✔ User ${email} is already a superadmin.${c.reset}\n`);
            } else {
                console.log(`${c.yellow}⚠ User ${email} already exists. Upgrading to superadmin...${c.reset}`);
                await db.update(parents)
                  .set({ role: 'superadmin', status: 'active', scope: 'full' })
                  .where(eq(parents.email, email));
                console.log(`${c.green}✔ Successfully upgraded ${email} to superadmin.${c.reset}\n`);
            }
        } else {
            console.log(`${c.dim}◌ Hashing password...${c.reset}`);
            const passwordHash = await argon2.hash(password);
            const id = crypto.randomUUID();

            console.log(`${c.dim}◌ Creating superadmin record...${c.reset}`);
            await db.insert(parents).values({
                id,
                name: 'System Superadmin',
                email,
                passwordHash,
                role: 'superadmin',
                status: 'active',
                scope: 'full'
            });

            console.log(`${c.green}✔ Successfully created superadmin account for ${email}!${c.reset}\n`);
        }

        console.log(`${c.magenta}${c.bold}You can now log in to the School Fees Dashboard.${c.reset}\n`);

    } catch (error) {
        console.error(`\n${c.red}✖ An error occurred during seeding:${c.reset}`);
        console.error(error);
    } finally {
        rl.close();
        await client.end();
        process.exit(0);
    }
}

main();
