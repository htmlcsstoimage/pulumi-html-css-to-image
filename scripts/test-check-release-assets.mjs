import assert from 'node:assert/strict';
import { test } from 'node:test';
import { missingAssets } from './check-release-assets.mjs';

const tag = 'v0.1.0';
const assets = missingAssets([], tag).map(name => ({ name, size: 123, state: 'uploaded' }));

test('a complete draft or published release skips rebuilding', () => {
  assert.equal(assets.length, 8);
  for (const draft of [true, false]) {
    assert.deepEqual(missingAssets([{ tag_name: tag, draft, assets }], tag), []);
  }
});

test('absent releases and other versions cannot satisfy the check', () => {
  assert.equal(missingAssets([], tag).length, 8);
  assert.equal(missingAssets([{ tag_name: 'v0.2.0', assets }], tag).length, 8);
});

test('every archive, schema and checksum file is required', () => {
  for (const asset of assets) {
    const remaining = assets.filter(other => other !== asset);
    assert.deepEqual(missingAssets([{ tag_name: tag, assets: remaining }], tag), [asset.name]);
  }
});

test('empty or unfinished uploads are not complete', () => {
  for (const invalid of [{ size: 0 }, { state: 'starter' }]) {
    const partial = assets.map((asset, index) => index === 0 ? { ...asset, ...invalid } : asset);
    assert.deepEqual(missingAssets([{ tag_name: tag, assets: partial }], tag), [assets[0].name]);
  }
});
