const handleChangeTheme = (isChecked) => {
    if (isChecked.checked) {
        document.body.setAttribute('light', '');
    } else {
        document.body.removeAttribute('light');
    }
}