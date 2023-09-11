document.addEventListener('DOMContentLoaded', () => {
    const accessSettingsButton = document.getElementById('access-settings-button');
    const accessSettingsModal = document.getElementById('access-settings-modal');
    const closeModalButton = document.getElementById('close-modal-button');
    const overlay = document.querySelector('.overlay'); // Выбираем подложку

    // Открыть модальное окно
    accessSettingsButton.addEventListener('click', () => {
        accessSettingsModal.style.display = 'block';
        overlay.style.display = 'block'; // Показываем подложку
    });

    // Закрыть модальное окно
    closeModalButton.addEventListener('click', () => {
        accessSettingsModal.style.display = 'none';
        overlay.style.display = 'none'; // Скрываем подложку
    });
});