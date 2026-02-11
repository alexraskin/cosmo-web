<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed } from "vue";

interface ImageItem {
  key: string;
  size: number;
  uploaded: string;
}

const images = ref<ImageItem[]>([]);
const loading = ref(true);
const error = ref("");
const selectedIndex = ref(-1);
const baseURL = ref("https://cdn.cosmothecat.net/cdn-cgi/image");

const selectedImage = computed(() =>
  selectedIndex.value >= 0 ? images.value[selectedIndex.value] : null,
);

const thumbnailURL = (key: string) =>
  `${baseURL.value}/width=400,height=400,fit=cover,quality=80/${key}`;

const fullURL = (key: string) =>
  `${baseURL.value}/width=1600,quality=85/${key}`;

const loadImages = async () => {
  try {
    const res = await fetch("/api/images");
    if (!res.ok) throw new Error(`Failed to load images: ${res.status}`);
    const data = await res.json();
    images.value = data.images || [];
    if (data.baseURL) baseURL.value = data.baseURL;
  } catch (e: any) {
    error.value = e.message;
    console.error("Error loading images:", e);
  } finally {
    loading.value = false;
  }
};

const openViewer = (index: number) => {
  selectedIndex.value = index;
  document.body.style.overflow = "hidden";
};

const closeViewer = () => {
  selectedIndex.value = -1;
  document.body.style.overflow = "";
};

const prevImage = () => {
  if (selectedIndex.value > 0) {
    selectedIndex.value--;
  } else {
    selectedIndex.value = images.value.length - 1;
  }
};

const nextImage = () => {
  if (selectedIndex.value < images.value.length - 1) {
    selectedIndex.value++;
  } else {
    selectedIndex.value = 0;
  }
};

const handleKeyDown = (e: KeyboardEvent) => {
  if (selectedIndex.value < 0) return;
  if (e.key === "Escape") closeViewer();
  if (e.key === "ArrowLeft") prevImage();
  if (e.key === "ArrowRight") nextImage();
};

onMounted(() => {
  loadImages();
  window.addEventListener("keydown", handleKeyDown);
});

onUnmounted(() => {
  window.removeEventListener("keydown", handleKeyDown);
  document.body.style.overflow = "";
});
</script>

<template>
  <div class="container">
    <!-- Gallery Panel -->
    <section class="panel">
      <div class="panel-header header-accent">
        <span>Cosmo the Cat</span>
        <span>🐱</span>
      </div>
      <div class="panel-body gallery-body">
        <div v-if="loading" class="loading-state">
          <span class="spinner"></span>
          <span>loading photos...</span>
        </div>

        <div v-else-if="error" class="error-state">
          <span>{{ error }}</span>
        </div>

        <div v-else-if="images.length === 0" class="empty-state">
          <span>no photos yet</span>
        </div>

        <div v-else class="gallery-grid">
          <button
            v-for="(image, index) in images"
            :key="image.key"
            class="gallery-item"
            @click="openViewer(index)"
          >
            <img
              :src="thumbnailURL(image.key)"
              :alt="image.key"
              loading="lazy"
              class="gallery-image"
            />
          </button>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="footer">powered by purrs & whiskers</footer>
  </div>

  <!-- Lightbox -->
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="selectedImage" class="lightbox" @click="closeViewer">
        <div class="lightbox-backdrop"></div>

        <button
          class="lightbox-nav lightbox-prev"
          @click.stop="prevImage"
          aria-label="Previous image"
        >
          &larr;
        </button>

        <div class="lightbox-content" @click.stop>
          <img
            :src="fullURL(selectedImage.key)"
            :alt="selectedImage.key"
            class="lightbox-image"
          />
        </div>

        <button
          class="lightbox-nav lightbox-next"
          @click.stop="nextImage"
          aria-label="Next image"
        >
          &rarr;
        </button>

        <button
          class="lightbox-close"
          @click.stop="closeViewer"
          aria-label="Close"
        >
          &times;
        </button>

        <div class="lightbox-counter">
          {{ selectedIndex + 1 }} / {{ images.length }}
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

