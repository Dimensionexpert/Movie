const TMDB_IMG_BASE = 'https://image.tmdb.org/t/p/w500';

function posterSrc(posterUrl) {
  if (!posterUrl) return null;
  return posterUrl.startsWith('http') ? posterUrl : TMDB_IMG_BASE + posterUrl;
}

function movieCard(movie) {
  const src = posterSrc(movie.posterUrl);
  const year = movie.releaseYear || null;

  const posterMarkup = src
    ? `<img src="${src}" alt="" loading="lazy">`
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

function escapeHtml(str) {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}

async function loadMovies() {
  const grid = document.getElementById('grid');
  const countEl = document.getElementById('film-count');

  try {
    const res = await fetch('/movies');
    if (!res.ok) throw new Error(`status ${res.status}`);
    const movies = await res.json();

    if (!movies || movies.length === 0) {
      grid.innerHTML = '<p class="state-message">Nothing in the archive yet — run a scan to add films.</p>';
      countEl.textContent = 'ARCHIVE — EMPTY';
      return;
    }

    countEl.textContent = `ARCHIVE — ${movies.length} FILM${movies.length === 1 ? '' : 'S'}`;
    movies.forEach(m => grid.appendChild(movieCard(m)));
  } catch (err) {
    grid.innerHTML = '<p class="state-message">Could not reach the archive. Is the server running?</p>';
    console.error('failed to load movies:', err);
  }
}

loadMovies();
