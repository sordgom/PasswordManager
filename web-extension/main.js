function isVisible(elem) {
    return !!(elem.offsetWidth || elem.offsetHeight || elem.getClientRects().length);
}

function autoLogin() {
    const inputTypes = ['text', 'email', 'tel'];
    const passwordInputs = Array.from(document.querySelectorAll('input[type="password"]'));

    passwordInputs.forEach(passwordInput => {
        const textInputs = Array.from(document.querySelectorAll('input')).filter(input => {
            return inputTypes.includes(input.type) && isVisible(input) && input.value.trim() === '';
        });

        const likelyLogin = textInputs.some(input => /user|login|email/i.test(input.name || input.id));
        const likelyNotSignup = !passwordInputs.some(input => input.innerHTML.toLowerCase().includes('confirm password'));

        if (textInputs.length === 1 && likelyLogin && likelyNotSignup) {
            autofillCredentials(textInputs[0], passwordInput);
        }
    });

    const submitButtons = document.querySelectorAll('button[type="submit"], input[type="submit"], button[id="submit"]');
    for (const button of submitButtons) {
        if (isVisible(button)) {
            // Communicate with background script to perform actions requiring background access
            chrome.runtime.sendMessage({ action: 'submitLoginForm' });
            console.log("Login form submitted.");
            return;
        }
    }
    console.log("No visible submit button found.");
}

function autofillCredentials(usernameInput, passwordInput) {
    chrome.runtime.sendMessage({ action: 'getCredentials', url: window.location.href }, credentials => {
        if (credentials) {
            usernameInput.value = credentials.name;  
            passwordInput.value = credentials.url;
        } else {
            console.error('Failed to retrieve credentials');
        }
    });
}

autoLogin();
