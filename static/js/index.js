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
            <span>👍 ${p.votes}</span>
            <span>💬 ${p.comment_count}</span>
        </div>
        `;

        container.appendChild(div);
    });
}

loadPitches();