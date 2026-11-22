// Renders real `afs` sessions as terminal-style SVGs for the README.
// Usage (from the repository root, after `make build`): node docs/images/source/terminal-shots.mjs
import { spawnSync } from 'node:child_process';
import { mkdtempSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';

const AFS = './bin/afs';
const HOME = mkdtempSync(join(tmpdir(), 'afs-shot-')); // empty profile, no credentials
const shots = [
  { file: 'afs-help.svg', title: 'afs help', runs: [['--version'], ['help']], maxLines: 30 },
  {
    file: 'afs-secret-guard.svg',
    title: 'Secret values never go on the command line',
    runs: [['secrets', 'create-secret', '--team-id', 'team_01JD2Q3R4S5T6V7W8X9Y0Z0001', '--name', 'github-pat', '--scope', 'team', '--secret-value', 'hunter2', '--endpoint-url', 'https://afs.example.com']],
    showExit: true,
  },
];

const esc = (s) => s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
for (const shot of shots) {
  const lines = [];
  for (const args of shot.runs) {
    lines.push({ t: `$ afs ${args.map((a) => (/\s/.test(a) ? `'${a}'` : a)).join(' ')}`, prompt: true });
    const r = spawnSync(AFS, args, { env: { HOME, PATH: process.env.PATH, NO_COLOR: '1', AFS_PAGER: '' }, encoding: 'utf8' });
    let out = (r.stdout + r.stderr).replace(/\s+$/, '').split('\n');
    if (shot.maxLines && out.length > shot.maxLines) out = [...out.slice(0, shot.maxLines), '…'];
    out.forEach((t) => lines.push({ t }));
    if (shot.showExit) lines.push({ t: `$ echo $?`, prompt: true }, { t: String(r.status) });
  }
  // wrap long lines at 118 columns
  const COLS = 118;
  const wrapped = lines.flatMap((l) => {
    const parts = []; let s = l.t;
    while (s.length > COLS) { let cut = s.lastIndexOf(' ', COLS); if (cut < 20) cut = COLS; parts.push(s.slice(0, cut)); s = `    ${s.slice(cut).trimStart()}`; }
    parts.push(s);
    return parts.map((t, i) => ({ t, prompt: l.prompt }));
  });
  const lh = 19, W = 980, H = 52 + wrapped.length * lh + 24;
  const body = wrapped.map((l, i) => `<text x="22" y="${62 + i * lh}" style="white-space:pre" fill="${l.prompt ? '#7ee787' : '#e6edf3'}" ${l.prompt ? 'font-weight="700"' : ''}>${esc(l.t)}</text>`).join('');
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="${W}" height="${H}" viewBox="0 0 ${W} ${H}" role="img" aria-label="${esc(shot.title)}"><title>${esc(shot.title)}</title>` +
    `<rect width="${W}" height="${H}" rx="10" fill="#0d1117"/><rect width="${W}" height="34" rx="10" fill="#21262d"/><rect y="24" width="${W}" height="10" fill="#21262d"/>` +
    `<circle cx="20" cy="17" r="6" fill="#ff5f57"/><circle cx="40" cy="17" r="6" fill="#febc2e"/><circle cx="60" cy="17" r="6" fill="#28c840"/>` +
    `<text x="${W / 2}" y="22" text-anchor="middle" font-family="-apple-system, 'Segoe UI', Helvetica, Arial, sans-serif" font-size="12.5" fill="#8b949e">${esc(shot.title)}</text>` +
    `<g font-family="ui-monospace, 'SFMono-Regular', Menlo, Consolas, monospace" font-size="13" style="white-space:pre" xml:space="preserve">${body}</g></svg>\n`;
  writeFileSync(join('docs/images', shot.file), svg);
  console.log('wrote', shot.file);
}
