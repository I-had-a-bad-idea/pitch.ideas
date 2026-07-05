import { loggedIn, getUsername } from "./cookie.js";


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

const form = document.querySelector(".comment-form");

async function showComments() {
    form.style.visibility = "visible";

    // For each comment check if username is the same as the logged in user, if so show edit and delete buttons
    const username = getUsername();
    const comments = document.querySelectorAll(".comment");
    comments.forEach(comment => {
        const commentMeta = comment.querySelector(".comment-meta");
        if (commentMeta && commentMeta.textContent.includes(`By ${username}`)) {
            const editButton = comment.querySelector(".edit-comment-btn");
            const deleteButton = comment.querySelector(".delete-comment-btn");
            editButton.style.visibility = "visible";
            deleteButton.style.visibility = "visible";
        }
    });
}

async function showPitchControls() {
    const username = getUsername();
    const pitchMeta = document.querySelector(".pitch-meta")
    const pitchFooter = document.querySelector(".pitch-footer");
    if (pitchFooter && pitchMeta && pitchMeta.textContent.includes(`Posted by ${username}`)) {
        const editPitchButton = pitchFooter.querySelector(".edit-pitch-btn");
        const deletePitchButton = pitchFooter.querySelector(".delete-pitch-btn");
        editPitchButton.style.visibility = "visible";
        deletePitchButton.style.visibility = "visible";
    }
}

async function checkAuth() {
    if (loggedIn()) {
        showLoggedInButtons();
        showComments();
        showPitchControls();
    } else {
        showAuthButtons();
    }
}

checkAuth();


form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const body = {
        content: document.getElementById("content").value.trim()
    };
    const pitchId = window.location.pathname.split("/").pop();

    try {
        const response = await fetch(`/pitches/${pitchId}/comments/add`, {
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

const editCommentButtons = document.querySelectorAll(".edit-comment-btn");
editCommentButtons.forEach(button => {
    button.addEventListener("click", async (e) => {
        const commentId = e.target.dataset.commentId;
        const newContent = prompt("Edit your comment:");

        if (newContent !== null && newContent.trim() !== "") {
            try {
                const response = await fetch(`/pitches/${window.location.pathname.split("/").pop()}/comments/${commentId}/edit`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({ content: newContent.trim() }),
                    credentials: "include",
                });

                const data = await response.json().catch(() => ({}));

                if (response.ok) {
                    window.location.reload();
                } else {
                    alert(data.message || "Failed to edit comment.");
                }
            } catch (err) {
                console.error(err);
                alert("Unable to connect to the server.");
            }
        }
    });
});

const deleteCommentButtons = document.querySelectorAll(".delete-comment-btn");
deleteCommentButtons.forEach(button => {
    button.addEventListener("click", async (e) => {
        const commentId = e.target.dataset.commentId;

        if (confirm("Are you sure you want to delete this comment?")) {
            try {
                const response = await fetch(`/pitches/${window.location.pathname.split("/").pop()}/comments/${commentId}/delete`, {
                    method: "DELETE",
                    credentials: "include",
                });

                const data = await response.json().catch(() => ({}));

                if (response.ok) {
                    window.location.reload();
                } else {
                    alert(data.message || "Failed to delete comment.");
                }
            } catch (err) {
                console.error(err);
                alert("Unable to connect to the server.");
            }
        }
    });
});