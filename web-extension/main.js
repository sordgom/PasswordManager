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

    // I'm Finding all elements that could potentially be the submit button
    // If there's more than one button, it might require specific logic to choose the right one
    const submitButtons = document.querySelectorAll('button[type="submit"], input[type="submit"], button[id="submit"]');
    for (const button of submitButtons) {
        if (isVisible(button)) {
            // Will uncomment this once the logic is sound
            // button.click();
            console.log("Login form submitted.");
            return;
        }
    }
    console.log("No visible submit button found.");

}

// Fetch credentials and autofill
function autofillCredentials(usernameInput, passwordInput) {
    getPassword().then(data => {
        usernameInput.value = data.name;  
        passwordInput.value = data.url;
    });
}

function getPassword() {
    return new Promise((resolve, reject) => {
        // Hardcode vault name and password name for now
        const vault_name = "test";
        const password_name = "test";
        const host = "http://localhost:8080";
        const url = `${host}/password?vault_name=${vault_name}&password_name=${password_name}`;

        fetch(url)
            .then(response => response.json())
            .then(data => resolve(data))
            .catch(error => reject(error));
    });
}

AutoLogin();