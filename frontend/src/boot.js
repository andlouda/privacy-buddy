import { GetLocalIP, GetPublicIP } from "../wailsjs/go/network/NetworkDashboardService";
import { GetSystemInfo } from "../wailsjs/go/system/SystemService";
import { showView } from "./navigation.js";
import { createBootBox, appendBootLine } from "./ui.js";

function delay(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function safeCall(fn, fallback = 'unbekannt') {
  try {
    return await fn();
  } catch (err) {
    console.warn(`Fehler bei ${fn.name}:`, err);
    return fallback;
  }
}

export async function runBootSequence() {
  // 🧠 Initialisierung
  const boot = createBootBox('boot-box', '🧠', 'System wird initialisiert...');
  boot.wrapper.classList.remove('hidden');
  await appendBootLine(boot.content, `Starte Privacy Buddy...`);
  await appendBootLine(boot.content, `Version: 0.0.1`);

  // 🌐 Netzwerk
  const net = createBootBox('net-box', '🌐', 'Initialisiere Netzwerk...');
  await delay(300);
  net.wrapper.classList.remove('hidden');
  const localIP = await safeCall(GetLocalIP);
  const publicIP = await safeCall(GetPublicIP);
  await appendBootLine(net.content, `Lokale IP: ${localIP}`);
  await appendBootLine(net.content, `Öffentliche IP: ${publicIP}`);

  // 💻 Systeminfo
  const sys = createBootBox('sys-box', '💻', 'Hole Systeminformationen...');
  await delay(300);
  sys.wrapper.classList.remove('hidden');
  const system = await safeCall(GetSystemInfo);
  await appendBootLine(sys.content, `Hostname: ${system.hostname}`);
  await appendBootLine(sys.content, `OS: ${system.os} (${system.arch})`);
  await appendBootLine(sys.content, `Benutzer: ${system.username}`);

  // ✅ Final
  const final = createBootBox('final-box', '🖥️', 'Zugriff gewährt');
  await delay(300);
  final.wrapper.classList.remove('hidden');
  await appendBootLine(final.content, `Sende Daten an NSA... ✅`);

  await delay(3000);

  const context = {
    sys: system,
    localIP,
    publicIP
  };
  return context;
}
