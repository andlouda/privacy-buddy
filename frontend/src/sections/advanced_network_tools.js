import { GetNetworkInterfaces } from '../../wailsjs/go/tools/NetworkToolsService';
import { GetPublicIPInfo } from '../../wailsjs/go/network/PublicIPService';
import {
  StartPacketCapture,
  StopPacketCapture,
  GetCaptureTemplates,
  SaveCaptureTemplate
} from '../../wailsjs/go/tools/AdvancedNetworkToolsService';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';

let eventListenerInitialized = false;

function formatMacAddress(mac) {
  if (!mac) return '-';
  // Assuming mac is a base64 encoded string of bytes, decode it first
  try {
    const byteCharacters = atob(mac);
    const byteNumbers = new Array(byteCharacters.length);
    for (let i = 0; i < byteCharacters.length; i++) {
      byteNumbers[i] = byteCharacters.charCodeAt(i);
    }
    const byteArray = new Uint8Array(byteNumbers);
    return Array.from(byteArray).map(b => b.toString(16).padStart(2, '0')).join(':');
  } catch (e) {
    console.warn("Failed to decode MAC address:", mac, e);
    return mac; // Return original if decoding fails
  }
}

function generateInterfaceTable(interfaces, externalInterfaceName) {
  console.debug('[generateInterfaceTable] externalInterfaceName:', externalInterfaceName);

  if (!interfaces || interfaces.length === 0) return '<p>No network interfaces found.</p>';

  const headers = ['Name', 'Description', 'MAC Address', 'IP Addresses', 'Status'];
  let table = '<table><thead><tr>' + headers.map(h => `<th>${h}</th>`).join('') + '</tr></thead><tbody>';

  interfaces.forEach(iface => {
    const name = iface.Name || '';
    const description = iface.Description || '';
    const isUplink = externalInterfaceName &&
      (name.toLowerCase().includes(externalInterfaceName.toLowerCase()) ||
       description.toLowerCase().includes(externalInterfaceName.toLowerCase()));
    const isUp = iface.IsUp;
    const statusText = isUplink ? '⚡ Uplink' : (isUp ? '🟢 Up' : '🔴 Down');
    const rowClass = isUplink ? 'uplink-interface' : (isUp ? 'up-interface' : 'down-interface');

    table += `<tr class="${rowClass}"><td>${name}</td><td>${description}</td><td>${formatMacAddress(iface.HardwareAddr) || '-'}</td><td>${(iface.Addrs || []).join('<br>') || '-'}</td><td>${statusText}</td></tr>`;
  });

  return table + '</tbody></table>';
}

function populateInterfaceDropdown(selectElement, interfaces) {
  console.debug('[populateInterfaceDropdown] interfaces:', interfaces);
  selectElement.innerHTML = '';

  const defaultOption = document.createElement('option');
  defaultOption.value = '';
  defaultOption.textContent = interfaces.length ? 'Select an interface' : 'No interfaces available';
  if (!interfaces.length) defaultOption.disabled = true;
  selectElement.appendChild(defaultOption);

  interfaces.forEach(iface => {
    const option = document.createElement('option');
    option.value = iface.Name;
    option.textContent = iface.DisplayName || iface.Name;
    selectElement.appendChild(option);
  });
}

function populateTemplateDropdown(selectElement, templates, bpfFilterInput, captureDurationInput) {
  console.debug('[populateTemplateDropdown] templates:', templates);
  selectElement.innerHTML = '<option value="">Select a template</option>';

  templates.forEach(template => {
    const option = document.createElement('option');
    option.value = template.name;
    option.textContent = `${template.name} (${template.description || ''})`;
    option.dataset.bpfFilter = template.bpfFilter;
    option.dataset.duration = template.duration || '';
    selectElement.appendChild(option);
  });

  selectElement.addEventListener('change', () => {
    const selected = selectElement.options[selectElement.selectedIndex];
    console.debug('[Template Change]', selected);
    bpfFilterInput.value = selected.dataset.bpfFilter || '';
    captureDurationInput.value = selected.dataset.duration || '10';
  });
}

function setupCaptureListeners(outputElement, startBtn, stopBtn) {
  if (eventListenerInitialized) return;
  eventListenerInitialized = true;

  console.debug('[setupCaptureListeners] Registering packetCaptureEvent + packetCaptureStopped');

  EventsOn('packetCaptureEvent', packet => {
    console.debug('[packetCaptureEvent]', packet);
    const line = document.createElement('p');
    line.textContent = `[${packet.Timestamp}] ${packet.Source} -> ${packet.Destination} (${packet.Protocol}) [${packet.Length}] ${packet.Summary}`;
    outputElement.appendChild(line);
    outputElement.scrollTop = outputElement.scrollHeight;
  });

  EventsOn('packetCaptureStopped', msg => {
    console.debug('[packetCaptureStopped] Received:', msg);
    const line = document.createElement('p');
    line.textContent = `🛑 Capture stopped: ${msg}`;
    outputElement.appendChild(line);
    startBtn.disabled = false;
    stopBtn.disabled = true;
    EventsOff('packetCaptureEvent');
    EventsOff('packetCaptureStopped');
    eventListenerInitialized = false;
  });
}

