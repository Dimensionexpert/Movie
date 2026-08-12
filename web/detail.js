const TMDB_IMG_BASE = 'https://image.tmdb.org/t/p/w780';

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
  div.textContent = str ?? '';
  return div.innerHTML;
}

function getMovieId() {
  const params = new URLSearchParams(window.location.search);
  return params.get('id');
}

function toggleVideoPlayback() {
  const video = document.querySelector('#screening video');
  if (!video) return;

  if (video.paused) {
    video.play();
  } else {
    video.pause();
  }
}

function render(movie) {
  const root = document.getElementById('detail-root');
  const src = posterSrc(movie.posterUrl);
  const duration = formatDuration(movie.durationSeconds);

  const posterMarkup = src
    ? `<img src="${src}" alt="${escapeHtml(movie.title)} poster">`
    : `<div class="no-poster">${escapeHtml(movie.title)}</div>`;

  const metaParts = [];
  if (movie.releaseYear) metaParts.push(`<span>${movie.releaseYear}</span>`);
  if (duration) metaParts.push(`<span>${duration}</span>`);
  if (metaParts.length === 0) metaParts.push('<span>In your library</span>');

  root.innerHTML = `
    <article class="movie-hero">
      <div class="poster-col">
        ${posterMarkup}
        <span class="poster-badge">Ready to play</span>
      </div>
      <div class="ticket-body">
        <p class="eyebrow">From your collection</p>
        <h1 class="movie-title">${escapeHtml(movie.title)}</h1>
        <div class="meta-row">${metaParts.join('')}</div>
        <div class="overview-panel">
          <span class="overview-label">Synopsis</span>
          <p class="overview">${escapeHtml(movie.overview || 'No synopsis on file.')}</p>
        </div>
        <div class="detail-actions">
          <button class="play-btn" id="play-btn"><span aria-hidden="true">▶</span> Play film</button>
          <span class="play-hint">Space bar toggles playback</span>
        </div>
      </div>
    </article>
    <section class="screening" id="screening" aria-live="polite"></section>
  `;

  const playBtn = document.getElementById('play-btn');
  const screening = document.getElementById('screening');

  if (playBtn && screening) {
    playBtn.addEventListener('click', () => {
      if (screening.classList.contains('is-active')) return;

      screening.innerHTML = `
        <div class="screening-head">
          <p class="eyebrow">Now playing</p>
          <span>${escapeHtml(movie.title)}</span>
        </div>
        <video controls autoplay src="/movies/${movie.id}/stream"></video>
      `;
      screening.classList.add('is-active');
      screening.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    });
  }

  document.addEventListener('keydown', (event) => {
    const tag = document.activeElement && document.activeElement.tagName;
    const isTypingField = tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT';

    if (event.code === 'Space' && !isTypingField) {
      event.preventDefault();
      toggleVideoPlayback();
    }
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
