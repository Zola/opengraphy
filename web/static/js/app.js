const $ = (q, root = document) => root.querySelector(q);

async function refreshStats() {
  const box = $('[data-stats]');
  if (!box) return;
  try {
    await fetch('/api/presence', {method: 'POST'});
    const stats = await fetch('/api/stats').then(r => r.json());
    $('[data-online]').textContent = stats.online ?? 0;
    $('[data-visits]').textContent = stats.visits_total ?? 0;
    $('[data-checks]').textContent = stats.checks_total ?? 0;
  } catch (_) {}
}
refreshStats();
setInterval(refreshStats, 30000);

document.addEventListener('click', async (event) => {
  const copy = event.target.closest('[data-copy]');
  if (copy) {
    const target = $(copy.dataset.copy);
    if (target) await navigator.clipboard.writeText(target.innerText);
  }
});

function getCookie(name) {
  return document.cookie.split('; ').find(row => row.startsWith(`${name}=`))?.split('=')[1] || '';
}

function setCookie(name, value, maxAge) {
  document.cookie = `${name}=${encodeURIComponent(value)}; Path=/; Max-Age=${maxAge}; SameSite=Lax`;
}

const consent = $('[data-consent]');
if (consent && !getCookie('og_consent')) consent.hidden = false;
document.addEventListener('click', (event) => {
  if (!consent) return;
  if (event.target.closest('[data-consent-accept]')) {
    setCookie('og_consent', 'accepted', 31536000);
    consent.hidden = true;
  }
  if (event.target.closest('[data-consent-reject]')) {
    setCookie('og_consent', 'necessary', 31536000);
    consent.hidden = true;
  }
});

const langForm = $('[data-lang-form]');
if (langForm) {
  langForm.addEventListener('change', async () => {
    const form = new FormData(langForm);
    await fetch('/api/lang', {method: 'POST', body: form});
    location.reload();
  });
}

function renderPKCard(button, work, side) {
  button.dataset.id = work.id;
  button.innerHTML = `<img src="${escapeAttr(work.image)}" alt=""><div><small>${escapeHTML(work.domain)}</small><h3>${escapeHTML(work.title)}</h3><p>${escapeHTML(work.description || '')}</p><span>${work.rating} pts</span></div>`;
  button.onclick = () => vote(side);
}

async function loadPair() {
  const left = $('[data-pk-left]');
  const right = $('[data-pk-right]');
  if (!left || !right) return;
  const pair = await fetch('/api/pk/pair').then(r => r.json()).catch(() => null);
  if (!pair || pair.error) return;
  renderPKCard(left, pair.left, 'left');
  renderPKCard(right, pair.right, 'right');
}

async function vote(side) {
  const left = $('[data-pk-left]');
  const right = $('[data-pk-right]');
  const winner = side === 'left' ? left.dataset.id : right.dataset.id;
  const loser = side === 'left' ? right.dataset.id : left.dataset.id;
  const pair = await fetch('/api/pk/vote', {
    method: 'POST',
    headers: {'content-type': 'application/json'},
    body: JSON.stringify({winner, loser})
  }).then(r => r.json()).catch(() => null);
  if (pair && !pair.error) {
    renderPKCard(left, pair.left, 'left');
    renderPKCard(right, pair.right, 'right');
  }
}
loadPair();

function escapeHTML(value) {
  return String(value).replace(/[&<>"']/g, m => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[m]));
}

function escapeAttr(value) {
  return escapeHTML(value).replace(/`/g, '&#96;');
}
