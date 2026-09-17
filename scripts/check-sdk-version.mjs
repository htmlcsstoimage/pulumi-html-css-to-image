import fs from 'node:fs';
const language = process.argv[2];
const version = fs.readFileSync('VERSION', 'utf8').trim();
const schema = JSON.parse(fs.readFileSync('internal/provider/schema.json', 'utf8'));
if (schema.version !== version) throw new Error('Schema and VERSION differ');
let actual;
let expected = version;
switch (language) {
    case 'nodejs': actual = JSON.parse(fs.readFileSync('sdk/nodejs/package.json')).version; break;
    case 'python':
        actual = fs.readFileSync('sdk/python/setup.py', 'utf8').match(/^VERSION = "([^"]+)"/m)?.[1];
        expected = version.replace(/-(alpha|beta|rc)\.(\d+)$/, (_, kind, n) => ({alpha: 'a', beta: 'b', rc: 'rc'}[kind]) + n);
        break;
    case 'dotnet':
        actual = fs.readFileSync('sdk/dotnet/Pulumi.HtmlCssToImage.csproj', 'utf8').match(/<Version>([^<]+)<\/Version>/)?.[1];
        if (fs.readFileSync('sdk/dotnet/version.txt', 'utf8').trim() !== version) throw new Error('.NET plugin version differs');
        break;
    case 'java': actual = fs.readFileSync('sdk/java/pom.xml', 'utf8').match(/<version>([^<]+)<\/version>/)?.[1]; break;
    case 'go': actual = JSON.parse(fs.readFileSync('sdk/go/pulumi-plugin.json')).version; break;
    default: throw new Error(`Unknown SDK: ${language}`);
}
if (actual !== expected) throw new Error(`${language} version ${actual} does not match ${expected}`);
