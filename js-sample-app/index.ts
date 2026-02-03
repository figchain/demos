
import { FigChainClient } from 'figchain-client';

async function main() {
  console.log("Starting FigChain JS Sample App...");
  try {
    const client = await FigChainClient.create('client-config.json');
    
    const figKey = 'test';
    console.log(`Listening for updates on key: ${figKey}`);
    
    // Listener
    client.on(figKey, (val: any) => {
        console.log(`>>> UPDATE RECEIVED for ${figKey} <<<`);
        console.log(`New Configuration:`, JSON.stringify(val, null, 2));
    });
    
    // Initial Get
    const val = await client.getFig(figKey);
    console.log(`Initial value fetched:`, JSON.stringify(val, null, 2));

    const once = process.argv.includes('--once');
    if (once) {
        console.log("Exiting because --once was specified.");
        client.close();
        process.exit(0);
    }
    
    // Keep alive
    process.on('SIGINT', () => {
        console.log("Shutting down...");
        client.close();
        process.exit(0);
    });
    
    // Prevent exit
    setInterval(() => {}, 60000);
  } catch (e) {
      console.error("Error:", e);
      process.exit(1);
  }
}

main();
