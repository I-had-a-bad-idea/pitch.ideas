import Toast from "./toast.js";

const logoutButton = document.getElementById('logout-button');

logoutButton.addEventListener('click', async () => {
    try {
        const response = await fetch('/auth/logout', {
            method: 'POST',
            credentials: 'include',
        });
    
        if (response.ok) {
            Toast.success("Logged out!");
            await new Promise(resolve => setTimeout(resolve, 250));
            window.location.href = '/';
        } else {
            Toast.error('Logout failed.');
        }
    } catch (err) {
        console.error(err);
        Toast.error('Unable to connect to the server.');
    }
});