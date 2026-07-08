const form = document.querySelector(".pitch-form");

form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const body = {
        title: document.getElementById("title").value.trim(),
        topic: document.getElementById("topic").value,
        description: document.getElementById("description").value.trim()
    };

    try {
        const response = await fetch("/pitches/create", {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(body),
            credentials: "include",
        });

        const data = await response.json().catch(() => ({}));

        if (response.ok) {
            window.location.href = "/";
        } else {
            alert(data.message || "Failed to create pitch.");
        }
    } catch (err) {
        console.error(err);
        alert("Unable to connect to the server.");
    }
});