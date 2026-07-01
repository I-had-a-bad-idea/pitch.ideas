const form = document.querySelector(".comment-form");

form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const body = {
        content: document.getElementById("content").value.trim()
    };
    const pitchId = window.location.pathname.split("/").pop();

    try {
        const response = await fetch(`/pitches/${pitchId}/comment`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(body),
            credentials: "include",
        });
    
    const data = await response.json().catch(() => ({}));

    if (response.ok) {
        window.location.reload();
    } else {
        alert(data.message || "Failed to add comment.");
    }
    } catch (err) {
        console.error(err);
        alert("Unable to connect to the server.");
    }
});

const voteBtn = document.querySelector(".vote-btn");

voteBtn.addEventListener("click", async (e) => {
    e.stopPropagation(); // Prevent the click from propagating to the parent div

    const pitchId = window.location.pathname.split("/").pop();

    try {
        const response = await fetch(`/pitches/${pitchId}/upvote`, {
            method: "POST",
            credentials: "include",
        });
        
        const data = await response.json().catch(() => ({}));

        if (response.ok) {
            // toggle UI state
            voteBtn.classList.toggle("voted");

            // update count to what backend returns
            voteBtn.querySelector(".vote-count").textContent = data.votes;
        } else {
            alert("Failed to upvote pitch.");
        }
    } catch (err) {
        console.error(err);
        alert("Unable to connect to the server.");
    }
});

const nav_right = document.querySelector("nav .nav-right");

async function showLoggedInButtons() {
    const createPitchlink = document.createElement("a");
    createPitchlink.href = "/create-pitch";
    createPitchlink.className = "btn btn-primary";
    createPitchlink.textContent = "Create Pitch";
    
    const logoutLink = document.createElement("a");
    logoutLink.href = "/auth/logout";
    logoutLink.className = "btn btn-secondary";
    logoutLink.textContent = "Logout";

    nav_right.appendChild(createPitchlink);
    nav_right.appendChild(logoutLink);
}

async function showAuthButtons() {
    const loginLink = document.createElement("a");
    loginLink.href = "/auth/login";
    loginLink.className = "btn btn-primary";
    loginLink.textContent = "Login";

    const registerLink = document.createElement("a");
    registerLink.href = "/auth/register";
    registerLink.className = "btn btn-primary";
    registerLink.textContent = "Register";

    nav_right.appendChild(loginLink);
    nav_right.appendChild(registerLink);
}

async function checkAuth() {
    const res = await fetch("/auth/status", {credentials: "include",});
    const data = await res.json();

    if (data.logged_in) {
        showLoggedInButtons();
    } else {
        showAuthButtons();
    }
}

checkAuth();