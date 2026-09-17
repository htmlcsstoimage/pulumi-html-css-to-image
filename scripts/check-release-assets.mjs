import { execFileSync } from 'node:child_process';
import { appendFileSync } from 'node:fs';
import { pathToFileURL } from 'node:url';

export function missingAssets(releases, tag) {
  const release = releases.find(release => release.tag_name === tag);
  const prefix = `pulumi-resource-html-css-to-image-${tag}`;
  const expected = ['schema.json', `${prefix}-checksums.txt`];
  for (const os of ['linux', 'darwin', 'windows']) {
    for (const arch of ['amd64', 'arm64']) {
      expected.push(`${prefix}-${os}-${arch}.tar.gz`);
    }
  }
  return expected.filter(name => !release?.assets?.some(asset =>
    asset.name === name && asset.state === 'uploaded' && asset.size > 0));
}

if (process.argv[1] && import.meta.url === pathToFileURL(process.argv[1]).href) {
  const { GITHUB_REPOSITORY, GITHUB_OUTPUT, TAG } = process.env;
  if (!GITHUB_REPOSITORY || !GITHUB_OUTPUT || !TAG) {
    throw new Error('GITHUB_REPOSITORY, GITHUB_OUTPUT, and TAG are required');
  }
  // Listing releases includes drafts. API/auth failures must fail this check,
  // rather than being mistaken for an absent release.
  const pages = JSON.parse(execFileSync('gh', [
    'api', '--paginate', '--slurp',
    `repos/${GITHUB_REPOSITORY}/releases?per_page=100`,
  ], { encoding: 'utf8' }));
  const missing = missingAssets(pages.flat(), TAG);
  const complete = missing.length === 0;
  console.log(complete
    ? `${TAG} already has all eight uploaded release assets; skipping GoReleaser.`
    : `${TAG} needs release assets: ${missing.join(', ')}`);
  appendFileSync(GITHUB_OUTPUT, `complete=${complete}\n`);
}
