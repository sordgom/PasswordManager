function isVisible(elem) {
    return !!(elem.offsetWidth || elem.offsetHeight || elem.getClientRects().length);
}

function AutoLogin() {
    const inputTypes = ['text', 'email', 'tel'];
    const passwordInputs = Array.from(document.querySelectorAll('input[type="password"]'));

    passwordInputs.forEach(passwordInput => {
        const textInputs = Array.from(document.querySelectorAll('input')).filter(input => {
            return inputTypes.includes(input.type) && isVisible(input) && input.value.trim() === '';
        });

        // Check for common login identifiers and absence of signup identifiers
        const likelyLogin = textInputs.some(input => /user|login|email/i.test(input.name || input.id));
        const likelyNotSignup = !passwordInputs.some(input => input.innerHTML.toLowerCase().includes('confirm password'));

        if (textInputs.length === 1 && likelyLogin && likelyNotSignup) {
            autofillCredentials(textInputs[0], passwordInput);
        }
    });

    // Find all elements that could potentially be the submit button
    const submitButtons = document.querySelectorAll('button[type="submit"], input[type="submit"], button[id="submit"]');
    // If there's more than one button, it might require specific logic to choose the right one
    for (const button of submitButtons) {
        if (isVisible(button)) {
            button.click();
            console.log("Login form submitted.");
            return;
        }
    }
    console.log("No visible submit button found.");

}

// Fetch credentials and autofill
function autofillCredentials(usernameInput, passwordInput) {
    usernameInput.value = 'student';  
    passwordInput.value = 'Password123';
}

AutoLogin()
