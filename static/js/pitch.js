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
        alert("Comment added successfully!");
        window.location.reload();
    } else {
        alert(data.message || "Failed to add comment.");
    }
    } catch (err) {
        console.error(err);
        alert("Unable to connect to the server.");
    }
});