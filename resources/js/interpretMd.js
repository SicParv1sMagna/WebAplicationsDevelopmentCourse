const interpretMd = (origin) => {
    const previewElement = document.getElementById('preview');
    previewElement.innerHTML = marked.parse(origin);

    hljs.highlightAll();

    // Прокрутите <div id="preview"> вниз
    previewElement.scrollTop = previewElement.scrollHeight;
};

document.addEventListener('DOMContentLoaded', () => {
    const editorElement = document.getElementById('editor');

    interpretMd(editorElement.value);

    editorElement.addEventListener('input', (e) => {
        interpretMd(e.target.value);
    });

    // Обработчик события keydown для нажатия клавиши Tab
    editorElement.addEventListener('keydown', (e) => {
        if (e.key === 'Tab') {
            e.preventDefault(); // Предотвращаем стандартное поведение (переключение фокуса)
            
            // Получаем текущее значение в <textarea>
            const value = e.target.value;
            
            // Получаем позицию курсора
            const start = e.target.selectionStart;
            const end = e.target.selectionEnd;
            
            // Вставляем символ табуляции (или пробелы) на место курсора
            const newValue = value.substring(0, start) + '\t' + value.substring(end);
            
            // Устанавливаем новое значение в <textarea>
            e.target.value = newValue;
            
            // Устанавливаем курсор после вставленной табуляции
            e.target.selectionStart = e.target.selectionEnd = start + 1;
            
            // Обновляем предварительный просмотр
            interpretMd(newValue);
        }
    });
});
