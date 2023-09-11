document.addEventListener('DOMContentLoaded', () => {
    const todoItems = document.querySelectorAll('.todo');

    todoItems.forEach((todoItem) => {
        const checkbox = todoItem.querySelector('.todo-checkbox');
        const todoText = todoItem.querySelector('.todo-text');

        // Добавляем обработчик события для чекбокса
        checkbox.addEventListener('change', () => {
            if (checkbox.checked) {
                todoItem.classList.add('checked', 'fade-out'); // Добавляем классы для зачеркивания и затемнения
            } else {
                todoItem.classList.remove('checked', 'fade-out'); // Удаляем классы
            }
        });
    });
});
