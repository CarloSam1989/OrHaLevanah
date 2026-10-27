const fs = require('node:fs');
const path = require('node:path');

const root = path.resolve(__dirname, '..');
const output = path.join(root, 'dist');

fs.rmSync(output, { recursive: true, force: true });
fs.mkdirSync(path.join(output, 'app'), { recursive: true });
fs.copyFileSync(
  path.join(root, 'web', 'templates', 'community.html'),
  path.join(output, 'index.html')
);
fs.copyFileSync(
  path.join(root, 'web', 'templates', 'index.html'),
  path.join(output, 'app', 'index.html')
);
fs.cpSync(
  path.join(root, 'web', 'static'),
  path.join(output, 'static'),
  { recursive: true }
);
fs.mkdirSync(path.join(output, 'data'), { recursive: true });
fs.copyFileSync(
  path.join(root, 'data', 'moon_sightings.json'),
  path.join(output, 'data', 'moon_sightings.json')
);
