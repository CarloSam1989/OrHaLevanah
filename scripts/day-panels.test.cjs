const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const path = require('node:path');
const root = path.resolve(__dirname, '..');

class Element {
  constructor(tag = 'section', text = '') { this.tag = tag; this.textContent = text; this.children = []; this.hidden = false; }
  append(...children) { this.children.push(...children); }
  replaceChildren(...children) { this.children = children; this.textContent = ''; }
  text() { return this.textContent + this.children.map(child => child.text()).join(' '); }
}
function setup(fetch) {
  const errors = [];
  const context = vm.createContext({ fetch, URL, Intl, Date, console: { error: (...args) => errors.push(args) }, document: { createElement: tag => new Element(tag) } });
  vm.runInContext(fs.readFileSync(path.join(root, 'web/static/js/day-panels.js'), 'utf8'), context);
  return { panels: context.DayPanels, errors };
}
function response(date, isHead, sources = []) {
  return { ok: true, json: async () => ({ success: true, data: { date, is_month_start: isHead, sources } }) };
}

test('ordinary day hides sources; a head without evidence explains the empty state', async () => {
  const { panels, errors } = setup(async url => response(new URL(url, 'https://local.test').searchParams.get('date'), url.includes('08-14')));
  const el = new Element();
  await panels.sources(el, '2026-08-14');
  assert.equal(el.hidden, false);
  assert.match(el.text(), /No existen fuentes registradas para este período/);
  await panels.sources(el, '2026-08-15');
  assert.equal(el.hidden, true);
  assert.equal(errors.length, 0);
});

test('sources distinguish estimates and safely render untrusted names and links', async () => {
  const { panels, errors } = setup(async () => response('2026-08-14', true, [
    { name: '<img onerror=alert(1)>', url: 'javascript:alert(1)', status: 'confirmed' },
    { name: 'NASA', url: 'https://example.test/forecast', status: 'estimated' }
  ]));
  const el = new Element();
  await panels.sources(el, '2026-08-14');
  assert.equal(el.children[1].children[0].tag, 'strong');
  assert.equal(el.children[1].children[0].textContent, '<img onerror=alert(1)>');
  assert.equal(el.children[2].children[0].rel, 'noopener noreferrer');
  assert.match(el.text(), /Pronóstico, pendiente de confirmación/);
  assert.equal(errors.length, 0);
});

test('rapid date changes discard a late response', async () => {
  let resolveOld;
  const { panels } = setup(url => url.includes('08-14') ? new Promise(resolve => { resolveOld = resolve; }) : Promise.resolve(response('2026-08-15', false)));
  const el = new Element();
  const first = panels.sources(el, '2026-08-14', '', true);
  await panels.sources(el, '2026-08-15');
  resolveOld(response('2026-08-14', true, [{ name: 'Old result' }]));
  await first;
  assert.equal(el.hidden, true);
  assert.doesNotMatch(el.text(), /Old result/);
});

test('storage/network error is not presented as an empty period', async () => {
  const { panels, errors } = setup(async () => ({ ok: false, status: 500 }));
  const el = new Element();
  await panels.sources(el, '2026-08-14');
  assert.match(el.text(), /No se pudieron consultar/);
  assert.doesNotMatch(el.text(), /No existen fuentes/);
  assert.equal(errors.length, 1);
});

test('changing calendar month cancels a pending source display without another request', async () => {
  let finish;
  let calls = 0;
  const { panels } = setup(() => { calls++; return new Promise(resolve => { finish = resolve; }); });
  const el = new Element();
  const request = panels.sources(el, '2026-08-14', '', true);
  await panels.sources(el, null);
  finish(response('2026-08-14', true, [{ name: 'Previous month' }]));
  await request;
  assert.equal(el.hidden, true);
  assert.doesNotMatch(el.text(), /Previous month/);
  assert.equal(calls, 1);
});

test('today follows the server date and updates after a date change', () => {
  const { panels } = setup();
  const el = new Element();
  panels.today(el, { civil_date: '2026-09-28', biblical_day: 17, sunset_time: '18:30', after_sunset: false });
  assert.match(el.text(), /lunes, 28 de septiembre de 2026/);
  assert.match(el.text(), /Día bíblico 17/);
  panels.today(el, { civil_date: '2026-09-29', biblical_day: 18, after_sunset: true }, { current_feasts: [{ name: 'Fiesta de prueba' }], omer: { is_omer_counting: true, omer_day_today: 4 } });
  assert.match(el.text(), /martes, 29 de septiembre de 2026/);
  assert.match(el.text(), /Después del atardecer/);
  assert.match(el.text(), /Omer: 4/);
  assert.doesNotMatch(el.text(), /lunes/);
  panels.today(el, { civil_date: '2026-09-29' }, { civil_date: '2026-09-28', current_feasts: [{ name: 'Old feast' }] });
  assert.doesNotMatch(el.text(), /Old feast/);
  assert.match(el.text(), /No se pudieron consultar/);
});

test('parallel consumers share one request without keeping stale daily data', async () => {
  let calls = 0;
  const { panels } = setup(async () => { calls++; return response('2026-08-14', true); });
  await Promise.all([panels.json('/same'), panels.json('/same')]);
  assert.equal(calls, 1);
  await panels.json('/same');
  assert.equal(calls, 2);
});

test('web and mobile use identical shared behavior; bundles include it before app.js', () => {
  const canonical = fs.readFileSync(path.join(root, 'web/static/js/day-panels.js'), 'utf8');
  for (const folder of ['mobile_web', 'www']) {
    assert.equal(fs.readFileSync(path.join(root, folder, 'day-panels.js'), 'utf8'), canonical);
    const html = fs.readFileSync(path.join(root, folder, 'index.html'), 'utf8');
    assert.ok(html.indexOf('src="day-panels.js"') < html.indexOf('src="app.js"'));
    assert.match(html, /id="monthSources"[^>]*hidden/);
  }
});
