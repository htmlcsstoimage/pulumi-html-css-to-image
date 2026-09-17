import fs from 'node:fs';
const schema = JSON.parse(fs.readFileSync('internal/provider/schema.json', 'utf8'));
const version = fs.readFileSync('VERSION', 'utf8').trim();
if (schema.version !== version) throw new Error('Schema and VERSION differ');
const template = fs.readFileSync('packaging/java/pom.xml', 'utf8');
fs.writeFileSync('sdk/java/pom.xml', template.replaceAll('@VERSION@', version));
fs.cpSync('packaging/java/tests', 'sdk/java/src/test/java', { recursive: true });
// Pulumi Java Utilities reads metadata using the schema name, regardless of the Java package override.
const resources = 'sdk/java/src/main/resources/com/htmlcsstoimage/html-css-to-image';
fs.mkdirSync(resources, { recursive: true });
fs.writeFileSync(`${resources}/version.txt`, version);
fs.writeFileSync(`${resources}/plugin.json`, JSON.stringify({
    resource: true, name: schema.name, version, server: schema.pluginDownloadURL,
}, null, 2) + '\n');
