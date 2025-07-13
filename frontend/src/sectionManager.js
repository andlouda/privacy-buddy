export function showSection(id) {
  document.querySelectorAll('.section').forEach(sec => sec.classList.remove('active'));
  const el = document.getElementById(id);
  if (el) el.classList.add('active');
}
