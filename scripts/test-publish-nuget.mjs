import { test } from 'node:test';
import assert from 'node:assert/strict';
import { mkdtempSync, mkdirSync, writeFileSync, existsSync, readFileSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { resolve, join } from 'node:path';
import { spawnSync } from 'node:child_process';

const script = resolve('scripts/publish-nuget.sh');
for (const [status, pushExit, expectedExit, pushed] of [
  ['200', 0, 0, false], ['404', 0, 0, true],
  ['404', 1, 1, true], ['403', 0, 1, false], ['500', 0, 1, false],
  ['network-error', 0, 6, false],
]) {
  test(`NuGet lookup ${status}, push exit ${pushExit}`, () => {
    const dir = mkdtempSync(join(tmpdir(), 'hcti-nuget-test-'));
    try {
      mkdirSync(join(dir, 'sdk/dotnet/artifacts'), { recursive: true });
      writeFileSync(join(dir, 'VERSION'), '0.1.1\n');
      writeFileSync(join(dir, 'sdk/dotnet/artifacts/HtmlCssToImage.Pulumi.0.1.1.nupkg'), 'fixture');
      writeFileSync(join(dir, 'curl'), '#!/bin/sh\n[ "$STATUS" != network-error ] || exit 6\nprintf "%s" "$STATUS"\n', { mode: 0o755 });
      writeFileSync(join(dir, 'dotnet'), '#!/bin/sh\nprintf "%s\\n" "$@" > pushed\nexit "$PUSH_EXIT"\n', { mode: 0o755 });
      const result = spawnSync('bash', [script], {
        cwd: dir, encoding: 'utf8',
        env: { ...process.env, PATH: `${dir}:${process.env.PATH}`, STATUS: status, PUSH_EXIT: String(pushExit), NUGET_API_KEY: 'test-token' },
      });
      assert.equal(result.status, expectedExit, result.stderr);
      assert.equal(existsSync(join(dir, 'pushed')), pushed);
      if (pushed) assert.ok(!readFileSync(join(dir, 'pushed'), 'utf8').includes('--skip-duplicate'));
    } finally {
      rmSync(dir, { recursive: true, force: true });
    }
  });
}
