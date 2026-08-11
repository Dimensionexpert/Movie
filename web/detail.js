const TMDB_IMG_BASE = 'https://image.tmdb.org/t/p/w500';

function posterSrc(posterUrl) {
  if (!posterUrl) return null;
  return posterUrl.startsWith('http') ? posterUrl : TMDB_IMG_BASE + posterUrl;
}

function formatDuration(totalSeconds) {
  if (!totalSeconds) return null;
  const h = Math.floor(totalSeconds / 3600);
  const m = Math.floor((totalSeconds % 3600) / 60);
  return h > 0 ? `${h}h ${m}m` : `${m}m`;
}

function escapeHtml(str) {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}

function getMovieId() {
  const params = new URLSearchParams(window.location.search);
  return params.get('id');
}

function render(movie) {
  const root = document.getElementById('detail-root');
  const src = posterSrc(movie.posterUrl);
  const duration = formatDuration(movie.durationSeconds);

  const posterMarkup = src
    ? `<img src="${src}" alt="">`
    : `<div class="no-poster">${escapeHtml(movie.title)}</div>`;

  const metaParts = [];
  if (movie.releaseYear) metaParts.push(`<span>${movie.releaseYear}</span>`);
  if (duration) metaParts.push(`<span>${duration}</span>`);

  root.innerHTML = `
    <div class="ticket">
      <div class="poster-col">${posterMarkup}</div>
      <div class="sprockets-vertical"></div>
      <div class="ticket-body">
        <h1 class="movie-title">${escapeHtml(movie.title)}</h1>
        <div class="meta-row">${metaParts.join('')}</div>
        <p class="overview">${escapeHtml(movie.overview || 'No synopsis on file.')}</p>
        <button class="play-btn" id="play-btn">Play</button>
        <div class="screening" id="screening"></div>
      </div>
    </div>
  `;

  document.getElementById('play-btn').addEventListener('click', () => {
    const screening = document.getElementById('screening');
    if (screening.classList.contains('is-active')) return;

    screening.innerHTML = `<video controls autoplay src="/movies/${movie.id}/stream"></video>`;
    screening.classList.add('is-active');
    screening.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
  });
}

async function loadMovie() {
  const root = document.getElementById('detail-root');
  const id = getMovieId();

  if (!id) {
    root.innerHTML = '<p class="detail-error">No film specified.</p>';
    return;
  }

  try {
    const res = await fetch(`/movies/${id}`);
    if (!res.ok) throw new Error(`status ${res.status}`);
    const movie = await res.json();
    document.title = `${movie.title} — movieLibrary`;
    render(movie);
  } catch (err) {
    root.innerHTML = '<p class="detail-error">Could not find this film in the archive.</p>';
    console.error('failed to load movie:', err);
  }
}

loadMovie();
