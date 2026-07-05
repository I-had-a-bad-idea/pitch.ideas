import { loggedIn, getUsername } from "./cookie.js";

const pitchId = window.location.pathname.split("/").pop();
const nav_right = document.querySelector("nav .nav-right");
const commentForm = document.querySelector(".comment-form");

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

async function showComments() {
    if (commentForm) {
        commentForm.style.visibility = "visible";
    }

    // For each comment check if username is the same as the logged in user, if so show edit and delete buttons
    const username = getUsername();
    const comments = document.querySelectorAll(".comment");
    comments.forEach((comment) => {
        const commentMeta = comment.querySelector(".comment-meta");
        const editButton = comment.querySelector(".edit-comment-btn");
        const deleteButton = comment.querySelector(".delete-comment-btn");
        const isOwner = commentMeta && commentMeta.textContent.includes(`By ${username}`);

        if (editButton) {
            editButton.classList.toggle("is-visible", Boolean(isOwner));
        }
        if (deleteButton) {
            deleteButton.classList.toggle("is-visible", Boolean(isOwner));
        }
    });
}

function showPitchControls() {
    const username = getUsername();
    const pitchMeta = document.querySelector(".pitch-meta");
    const pitchFooter = document.querySelector(".pitch-footer");
    if (pitchFooter && pitchMeta && pitchMeta.textContent.includes(`Posted by ${username}`)) {
        const editPitchButton = pitchFooter.querySelector(".edit-pitch-btn");
        const deletePitchButton = pitchFooter.querySelector(".delete-pitch-btn");
        if (editPitchButton) {
            editPitchButton.classList.add("is-visible");
        }
        if (deletePitchButton) {
            deletePitchButton.classList.add("is-visible");
        }
    }
}

function checkAuth() {
    if (loggedIn()) {
        showLoggedInButtons();
        showComments();
        showPitchControls();
    } else {
        showAuthButtons();
    }
}

function setupPitchEditor() {
    const editPitchButton = document.querySelector(".edit-pitch-btn");
    const savePitchButton = document.querySelector(".save-pitch-btn");
    const cancelPitchButton = document.querySelector(".cancel-pitch-btn");
    const pitchEditActions = document.querySelector(".pitch-edit-actions");
    const titleText = document.querySelector(".pitch-title-text");
    const titleInput = document.querySelector(".pitch-title-input");
    const topicDisplay = document.querySelector(".pitch-topic-display");
    const topicInput = document.querySelector(".pitch-topic-input");
    const descriptionText = document.querySelector(".description");
    const descriptionEditor = document.querySelector(".description-editor");

    if (!editPitchButton || !savePitchButton || !cancelPitchButton) {
        return;
    }

    const originalValues = {
        title: titleText?.textContent.trim() || "",
        topic: topicDisplay?.textContent.trim() || "",
        description: descriptionText?.textContent.trim() || "",
    };

    const showEditor = () => {
        if (titleText) titleText.style.display = "none";
        if (titleInput) titleInput.style.display = "block";
        if (topicDisplay) topicDisplay.style.display = "none";
        if (topicInput) topicInput.style.display = "block";
        if (descriptionText) descriptionText.style.display = "none";
        if (descriptionEditor) descriptionEditor.style.display = "block";
        if (pitchEditActions) pitchEditActions.style.display = "block";
        if (titleInput) titleInput.focus();
    };

    const hideEditor = () => {
        if (titleText) titleText.style.display = "block";
        if (titleInput) titleInput.style.display = "none";
        if (topicDisplay) topicDisplay.style.display = "block";
        if (topicInput) topicInput.style.display = "none";
        if (descriptionText) descriptionText.style.display = "block";
        if (descriptionEditor) descriptionEditor.style.display = "none";
        if (pitchEditActions) pitchEditActions.style.display = "none";
    };

    editPitchButton.addEventListener("click", () => {
        if (titleInput) titleInput.value = originalValues.title;
        if (topicInput) topicInput.value = originalValues.topic;
        if (descriptionEditor) descriptionEditor.value = originalValues.description;
        showEditor();
    });

    cancelPitchButton.addEventListener("click", () => {
        hideEditor();
    });

    savePitchButton.addEventListener("click", async () => {
        const newTitle = titleInput?.value.trim() || "";
        const newDescription = descriptionEditor?.value.trim() || "";
        const newTopic = topicInput?.value.trim() || "";

        if (!newTitle || !newDescription || !newTopic) {
            alert("Please fill out the title, topic, and description.");
            return;
        }

        try {
            const response = await fetch(`/pitches/${pitchId}/edit`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ title: newTitle, description: newDescription, topic: newTopic }),
                credentials: "include",});
            const data = await response.json().catch(() => ({}));
            if (response.ok) {
                window.location.reload();
            } else {
                alert(data.message || "Failed to edit pitch.");}
        } catch (err) {
            console.error(err);
            alert("Unable to connect to the server.");}
    });
}

