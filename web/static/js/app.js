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
    await fetch('/api/lang', {method: 'POST', headers: {'accept': 'application/json'}, body: form});
    location.reload();
  });
}

document.addEventListener('error', (event) => {
  const img = event.target;
  if (!(img instanceof HTMLImageElement)) return;
  if (!img.closest('.work-card, .pk-card, .leader-row, .social-card')) return;
  img.classList.add('image-fallback');
  img.removeAttribute('src');
}, true);

let pkVotes = 0;
let pkLocked = false;
let nextPair = null;
let pkFinished = false;

function renderPKCard(button, work, side) {
  button.dataset.id = work.id;
  button.classList.remove('pk-winner', 'pk-loser', 'pk-loading');
  button.setAttribute('aria-disabled', 'false');
  const visitURL = work.final_url || work.url || '';
  const openLabel = $('[data-open-link]')?.dataset.openLink || 'Open';
  button.innerHTML = `<img src="${escapeAttr(work.image)}" alt=""><div><small>${escapeHTML(work.domain)}</small><h3>${escapeHTML(work.title)}</h3><p>${escapeHTML(work.description || '')}</p><a class="pk-visit-link" href="${escapeAttr(visitURL)}" target="_blank" rel="noopener noreferrer" title="${escapeAttr(openLabel)} ${escapeAttr(visitURL)}">${escapeHTML(work.domain || visitURL)}</a></div>`;
  button.onclick = () => vote(side);
  const link = $('.pk-visit-link', button);
  if (link) link.onclick = event => event.stopPropagation();
  button.onkeydown = event => {
    if (event.target.closest('.pk-visit-link')) return;
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      vote(side);
    }
  };
}

function pkSection() {
  return $('[data-pk]');
}

function setPKStatus(message, mode = 'notice') {
  const section = pkSection();
  const status = $('[data-pk-status]');
  const text = $('[data-pk-message]');
  if (!section || !status || !text) return;
  section.dataset.pkMode = mode;
  text.textContent = message || '';
  status.hidden = !message;
}

function clearPKStatus() {
  setPKStatus('');
}

function setPKDisabled(disabled) {
  const left = $('[data-pk-left]');
  const right = $('[data-pk-right]');
  for (const card of [left, right]) {
    if (!card) continue;
    card.setAttribute('aria-disabled', disabled ? 'true' : 'false');
  }
}

async function fetchPair() {
  return fetch('/api/pk/pair').then(r => r.json()).catch(() => null);
}

function showPair(pair) {
  const left = $('[data-pk-left]');
  const right = $('[data-pk-right]');
  if (!left || !right) return;
  if (!pair || pair.error) {
    const section = pkSection();
    setPKStatus(section?.dataset.pkError || 'Not enough works yet.', 'paused');
    return;
  }
  renderPKCard(left, pair.left, 'left');
  renderPKCard(right, pair.right, 'right');
  setPKDisabled(false);
}

async function loadPair() {
  showPair(nextPair || await fetchPair());
  nextPair = null;
}

async function vote(side) {
  if (pkLocked || pkFinished) return;
  const left = $('[data-pk-left]');
  const right = $('[data-pk-right]');
  const section = pkSection();
  if (!left || !right || !left.dataset.id || !right.dataset.id) return;
  const winner = side === 'left' ? left.dataset.id : right.dataset.id;
  const loser = side === 'left' ? right.dataset.id : left.dataset.id;
  const winnerCard = side === 'left' ? left : right;
  const loserCard = side === 'left' ? right : left;

  pkLocked = true;
  clearPKStatus();
  setPKDisabled(true);
  winnerCard.classList.add('pk-winner');
  loserCard.classList.add('pk-loser');

  const pairRequest = fetch('/api/pk/vote', {
    method: 'POST',
    headers: {'content-type': 'application/json'},
    body: JSON.stringify({winner, loser})
  }).then(r => r.json()).catch(() => null);

  setTimeout(async () => {
    pkVotes++;
    nextPair = await pairRequest;
    winnerCard.classList.remove('pk-winner');
    loserCard.classList.remove('pk-loser');

    if (!nextPair || nextPair.error) {
      pkLocked = false;
      setPKStatus(section?.dataset.pkError || 'Not enough works yet.', 'paused');
      return;
    }

    if (pkVotes >= 30) {
      pkFinished = true;
      pkLocked = false;
      setPKDisabled(true);
      setPKStatus(section?.dataset.pkFinished || 'Thanks for playing.', 'finished');
      return;
    }

    if (pkVotes === 10) {
      pkLocked = false;
      setPKStatus(section?.dataset.pkTen || 'Keep going or view the leaderboard.', 'paused');
      return;
    }

    if (pkVotes === 5) {
      setPKStatus(section?.dataset.pkEncourage || 'Nice picks. Keep going!', 'notice');
    }

    showPair(nextPair);
    nextPair = null;
    pkLocked = false;
  }, 400);
}

const pkContinue = $('[data-pk-continue]');
if (pkContinue) {
  pkContinue.addEventListener('click', () => {
    if (pkFinished) return;
    clearPKStatus();
    loadPair();
    pkLocked = false;
  });
}

const pkLeaderboard = $('[data-pk-leaderboard]');
if (pkLeaderboard) {
  pkLeaderboard.addEventListener('click', () => {
    const target = $('#leaderboard');
    if (target) target.scrollIntoView({behavior: 'smooth', block: 'start'});
  });
}

document.addEventListener('click', event => {
  const card = event.target.closest('.work-card');
  if (!card || event.target.closest('a')) return;
  const section = pkSection();
  if (section) {
    section.scrollIntoView({behavior: 'smooth', block: 'start'});
  }
});
loadPair();

function escapeHTML(value) {
  return String(value).replace(/[&<>"']/g, m => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[m]));
}

function escapeAttr(value) {
  return escapeHTML(value).replace(/`/g, '&#96;');
}