:root {
  --bg: #111;
  --panel-bg: #f5f5f5;
  --panel-border: 3px solid #111;
  --text-dark: #111;
  --text-light: #fff;
  --text-muted: #555;
  --accent: #8b5cf6;
  --font-heading: "Anton", "Impact", "Arial Black", sans-serif;
  --font-body: "JetBrains Mono", "SF Mono", "Fira Code", "Consolas", monospace;
}

html {
  font-size: 16px;
  -webkit-text-size-adjust: 100%;
}

body {
  font-family: var(--font-body);
  background: var(--bg);
  min-height: 100vh;
  color: var(--text-light);
  line-height: 1.5;
}

/* ── Container ── */

.container {
  max-width: 720px;
  margin: 0 auto;
  padding: 2rem 1rem 3rem;
}

/* ── Panel ── */

.panel {
  background: var(--panel-bg);
  border: var(--panel-border);
  margin-bottom: 1rem;
}

.panel-header {
  background: var(--text-dark);
  color: var(--text-light);
  padding: 0.7rem 1rem;
  font-family: var(--font-heading);
  font-size: 1.15rem;
  font-weight: 400;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: var(--panel-border);
}

.header-accent {
  background: var(--accent);
}

.panel-body {
  padding: 1.25rem;
}

/* ── Gallery ── */

.gallery-body {
  padding: 0.75rem;
}

.gallery-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 0.5rem;
}

.gallery-item {
  aspect-ratio: 1;
  overflow: hidden;
  border: 2px solid var(--text-dark);
  cursor: pointer;
  background: #ddd;
  padding: 0;
  transition: transform 0.12s, box-shadow 0.12s;
}

.gallery-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.gallery-item:active {
  transform: scale(0.98);
}

.gallery-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

/* ── States ── */

.loading-state,
.error-state,
.empty-state {
  text-align: center;
  padding: 2rem 1rem;
  color: var(--text-muted);
  font-size: 0.82rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.error-state {
  color: #e53e3e;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid var(--text-muted);
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* ── Footer ── */

.footer {
  text-align: center;
  margin-top: 1.5rem;
  padding-top: 1.25rem;
  border-top: 3px solid #333;
  font-size: 0.7rem;
  color: #777;
}

.footer a {
  color: var(--accent);
  text-decoration: none;
}

.footer a:hover {
  text-decoration: underline;
}

/* ── Lightbox ── */

.lightbox {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.lightbox-backdrop {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.95);
}

.lightbox-content {
  position: relative;
  z-index: 1;
  max-width: 90vw;
  max-height: 85vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.lightbox-image {
  max-width: 100%;
  max-height: 85vh;
  object-fit: contain;
  border: var(--panel-border);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
}

.lightbox-close {
  position: absolute;
  top: 1rem;
  right: 1rem;
  z-index: 2;
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.08);
  color: var(--text-light);
  border: 2px solid rgba(255, 255, 255, 0.2);
  font-size: 1.5rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-body);
  transition: background 0.15s, border-color 0.15s;
}

.lightbox-close:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.4);
}

.lightbox-nav {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  z-index: 2;
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.08);
  color: var(--text-light);
  border: 2px solid rgba(255, 255, 255, 0.2);
  font-size: 1.25rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-body);
  font-weight: 700;
  transition: background 0.15s, border-color 0.15s;
}

.lightbox-nav:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.4);
}

.lightbox-prev {
  left: 1rem;
}

.lightbox-next {
  right: 1rem;
}

.lightbox-counter {
  position: absolute;
  bottom: 1.5rem;
  left: 50%;
  transform: translateX(-50%);
  z-index: 2;
  font-family: var(--font-body);
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.6);
  background: rgba(0, 0, 0, 0.5);
  padding: 0.25rem 0.75rem;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

/* ── Transitions ── */

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* ── Responsive ── */

@media (min-width: 480px) {
  .container {
    padding: 3rem 1.5rem 4rem;
  }

  .gallery-grid {
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  }
}

@media (max-width: 480px) {
  .gallery-grid {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  }

  .lightbox-nav {
    width: 36px;
    height: 36px;
    font-size: 1rem;
  }

  .lightbox-close {
    width: 36px;
    height: 36px;
    font-size: 1.25rem;
  }
}
</style>
