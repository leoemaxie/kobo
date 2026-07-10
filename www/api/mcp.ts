import docs from '../build/mcp/docs.json';
import searchIndex from '../build/mcp/search-index.json';

// We initialize this lazily to avoid ERR_REQUIRE_ESM in CommonJS environments
let handler: ((req: Request) => Promise<Response>) | null = null;

export default async function(req: Request) {
  if (!handler) {
    const { createWebRequestHandler } = await import('docusaurus-plugin-mcp-server/adapters');
    handler = createWebRequestHandler({
      docs,
      searchIndexData: searchIndex,
      name: 'kobo-docs',
      baseUrl: 'https://kobo.triumphsystems.tech',
      instructions: 'Search the Kobo API documentation and Virtual Account Infrastructure docs. Use docs_search to find pages, then docs_fetch for full content.',
    });
  }
  return handler(req);
}
