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

const comments = [];
const authors = [];
const dates = [];

const voteBtn = document.querySelector(".vote-btn");
var votes = parseInt(voreBtn.textContent.split(" ")[1]);

const totalLikes = 150;
const totalComments = 50; // should always be less then likes 
const waitMsBetweenComments = 10;
const likesPerComment = totalLikes / totalComments;

async function flood_with_likes_and_comments() {
    for (let i = 0; i < totalComments; i++) {
        await new Promise(resolve => setTimeout(resolve, waitMsBetweenComments)); // wait
        const author = authors[Math.floor(Math.random() * authors.length)];
        const comment = comments[Math.floor(Math.random() * comments.length)];
        const date = dates[Math.floor(Math.random() * dates.length)];

        const li = document.createElement("li");
        li.className = "comment"

        li.innerHTML = `
        <p class="comment-text">${comment}</p>
        <p class"comment-meta">By ${author } on ${date}</p>
        `;
    }
}

// After showing the pitch, we want to flood it (for the ad)
flood_with_likes_and_comments();