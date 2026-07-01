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
            body: JSON.stringify(body)
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

const voreBtn = document.querySelector(".vote-btn");

voreBtn.addEventListener("click", async (e) => {
    e.stopPropagation(); // Prevent the click from propagating to the parent div

    const pitchId = window.location.pathname.split("/").pop();

    try {
        const response = await fetch(`/pitches/${pitchId}/upvote`, {
            method: "POST"
        });

        if (response.ok) {
            const currentVotes = parseInt(voreBtn.textContent.split(" ")[1]);
            voreBtn.textContent = `👍 ${currentVotes + 1}`;
        } else {
            alert("Failed to upvote pitch.");
        }
    } catch (err) {
        console.error(err);
        alert("Unable to connect to the server.");
    }
});

const nav_right = document.querySelector("nav .nav-right");

async function showCreatePitchButton() {
    const a = document.createElement("a");
    a.href = "/create-pitch";

    const button = document.createElement("button");
    button.className = "btn btn-primary";
    button.textContent = "Create Pitch";

    a.appendChild(button);
    nav_right.appendChild(a);
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
    const res = await fetch("/auth/status");
    const data = await res.json();

    if (data.logged_in) {
        showCreatePitchButton();
    } else {
        showAuthButtons();
    }
}

checkAuth();