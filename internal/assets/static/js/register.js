import Toast from "./toast.js";

const form = document.querySelector(".auth-form");

form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const body = {
        username: document.getElementById("username").value.trim(),
        password: document.getElementById("password").value
    };

    try {
        const response = await fetch("/auth/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(body),
            credentials: "include",
        });

        const data = await response.json().catch(() => ({}));

        if (response.ok) {
            Toast.success("Registered!");
            await new Promise(resolve => setTimeout(resolve, 250));
            window.location.href = "/";
        } else {
            Toast.error(data.message || "Registration failed.");
        }

    } catch (err) {
        console.error(err);
        Toast.error("Unable to connect to the server.");
    }
});