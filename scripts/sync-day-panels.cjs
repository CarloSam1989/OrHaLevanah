const fs = require('node:fs');
const path = require('node:path');
const root = path.resolve(__dirname, '..');
for (const target of ['mobile_web', 'www']) {
  fs.copyFileSync(path.join(root, 'web/static/js/day-panels.js'), path.join(root, target, 'day-panels.js'));
}
