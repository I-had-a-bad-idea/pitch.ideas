async function loadPitches() {
    const res = await fetch("/pitches");
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

        <p>${p.description}</p>
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
                    method: "POST"
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

loadPitches();