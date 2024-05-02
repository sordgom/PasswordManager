// Listen for a message from the content script
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    if (message.action === 'getCredentials') {
        getPassword(message.url)
            .then(credentials => {
                sendResponse(credentials);
            })
            .catch(error => {
                console.error('Error getting credentials:', error);
                sendResponse({ error: 'Failed to get credentials' });
            });
        return true;
    }
});

async function getPassword(url) {
    const vault_name = "test";
    const host = "http://localhost:8080";
    const hostUrl = `${host}/password/url?vault_name=${vault_name}&url=${url}`;

    const response = await fetch(hostUrl);
    if (!response.ok) {
        throw new Error('Failed to fetch credentials');
    }
    return await response.json();
}
