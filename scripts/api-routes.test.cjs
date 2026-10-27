const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');

const root = path.resolve(__dirname, '..');
const routes = fs.readFileSync(path.join(root, 'internal/api/routes.go'), 'utf8');
const webApp = fs.readFileSync(path.join(root, 'web/static/js/app.js'), 'utf8');
const sharedPanels = fs.readFileSync(path.join(root, 'web/static/js/day-panels.js'), 'utf8');
const server = fs.readFileSync(path.join(root, 'cmd/server/main.go'), 'utf8');
const vercelBuild = fs.readFileSync(path.join(root, 'scripts/build-vercel.cjs'), 'utf8');

test('web calendar endpoints are registered by the Go API', () => {
  const registered = new Set([...routes.matchAll(/HandleFunc\("(\/api\/[^"]+)"/g)].map(match => match[1]));
  const requested = [
    '/api/biblical/jerusalem/sources',
    '/api/moon/jerusalem/sightings',
    '/api/feasts/jerusalem',
    '/api/biblical/jerusalem/month',
    '/api/biblical/jerusalem/calendar',
  ];

  for (const endpoint of requested) {
    assert.ok(registered.has(endpoint), `Go route missing: ${endpoint}`);
  }
  assert.equal(registered.size, 17);
});

test('stored moon sightings are available to the static app and Go server', () => {
  assert.match(vercelBuild, /data['"],\s*'moon_sightings\.json'/);
  assert.match(server, /HandleFunc\("\/data\/moon_sightings\.json"/);
  assert.ok(webApp.includes('/data/moon_sightings.json'));
  assert.ok(webApp.includes('/api/moon/jerusalem/sightings'));
});

test('frontend API paths use the registered /api prefix and expected route names', () => {
  for (const endpoint of [
    '/api/feasts/jerusalem',
    '/api/biblical/jerusalem/month',
    '/api/biblical/jerusalem/calendar',
    '/api/moon/jerusalem/sightings',
  ]) {
    assert.ok(webApp.includes(endpoint), `Frontend endpoint missing: ${endpoint}`);
  }
  assert.ok(sharedPanels.includes('/api/biblical/jerusalem/sources'));
  assert.doesNotMatch(webApp, /["'`]\/(?:app\/)?api\/calendar/);
});
