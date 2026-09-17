import { copyFileSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';

const language = process.argv[2];
if (!['nodejs', 'python', 'go', 'dotnet', 'java'].includes(language)) {
  throw new Error('Specify nodejs, python, go, dotnet, or java');
}
const repository = 'https://github.com/htmlcsstoimage/pulumi-html-css-to-image';
let example;
if (['nodejs', 'python'].includes(language)) {
  example = readFileSync(`packaging/${language}/example.md`, 'utf8');
} else {
  const codeLanguage = language === 'dotnet' ? 'csharp' : language;
  const docs = readFileSync('docs/_index.md', 'utf8');
  const code = docs.match(new RegExp('```' + codeLanguage + '\\n[\\s\\S]*?```'))?.[0];
  if (!code) throw new Error(`Missing ${codeLanguage} example in docs/_index.md`);
  example = `## Example\n\nIn a Pulumi project for your language, install the SDK and configure credentials:\n\n` +
    '```sh\npulumi config set html-css-to-image:apiId YOUR_API_ID\npulumi config set --secret html-css-to-image:apiKey\n```\n\n' +
    code + '\n\nRun `pulumi preview` and `pulumi up`. Open the returned URL to render the image. Creation and refresh save/read metadata without rendering bytes. Input changes replace image definitions; template edits create new versions under a stable ID.\n';
}
const readme = readFileSync('README.md', 'utf8');
if (!readme.includes('## Example\n') || !readme.includes('## Configuration\n')) {
  throw new Error('README Example/Configuration sections are required');
}
const content = readme
  .replace(/## Example\n[\s\S]*?(?=## Configuration\n)/, `${example.trim()}\n\n`)
  // Registry pages cannot resolve repository-relative documentation links.
  .replace(/\]\((docs\/[^)]+|examples)\)/g, `](${repository}/tree/main/$1)`);
writeFileSync(`sdk/${language}/README.md`, content);
if (language === 'dotnet') {
  copyFileSync('packaging/dotnet/Directory.Build.props', 'sdk/dotnet/Directory.Build.props');
}
if (language === 'java') {
  const resources = 'sdk/java/src/main/resources/META-INF';
  mkdirSync(resources, { recursive: true });
  writeFileSync(`${resources}/README.md`, content);
}
