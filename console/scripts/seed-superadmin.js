import { eq } from 'drizzle-orm';
import { drizzle } from 'drizzle-orm/postgres-js';
import postgres from 'postgres';
import * as schema from '../src/lib/server/db/schema.js';
import * as argon2 from 'argon2';
import readline from 'readline/promises';
import { stdin as input, stdout as output } from 'process';
import fs from 'fs';
import path from 'path';

// Minimal .env parser for the script
const envPath = path.join(process.cwd(), '.env');
if (fs.existsSync(envPath)) {
  const envContent = fs.readFileSync(envPath, 'utf8');
  envContent.split('\n').forEach(line => {
    const match = line.match(/^\s*([\w.-]+)\s*=\s*(.*)?\s*$/);
    if (match) {
      const key = match[1];
      let value = match[2] || '';
      if (value.startsWith('"') && value.endsWith('"')) {
         value = value.slice(1, -1);
      } else if (value.startsWith("'") && value.endsWith("'")) {
         value = value.slice(1, -1);
      }
      if (!process.env[key]) {
        process.env[key] = value;
      }
    }
  });
}

const connectionString = process.env.DATABASE_URL || 'postgres://kobo_console_app:pass@localhost:5432/kobo';

const client = postgres(connectionString);
const db = drizzle(client, { schema });

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
  _         _          
 | |       | |         
 | | _____ | |__   ___ 
 | |/ / _ \\| '_ \\ / _ \\
 |   < (_) | |_) | (_) |
 |_|\\_\\___/|_.__/ \\___/ 
                        
 Console Superadmin Seeder
${c.reset}
`;

const askPassword = (query) => {
  return new Promise((resolve) => {
    output.write(query);
    input.setRawMode(true);
    input.resume();
    input.setEncoding('utf8');

    let password = '';

    const onData = (char) => {
      if (char === '\n' || char === '\r' || char === '\u0004') {
        input.removeListener('data', onData);
        input.setRawMode(false);
        input.pause();
        output.write('\n');
        resolve(password);
        return;
      }
      if (char === '\u0003') {
        input.removeListener('data', onData);
        input.setRawMode(false);
        input.pause();
        output.write('\n');
        process.exit(1);
      }
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
  
  const rl = readline.createInterface({ input, output });

  try {
    console.log(`${c.cyan}Initialize your root superadmin account.${c.reset}\n`);

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
    
    // Check if email already exists
    const existingUsers = await db.select().from(schema.users).where(eq(schema.users.email, email));
    
    if (existingUsers.length > 0) {
      const user = existingUsers[0];
      if (user.role === 'superadmin') {
        console.log(`${c.green}✔ User ${email} is already a superadmin.${c.reset}\n`);
      } else {
        console.log(`${c.yellow}⚠ User ${email} already exists. Upgrading to superadmin...${c.reset}`);
        await db.update(schema.users)
          .set({ role: 'superadmin' })
          .where(eq(schema.users.email, email));
        console.log(`${c.green}✔ Successfully upgraded ${email} to superadmin.${c.reset}\n`);
      }
    } else {
      console.log(`${c.dim}◌ Hashing password...${c.reset}`);
      const passwordHash = await argon2.hash(password);
      
      console.log(`${c.dim}◌ Creating superadmin record...${c.reset}`);
      await db.insert(schema.users).values({
        email,
        passwordHash,
        role: 'superadmin',
        emailVerifiedAt: new Date(),
      });
      console.log(`${c.green}✔ Successfully created superadmin account for ${email}!${c.reset}\n`);
    }

    console.log(`${c.magenta}${c.bold}You can now log in to the Kobo Console.${c.reset}\n`);

  } catch (error) {
    console.error(`\n${c.red}✖ An error occurred during seeding:${c.reset}`);
    console.error(error);
  } finally {
    rl.close();
    await client.end();
  }
}

main();
