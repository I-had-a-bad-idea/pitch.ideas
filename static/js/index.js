async function loadPitches() {
    const res = await fetch("/pitches", {credentials: "include",});
    const data = await res.json();

    const container = document.querySelector(".feed .container");

    data.pitches.forEach(p => {
        const div = document.createElement("div");
        div.className = "pitch";

        div.addEventListener("click", () => {
            window.location.href = `/pitches/${p.id}`;
        });

        div.innerHTML = `
        <div class="pitch-header">
            <div class="pitch-title">${p.title}</div>
            <div class="tag">${p.topic}</div>
        </div>

        <p class="description">${p.description}</p>
        <div class="pitch-footer">
            <span class="vote-btn">👍 ${p.votes}</span>
            <span>💬 ${p.comment_count}</span>
        </div>
        `;

        container.appendChild(div);

        const voteBtn = div.querySelector(".vote-btn");
        voteBtn.addEventListener("click", async (e) => {
            e.stopPropagation(); // Prevent the click from propagating to the parent div

            try {
                const response = await fetch(`/pitches/${p.id}/upvote`, {
                    method: "POST",
                    credentials: "include",
                });

                if (response.ok) {
                    p.votes = Number(p.votes);
                    p.votes += 1; // Increment the vote count locally
                    voteBtn.textContent = `👍 ${p.votes}`; // Update the button text
                } else {
                    alert("Failed to upvote pitch.");
                }
            } catch (err) {
                console.error(err);
                alert("Unable to connect to the server.");
            }
        });
    });


}

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
loadPitches();