function setupCommentEditors() {
    const commentCards = document.querySelectorAll(".comment");
    commentCards.forEach((comment) => {
        const editButton = comment.querySelector(".edit-comment-btn");
        const saveButton = comment.querySelector(".save-comment-btn");
        const cancelButton = comment.querySelector(".cancel-comment-btn");
        const commentText = comment.querySelector(".comment-text");
        const commentEditor = comment.querySelector(".comment-editor");
        const commentEditActions = comment.querySelector(".comment-edit-actions");
        const commentId = editButton?.dataset.commentId;

        if (!editButton || !saveButton || !cancelButton || !commentText || !commentEditor || !commentEditActions) {
            return;
        }

        const originalValue = commentText.textContent.trim();

        const showEditor = () => {
            commentText.style.display = "none";
            commentEditor.style.display = "block";
            commentEditActions.style.display = "block";
            commentEditor.focus();
            commentEditor.setSelectionRange(commentEditor.value.length, commentEditor.value.length);
        };

        const hideEditor = () => {
            commentText.style.display = "block";
            commentEditor.style.display = "none";
            commentEditActions.style.display = "none";
            commentEditor.value = originalValue;
        };

        editButton.addEventListener("click", () => {
            commentEditor.value = originalValue;
            showEditor();
        });

        cancelButton.addEventListener("click", () => {
            hideEditor();
        });

        saveButton.addEventListener("click", async () => {
            const newContent = commentEditor.value.trim();
            if (!newContent) {
                alert("Comment cannot be empty.");
                return;
            }

            try {
                const response = await fetch(`/pitches/${pitchId}/comments/${commentId}/edit`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ content: newContent }),
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
        });
    });
}

checkAuth();
setupPitchEditor();
setupCommentEditors();

if (commentForm) {
    commentForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const body = { content: document.getElementById("content").value.trim() };
        try {
            const response = await fetch(`/pitches/${pitchId}/comments/add`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
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
}

const voteBtn = document.querySelector(".vote-btn");
if (voteBtn) {
    voteBtn.addEventListener("click", async (e) => {
        e.stopPropagation();
        try {
            const response = await fetch(`/pitches/${pitchId}/upvote`, {
                method: "POST",
                credentials: "include",
            });

            const data = await response.json().catch(() => ({}));
            if (response.ok) {
                voteBtn.classList.toggle("voted");
                voteBtn.querySelector(".vote-count").textContent = data.votes;
            } else {
                alert("Failed to upvote pitch.");
            }
        } catch (err) {
            console.error(err);
            alert("Unable to connect to the server.");
        }
    });
}

const deletePitchButton = document.querySelector(".delete-pitch-btn");
if (deletePitchButton) {
    deletePitchButton.addEventListener("click", async () => {
        if (confirm("Are you sure you want to delete this pitch?")) {
            try {
                const response = await fetch(`/pitches/${pitchId}/delete`, {
                    method: "DELETE",
                    credentials: "include",
                });
                const data = await response.json().catch(() => ({}));
                if (response.ok) {
                    window.location.href = "/";
                } else {
                    alert(data.message || "Failed to delete pitch.");
                }
            } catch (err) {
                console.error(err);
                alert("Unable to connect to the server.");
            }
        }
    });
}

const deleteCommentButtons = document.querySelectorAll(".delete-comment-btn");
deleteCommentButtons.forEach((button) => {
    button.addEventListener("click", async () => {
        const commentId = button.dataset.commentId;

        if (confirm("Are you sure you want to delete this comment?")) {
            try {
                const response = await fetch(`/pitches/${pitchId}/comments/${commentId}/delete`, {
                    method: "DELETE",
                    credentials: "include",});
                const data = await response.json().catch(() => ({}));
                if (response.ok) {
                    window.location.reload();
                } else {
                    alert(data.message || "Failed to delete comment.");}
            } catch (err) {
                console.error(err);
                alert("Unable to connect to the server.");}}
    });
});