const copiedFor = 1600;

for (const button of document.querySelectorAll('[data-copies]')) {
  const resting = button.textContent;
  const source = document.getElementById(button.dataset.copies);
  let restore;

  button.addEventListener('click', async () => {
    try {
      await navigator.clipboard.writeText(source.textContent.trim());
      button.textContent = 'Copied';
    } catch {
      button.textContent = 'Select and copy it';
    }
    clearTimeout(restore);
    restore = setTimeout(() => {
      button.textContent = resting;
    }, copiedFor);
  });
}
