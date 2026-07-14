const logoutButton = document.getElementById('logout-button');

logoutButton.addEventListener('click', async () => {
    try {
        const response = await fetch('/auth/logout', {
            method: 'POST',
            credentials: 'include',
        });
    
        if (response.ok) {
            window.location.href = '/';
        } else {
            alert('Logout failed.');
        }
    } catch (err) {
        console.error(err);
        alert('Unable to connect to the server.');
    }
});