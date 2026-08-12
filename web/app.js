const TMDB_IMG_BASE = 'https://image.tmdb.org/t/p/w780';

let allMovies = [];

function posterSrc(posterUrl) {
  if (!posterUrl) return null;
  return posterUrl.startsWith('http') ? posterUrl : TMDB_IMG_BASE + posterUrl;
}

function escapeHtml(str) {
  const div = document.createElement('div');
  div.textContent = str ?? '';
  return div.innerHTML;
}

function movieCard(movie) {
  const src = posterSrc(movie.posterUrl);
  const year = movie.releaseYear || null;

  const posterMarkup = src
    ? `<img src="${src}" alt="${escapeHtml(movie.title)} poster" loading="lazy">`
    : `<div class="no-poster">${escapeHtml(movie.title)}</div>`;

  const yearMarkup = year ? `<span class="card-year">${year}</span>` : '';

  const card = document.createElement('a');
  card.className = 'movie-card';
  card.href = `/detail.html?id=${movie.id}`;
  card.innerHTML = `
    <div class="poster-wrap">${posterMarkup}</div>
    <div class="card-label">
      <p class="card-title">${escapeHtml(movie.title)}</p>
      ${yearMarkup}
    </div>
  `;
  return card;
}

function renderFeatured(movie) {
  const featured = document.getElementById('featured');
  if (!featured) return;

  if (!movie) {
    featured.innerHTML = '<div class="featured-empty">No films in the archive yet.</div>';
    return;
  }

  const src = posterSrc(movie.posterUrl);
  const overview = movie.overview || 'No synopsis on file.';
  const runtime = movie.durationSeconds ? formatDuration(movie.durationSeconds) : null;

  featured.innerHTML = `
    <div class="featured-poster">
      ${src ? `<img src="${src}" alt="${escapeHtml(movie.title)} poster">` : `<div class="no-poster">${escapeHtml(movie.title)}</div>`}
    </div>
    <div class="featured-copy">
      <p class="kicker">Featured film</p>
      <h2>${escapeHtml(movie.title)}</h2>
      <div class="featured-meta">
        ${movie.releaseYear ? `<span>${movie.releaseYear}</span>` : ''}
        ${runtime ? `<span>${runtime}</span>` : ''}
      </div>
      <p class="featured-overview">${escapeHtml(overview)}</p>
      <div class="featured-actions">
        <a class="button-link" href="/detail.html?id=${movie.id}">Open details</a>
      </div>
    </div>
  `;
}

function formatDuration(totalSeconds) {
  if (!totalSeconds) return null;
  const h = Math.floor(totalSeconds / 3600);
  const m = Math.floor((totalSeconds % 3600) / 60);
  return h > 0 ? `${h}h ${m}m` : `${m}m`;
}

function renderMovies(list) {
  const grid = document.getElementById('grid');
  const resultsCopy = document.getElementById('results-copy');

  if (!grid) return;

  grid.innerHTML = '';

  if (!list || list.length === 0) {
    grid.innerHTML = '<p class="state-message">Nothing in the archive yet — run a scan to add films.</p>';
    if (resultsCopy) resultsCopy.textContent = 'No results';
    return;
  }

  list.forEach((movie) => grid.appendChild(movieCard(movie)));

  if (resultsCopy) {
    resultsCopy.textContent = `${list.length} ${list.length === 1 ? 'film' : 'films'} in view`;
  }
}

function applySearch(query) {
  const cleanQuery = query.trim().toLowerCase();
  const filtered = !cleanQuery
    ? allMovies
    : allMovies.filter((movie) => movie.title.toLowerCase().includes(cleanQuery));

  renderMovies(filtered);

  if (filtered.length > 0) {
    renderFeatured(filtered[0]);
  } else {
    renderFeatured(null);
  }
}

async function loadMovies() {
  const countEl = document.getElementById('film-count');
  const searchInput = document.getElementById('movie-search');

  try {
    const res = await fetch('/movies');
    if (!res.ok) throw new Error(`status ${res.status}`);
    const movies = await res.json();

    allMovies = Array.isArray(movies) ? movies : [];

    if (allMovies.length === 0) {
      countEl.textContent = '0 films';
      renderFeatured(null);
      renderMovies([]);
      return;
    }

    countEl.textContent = `${allMovies.length} film${allMovies.length === 1 ? '' : 's'}`;
    renderFeatured(allMovies[0]);
    renderMovies(allMovies);

    if (searchInput) {
      searchInput.addEventListener('input', (event) => applySearch(event.target.value));
    }
  } catch (err) {
    const grid = document.getElementById('grid');
    if (grid) {
      grid.innerHTML = '<p class="state-message">Could not reach the archive. Is the server running?</p>';
    }
    if (countEl) countEl.textContent = 'Archive unavailable';
    console.error('failed to load movies:', err);
  }
}

loadMovies();
