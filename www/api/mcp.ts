import { createWebRequestHandler } from 'docusaurus-plugin-mcp-server/adapters';
import docs from '../build/mcp/docs.json';
import searchIndex from '../build/mcp/search-index.json';



const handler = createWebRequestHandler({
  docs,
  searchIndexData: searchIndex,
  name: 'kobo-docs',
  baseUrl: 'https://kobo.triumphsystems.tech',
  instructions: 'Search the Kobo API documentation and Virtual Account Infrastructure docs. Use docs_search to find pages, then docs_fetch for full content.',
});

export default function(req: Request) {
  return handler(req);
}
