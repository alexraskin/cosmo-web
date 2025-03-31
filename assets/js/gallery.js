document.addEventListener('DOMContentLoaded', () => {
  const galleryItems = document.querySelectorAll('.gallery-item');
  const modal = document.querySelector('.modal');
  const modalImg = document.querySelector('.modal-content');
  const closeButton = document.querySelector('.close');
  let currentIndex = 0;
  
  galleryItems.forEach((item, index) => {
    item.addEventListener('click', () => {
      currentIndex = index;
      const fullsizeUrl = item.getAttribute('data-fullsize');
      openModal(fullsizeUrl);
    });
  });
  
  if (closeButton) {
    closeButton.addEventListener('click', closeModal);
  }
  
  if (modal) {
    modal.addEventListener('click', (e) => {
      if (e.target === modal) {
        closeModal();
      }
    });
  }
  
  document.addEventListener('keydown', (e) => {
    if (!modal || modal.style.display !== 'flex') return;
    
    if (e.key === 'Escape') {
      closeModal();
    } else if (e.key === 'ArrowRight') {
      showNextImage();
    } else if (e.key === 'ArrowLeft') {
      showPrevImage();
    }
  });
  
  function showNextImage() {
    currentIndex = (currentIndex + 1) % galleryItems.length;
    const nextFullsizeUrl = galleryItems[currentIndex].getAttribute('data-fullsize');
    modalImg.src = nextFullsizeUrl;
  }
  
  function showPrevImage() {
    currentIndex = (currentIndex - 1 + galleryItems.length) % galleryItems.length;
    const prevFullsizeUrl = galleryItems[currentIndex].getAttribute('data-fullsize');
    modalImg.src = prevFullsizeUrl;
  }
  
  function openModal(imageUrl) {
    if (modal && modalImg) {
      modal.style.display = 'flex';
      modalImg.src = imageUrl;
      document.body.style.overflow = 'hidden';
    }
  }
  
  function closeModal() {
    if (modal) {
      modal.style.display = 'none';
      document.body.style.overflow = '';
    }
  }
}); 