export function initializeAdvancedNetworkTools(sectionElement) {
  console.debug('[init] Initializing advanced network tools');
  if (!sectionElement) return console.error('[init] Missing section element');

  const getBtn = sectionElement.querySelector('#get-interfaces-btn');
  const interfacesOutput = sectionElement.querySelector('#interfaces-output');
  const ifaceSelect = sectionElement.querySelector('#interface-select');
  const bpfInput = sectionElement.querySelector('#bpf-filter');
  const durationInput = sectionElement.querySelector('#capture-duration');
  const startBtn = sectionElement.querySelector('#start-capture-btn');
  const stopBtn = sectionElement.querySelector('#stop-capture-btn');
  const output = sectionElement.querySelector('#packet-capture-output');
  const templateSelect = sectionElement.querySelector('#template-select');
  const tplName = sectionElement.querySelector('#new-template-name');
  const tplDesc = sectionElement.querySelector('#new-template-description');
  const saveBtn = sectionElement.querySelector('#save-template-btn');
  const saveOutput = sectionElement.querySelector('#save-template-output');

  if (!getBtn || !startBtn || !stopBtn || !output || !ifaceSelect || !bpfInput || !durationInput || !templateSelect || !tplName || !tplDesc || !saveBtn || !saveOutput) {
    return console.error('[init] Required DOM elements missing');
  }

  stopBtn.disabled = true;

  getBtn.addEventListener('click', () => {
    console.debug('[getBtn] Clicked');
    interfacesOutput.textContent = 'Loading interfaces...';
    ifaceSelect.disabled = true;

    Promise.all([GetNetworkInterfaces(), GetPublicIPInfo()]).then(([ifaces, ipInfo]) => {
      console.debug('[getBtn] Interfaces:', ifaces);
      console.debug('[getBtn] PublicIP Info:', ipInfo);
      interfacesOutput.innerHTML = generateInterfaceTable(ifaces, ipInfo?.interfaceName);
      populateInterfaceDropdown(ifaceSelect, ifaces);
      ifaceSelect.disabled = false;
    }).catch(err => {
      console.error('[getBtn] Error:', err);
      interfacesOutput.textContent = `Error loading interfaces: ${err}`;
      ifaceSelect.disabled = false;
    });
  });

  startBtn.addEventListener('click', async () => {
    const selected = ifaceSelect.value;
    const bpf = bpfInput.value;
    const dur = parseInt(durationInput.value, 10);

    if (!selected || isNaN(dur)) return;

    output.innerHTML = '⏳ Starting packet capture...';
    startBtn.disabled = true;
    stopBtn.disabled = false;

    setupCaptureListeners(output, startBtn, stopBtn);

    try {
      await StartPacketCapture(selected, bpf, dur);
      console.debug('[startBtn] Capture started successfully');
    } catch (e) {
      console.error('[startBtn] Capture error:', e);
      output.textContent = `❌ Error: ${e}`;
      startBtn.disabled = false;
      stopBtn.disabled = true;
    }
  });

  stopBtn.addEventListener('click', async () => {
    console.debug('[stopBtn] Clicked');
    output.innerHTML += '\n⏹️ Stopping capture...';
    try {
      await StopPacketCapture();
    } catch (e) {
      console.error('[stopBtn] Error:', e);
    }
  });

  saveBtn.addEventListener('click', async () => {
    const name = tplName.value.trim();
    const desc = tplDesc.value.trim();
    const bpf = bpfInput.value.trim();
    const dur = parseInt(durationInput.value, 10);

    if (!name || !bpf) {
      saveOutput.textContent = '❌ Name and BPF Filter required';
      return;
    }

    const tpl = { name, description: desc, bpfFilter: bpf, duration: dur };
    console.debug('[saveBtn] Saving:', tpl);

    try {
      await SaveCaptureTemplate(tpl);
      saveOutput.textContent = `✅ Saved template '${name}'`;
      tplName.value = '';
      tplDesc.value = '';
      const updated = await GetCaptureTemplates();
      populateTemplateDropdown(templateSelect, updated, bpfInput, durationInput);
    } catch (e) {
      console.error('[saveBtn] Error:', e);
      saveOutput.textContent = `❌ Save failed: ${e}`;
    }
  });

  GetCaptureTemplates().then(tpls => {
    console.debug('[init] Templates loaded:', tpls);
    populateTemplateDropdown(templateSelect, tpls, bpfInput, durationInput);
  }).catch(err => console.error('[init] Template fetch error:', err));
}