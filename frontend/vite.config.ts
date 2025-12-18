import fs from 'fs';
import path from 'path';
import {defineConfig} from 'vite';
import react from '@vitejs/plugin-react';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    configureServer(server) {
      const iconRoot = resolveIconRoot();
      if (!iconRoot) {
        return;
      }

      server.middlewares.use('/icons', (req, res, next) => {
        if (!req.url) {
          next();
          return;
        }

        const requestPath = req.url.split('?')[0] || '';
        let decodedPath = requestPath;
        try {
          decodedPath = decodeURIComponent(requestPath);
        } catch {
          res.statusCode = 400;
          res.end('Bad Request');
          return;
        }

        const safePath = path.normalize(decodedPath).replace(/^(\.\.(\/|\\|$))+/, '');
        const relativePath = safePath.replace(/^[/\\]+/, '');
        const filePath = path.join(iconRoot, relativePath);

        if (!filePath.startsWith(iconRoot)) {
          res.statusCode = 403;
          res.end('Forbidden');
          return;
        }

        fs.stat(filePath, (err, stat) => {
          if (err || !stat.isFile()) {
            next();
            return;
          }

          res.setHeader('Content-Type', 'image/png');
          res.setHeader('Cache-Control', 'public, max-age=86400');
          fs.createReadStream(filePath).pipe(res);
        });
      });
    },
  },
});

function resolveIconRoot(): string | null {
  const base = process.env.APPDATA || process.env.LOCALAPPDATA;
  if (!base) {
    return null;
  }
  return path.join(base, 'rungrid', 'icons');
}
