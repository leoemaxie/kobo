import test from 'node:test';
import assert from 'node:assert';
import { createServer, IncomingMessage, ServerResponse } from 'node:http';
import { KoboClient } from './client.js';

test('KoboClient', async (t: any) => {
  await t.test('Client construction', () => {
    const client = new KoboClient('pk_test', 'sk_test');
    assert.ok(client);
    
    const sandboxClient = KoboClient.sandbox('pk_test', 'sk_test');
    assert.ok(sandboxClient);
  });

  await t.test('Health check', async () => {
    const server = createServer((req: IncomingMessage, res: ServerResponse) => {
      assert.strictEqual(req.url, '/healthz');
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ status: 'ok', db: 'ok' }));
    });

    await new Promise<void>((resolve) => server.listen(0, resolve));
    const port = (server.address() as any).port;
    const baseUrl = `http://localhost:${port}`;

    const client = new KoboClient('pk', 'sk', { baseUrl });
    const response = await client.health();
    
    assert.strictEqual(response.status, 'ok');
    assert.strictEqual(response.db, 'ok');

    server.close();
  });

  await t.test('Identities create', async () => {
    const server = createServer((req: IncomingMessage, res: ServerResponse) => {
      assert.strictEqual(req.method, 'POST');
      assert.strictEqual(req.url, '/identities');
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ id: 'id_123', state: 'pending', display_name: 'Test User' }));
    });

    await new Promise<void>((resolve) => server.listen(0, resolve));
    const port = (server.address() as any).port;
    const baseUrl = `http://localhost:${port}`;

    const client = new KoboClient('pk', 'sk', { baseUrl });
    const identity = await client.identities.create({
      display_name: 'Test User',
      external_reference: 'ext_123'
    });
    
    assert.strictEqual(identity.id, 'id_123');
    assert.strictEqual(identity.display_name, 'Test User');
    assert.strictEqual(identity.state, 'pending');

    server.close();
  });
});
