import { chromium } from 'playwright';

const pages = [
  { url: 'https://nieboguard.pl/', name: 'main' },
  { url: 'https://nieboguard.pl/ar/', name: 'ar' },
  { url: 'https://nieboguard.pl/bezpieczne-miejsca/', name: 'safety' },
  { url: 'https://nieboguard.pl/obserwacje-lotnicze/', name: 'explorer' },
];

(async () => {
  const browser = await chromium.launch();
  for (const p of pages) {
    const page = await browser.newPage({ viewport: { width: 1280, height: 720 } });
    await page.goto(p.url, { waitUntil: 'networkidle', timeout: 30000 });
    await page.screenshot({ path: `docs/screenshots/${p.name}.png`, fullPage: false });
    await page.close();
    console.log('Screenshot:', p.name);
  }
  await browser.close();
})();
