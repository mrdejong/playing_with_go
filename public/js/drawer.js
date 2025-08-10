(function () {
    const drawers = new Map();

    // Update trigger states
    function updateTriggers(drawerId, isOpen) {
        document
            .querySelectorAll(`[data-drawer-trigger="${drawerId}"]`)
            .forEach((trigger) => {
                trigger.setAttribute("data-open", isOpen);
            });
    }

    // Get transform value based on position
    function getTransform(position, isOpen) {
        if (isOpen) return "translate(0)";

        switch (position) {
            case "left":
                return "translateX(-100%)";
            case "right":
                return "translateX(100%)";
            case "top":
                return "translateY(-100%)";
            case "bottom":
                return "translateY(100%)";
            default:
                return "translateX(100%)";
        }
    }

    // Create drawer instance
    function createDrawer(backdrop) {
        if (!backdrop || backdrop.hasAttribute("data-initialized")) return null;
        backdrop.setAttribute("data-initialized", "true");

        const drawerId = backdrop.id;
        const content = document.getElementById(drawerId + "-content");
        const content2 = document.getElementsByClassName('.select-item');
        const position = content?.getAttribute("data-drawer-position") || "right";
        const isInitiallyOpen = backdrop.hasAttribute("data-initial-open");

        if (!content || !drawerId) return null;

        let isOpen = isInitiallyOpen;

        // Set initial state
        function setState(open) {
            isOpen = open;
            const display = open ? "block" : "none";
            const opacity = open ? "1" : "0";

            if (open) {
                backdrop.style.display = display;
                content.style.display = display;

                // Force reflow
                content.offsetHeight;

                // Add transitions
                backdrop.style.transition = "opacity 300ms ease";
                content.style.transition = "opacity 300ms ease, transform 300ms ease";
            }

            backdrop.style.opacity = opacity;
            content.style.opacity = opacity;
            content.style.transform = getTransform(position, open);

            backdrop.setAttribute("data-open", open);
            updateTriggers(drawerId, open);

            document.body.style.overflow = open ? "hidden" : "";
        }

        // Open drawer
        function open() {
            setState(true);

            // Add event listeners
            backdrop.addEventListener("click", close);
            document.addEventListener("keydown", handleEsc);
            document.addEventListener("click", handleClickAway);
            document.addEventListener("close", close);
        }

        // Close drawer
        function close() {
            setState(false);

            // Remove event listeners
            backdrop.removeEventListener("click", close);
            document.removeEventListener("keydown", handleEsc);
            document.removeEventListener("click", handleClickAway);
            document.removeEventListener("close", close);

            // Hide after animation
            setTimeout(() => {
                if (!isOpen) {
                    backdrop.style.display = "none";
                    content.style.display = "none";
                }
            }, 300);
        }

        // Toggle drawer
        function toggle() {
            isOpen ? close() : open();
        }

        // Handle escape key
        function handleEsc(e) {
            if (e.key === "Escape" && isOpen) close();
        }

        // Handle click away
        function handleClickAway(e) {
            if (document.querySelector('div[data-popover-portal-container]').contains(e.target)) {
                return;
            }

            // data-popover-portal-container
            if (!content.contains(e.target) && !isTriggerClick(e.target)) {
                close();
            }
        }

        // Check if click is on a trigger
        function isTriggerClick(target) {
            const triggers = document.querySelectorAll(
                `[data-drawer-trigger="${drawerId}"]`
            );
            return Array.from(triggers).some((trigger) => trigger.contains(target));
        }

        // Setup close buttons
        content.querySelectorAll("[data-drawer-close]").forEach((btn) => {
            btn.addEventListener("click", close);
        });

        // Prevent backdrop clicks on content
        content
            .querySelector("[data-drawer-inner]")
            ?.addEventListener("click", (e) => {
                e.stopPropagation();
            });

        // Set initial state
        setState(isInitiallyOpen);

        return { open, close, toggle };
    }

    // Initialize all drawers and triggers
    function init(root = document) {
        // Find and initialize drawers
        root.querySelectorAll('[data-component="drawer"]:not([data-initialized])').forEach((backdrop) => {
            const drawer = createDrawer(backdrop);
            if (drawer && backdrop.id) {
                drawers.set(backdrop.id, drawer);
            }
        });

        // Setup trigger clicks
        root.querySelectorAll("[data-drawer-trigger]").forEach((trigger) => {
            if (trigger.dataset.initialized) return;
            trigger.dataset.initialized = "true";

            const drawerId = trigger.getAttribute("data-drawer-trigger");
            trigger.addEventListener("click", () => {
                drawers.get(drawerId)?.toggle();
            });
        });
    }

    // Export
    window.templUI = window.templUI || {};
    window.templUI.drawer = { init: init };

    // Auto-initialize
    document.addEventListener("DOMContentLoaded", () => init());
})();