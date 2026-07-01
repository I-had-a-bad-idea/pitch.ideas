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

function randomDateTime() {
    const now = new Date();

    // Anywhere from now to 30 minutes ago
    const date = new Date(now.getTime() - Math.random() * 30 * 60 * 1000);

    const months = [
        "Jan", "Feb", "Mar", "Apr", "May", "Jun",
        "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"
    ];

    const day = date.getDate();
    const month = months[date.getMonth()];
    const year = date.getFullYear();
    const hours = String(date.getHours()).padStart(2, "0");
    const minutes = String(date.getMinutes()).padStart(2, "0");

    return `${day} ${month} ${year}, ${hours}:${minutes}`;
}

const comments = [
    "This is exactly what I needed!",
    "Amazing idea 👏",
    "I'd definitely use this.",
    "Take my upvote!",
    "This deserves more attention.",
    "Really well thought out.",
    "Love this concept.",
    "Great work!",
    "I'd invest in this.",
    "Fantastic pitch!"
];
const authors = [
    "Alex",
    "Sarah",
    "Michael",
    "Emma",
    "Noah",
    "Olivia",
    "Liam",
    "Sophia",
    "Daniel",
    "Mia",
    "James",
    "Ava"
];

const voteBtn = document.querySelector(".vote-btn");
var votes = parseInt(voreBtn.textContent.split(" ")[1]);
votes = 0;
const commentCount = document.querySelector(".comment-count");
const commentList = document.querySelector(".comment-list");

const totalLikes = 231;
const totalComments = 57; // should always be less then likes 
const waitMsBetweenComments = 100;
const likesPerComment = totalLikes / totalComments;

async function flood_with_likes_and_comments() {
    await new Promise(resolve => setTimeout(resolve, waitMsBetweenComments)); // wait before starting
    for (let i = 0; i < totalComments; i++) {
        const author = authors[Math.floor(Math.random() * authors.length)];
        const comment = comments[Math.floor(Math.random() * comments.length)];
        const date = randomDateTime();

        const li = document.createElement("li");
        li.className = "comment"

        li.innerHTML = `
        <p class="comment-text">${comment}</p>
        <p class="comment-meta">By ${author } on ${date}</p>
        `;

        commentList.prepend(li);
        votes += likesPerComment;
        voteBtn.textContent = `👍 ${Math.round(votes)}`;
        commentCount.textContent = `💬 ${i + 1}`;

        await new Promise(resolve => setTimeout(resolve, waitMsBetweenComments)); // wait
    }
}

// After showing the pitch, we want to flood it (for the ad)
flood_with_likes_and_comments();
