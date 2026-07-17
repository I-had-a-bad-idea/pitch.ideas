// Inject styles once
function injectToastStyles() {
    if (document.getElementById("toast-styles")) return;

    const style = document.createElement("style");
    style.id = "toast-styles";

    style.textContent = `
        #toast-container {
            position: fixed;
            bottom: 20px;
            right: 20px;

            display: flex;
            flex-direction: column;
            gap: 10px;

            z-index: 9999;
        }

        .toast {
            min-width: 250px;
            max-width: 400px;

            padding: 14px 18px;
            border-radius: 8px;

            color: white;
            font-family: sans-serif;

            opacity: 0;
            transform: translateY(20px);

            transition: opacity 0.25s ease,
                        transform 0.25s ease;
        }

        .toast.show {
            opacity: 1;
            transform: translateY(0);
        }

        .toast.success {
            background: #16a34a;
        }

        .toast.error {
            background: #dc2626;
        }

        .toast.warning {
            background: #d97706;
        }

        .toast.info {
            background: #2563eb;
        }
    `;

    document.head.appendChild(style);
}


class Toast {
    static init() {
        injectToastStyles();
    }

    static show(message, type = "info", duration = 3000) {
        Toast.init();

        let container = document.getElementById("toast-container");

        if (!container) {
            container = document.createElement("div");
            container.id = "toast-container";
            document.body.appendChild(container);
        }

        const toast = document.createElement("div");
        toast.className = `toast ${type}`;
        toast.textContent = message;

        container.appendChild(toast);

        requestAnimationFrame(() => {
            toast.classList.add("show");
        });

        setTimeout(() => {
            toast.classList.remove("show");

            setTimeout(() => {
                toast.remove();
            }, 250);

        }, duration);
    }

    static success(message, duration) {
        this.show(message, "success", duration);
    }

    static error(message, duration) {
        this.show(message, "error", duration);
    }

    static warning(message, duration) {
        this.show(message, "warning", duration);
    }

    static info(message, duration) {
        this.show(message, "info", duration);
    }
}

export default Toast;