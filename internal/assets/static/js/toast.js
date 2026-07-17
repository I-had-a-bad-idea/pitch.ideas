// Inject styles once
function injectToastStyles() {
    if (document.getElementById("toast-styles")) return;

    const style = document.createElement("style");
    style.id = "toast-styles";

    style.textContent = `
        #toast-container {
            position: fixed;
            bottom: 24px;
            right: 24px;

            display: flex;
            flex-direction: column;
            gap: 12px;

            z-index: 9999;
        }

        .pitch-toast {
            min-width: 280px;
            max-width: 420px;

            padding: 14px 18px;

            border-radius: 10px;
            border: 1px solid var(--border);

            background: var(--white);
            color: var(--primary);

            font-family: 'Poppins', sans-serif;
            font-size: 0.95rem;
            font-weight: 500;

            display: flex;
            align-items: center;
            gap: 10px;

            box-shadow: 0 8px 25px rgba(15, 27, 61, 0.12);

            opacity: 0;
            transform: translateY(20px);

            transition:
                opacity 0.25s ease,
                transform 0.25s ease;
        }

        .pitch-toast.show {
            opacity: 1;
            transform: translateY(0);
        }

        .pitch-toast.success {
            border-left: 5px solid #16a34a;
        }

        .pitch-toast.error {
            border-left: 5px solid #dc2626;
        }

        .pitch-toast.warning {
            border-left: 5px solid var(--accent);
        }

        .pitch-toast.info {
            border-left: 5px solid #2563eb;
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
        toast.className = `pitch-toast ${type}`;
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