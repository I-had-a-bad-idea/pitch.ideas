import { loggedIn } from "./cookie.js";

const orderBy = document.getElementById("order_by");
orderBy.addEventListener("change", () => {
    loadPitches();
})

async function loadPitches() {
    const order_by = orderBy.value;

    const res = await fetch(`/pitches?order_by=${order_by}`, {credentials: "include",});
    const data = await res.json();

    const container = document.querySelector(".feed .container");

    // Remove all pitches
    container.querySelectorAll(".pitch").forEach(p => p.remove());

    data.pitches.forEach(p => {
        const maxLength = 300;
        const description =
        p.description.length > maxLength
            ? p.description.slice(0, maxLength) + "..."
            : p.description;

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

        <p class="description">${description}</p>
        <div class="pitch-footer">
            <span class="vote-btn ${p.voted_by_user ? "" : "voted"}">
                👍 <span class="vote-count">${p.votes}</span>
            </span>
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
                
                const data = await response.json().catch(() => ({}));
                
                if (response.ok) {
                    // toggle UI state
                    voteBtn.classList.toggle("voted");
                    console.log(data.votes);
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
    });


}

const nav_right = document.querySelector("nav .nav-right");

async function showLoggedInButtons() {
    const createPitchlink = document.createElement("a");
    createPitchlink.href = "/pitches/create";
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
    if (loggedIn()) {
        showLoggedInButtons();
    } else {
        showAuthButtons();
    }
}

checkAuth();
loadPitches();