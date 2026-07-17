import Toast from "./toast.js";

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
            Toast.success("Pitch created!");
            await new Promise(resolve => setTimeout(resolve, 250));
            window.location.href = `/pitches/${data.idea_id}`;
        } else {
            Toast.error(data.message || "Failed to create pitch.");
        }
    } catch (err) {
        console.error(err);
        Toast.error("Unable to connect to the server.");
    }
});