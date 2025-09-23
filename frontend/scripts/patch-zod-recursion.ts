import fs from 'fs';
import path from 'path';

const apiDir = path.resolve('src/api-client');

function patchFile(filePath: string) {
  let content = fs.readFileSync(filePath, 'utf8');

  console.log(`Patching file: ${filePath}`);

  // naive example: replace ".array(CategorySchema)" with ".array(z.lazy(() => CategorySchema))"
  content = content.replace(/:\s*(z\w+)\s*,/g, ': z.lazy(() => $1),');

  fs.writeFileSync(filePath, content, 'utf8');
}

fs.readdirSync(apiDir)
  .filter((f) => f === 'zod.gen.ts')
  .forEach((f) => patchFile(path.join(apiDir, f)));
