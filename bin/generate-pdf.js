#!nodejs
const YAML = require('yaml'),
  CDP = require('chrome-remote-interface'),
  fs = require('fs')

// args are: nodejs index.js config.yaml
if (process.argv.length !== 2) {
  console.log("generate-pdf.js config.yaml")
  process.exit(1)
}

// Load config
const config = YAML.parse(fs.readFileSync(process.argv[2], 'utf8')),
  tasks = config.modules;

console.log("Connecting to chromium")
CDP(async (client) => {
  console.log("Connected!");

  const {Network, Page, Runtime} = client;
  config.client = client;

  try {
    await Network.enable();
    await Page.enable();

    // Capture the browser logs on our console
    await Runtime.enable()
    await Runtime.consoleAPICalled((p) => {
      console.log(p);
    });

    console.log(' Completed');
  } catch (err) {
    console.error(err);
  } finally {
    client.close();
  }
}).on('error', (err) => {
  console.error(err);
});
