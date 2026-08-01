export interface SpatialNavOptions {
    enabled?: boolean;
}

export function spatialNavigation(node: HTMLElement, options: SpatialNavOptions = { enabled: true }) {
    let isEnabled = options.enabled ?? true;

    function handleKeyDown(e: KeyboardEvent) {
        if (!isEnabled) return;

        let direction: "up" | "right" | "down" | "left";

        switch (e.key) {
            case "ArrowUp":
                direction = "up";
                break;
            case "ArrowRight":
                direction = "right";
                break;
            case "ArrowDown":
                direction = "down";
                break;
            case "ArrowLeft":
                direction = "left";
                break;
            default: return;
        }

        const active = document.activeElement as HTMLElement;

        const focusables = Array.from(
            node.querySelectorAll<HTMLElement>("button:not([disabled]), [tabindex='0']:not([disabled])")
        );
        if (focusables.length === 0) return;

        if (!active || !node.contains(active)) {
            e.preventDefault();
            focusables[0].focus();
            return;
        }

        const currentRect = active.getBoundingClientRect();

        const target = findClosest(focusables, currentRect, direction);
        if (target) {
            e.preventDefault();
            target.focus();
        }
    }

    function findClosest(
        elements: HTMLElement[],
        currentRect: DOMRect,
        direction: "up" | "right" | "down" | "left"
    ): HTMLElement | null {
        let closest: HTMLElement | null = null;
        let minDistance = Infinity;

        for (const element of elements) {
            if (element === document.activeElement) continue;

            const rect = element.getBoundingClientRect();

            const isUp = direction === "up" && rect.bottom <= currentRect.top + 5;
            const isRight = direction === "right" && rect.left >= currentRect.right - 5;
            const isDown = direction === "down" && rect.top >= currentRect.bottom - 5;
            const isLeft = direction === "left" && rect.right <= currentRect.left + 5;

            if (isUp || isRight || isDown || isLeft) {
                const dx = rect.left + rect.width / 2 - (currentRect.left + currentRect.width / 2);
                const dy = rect.top + rect.height / 2 - (currentRect.top + currentRect.height / 2);
                const distance = Math.sqrt(dx * dx + dy * dy);

                if (distance < minDistance) {
                    minDistance = distance;
                    closest = element;
                }
            }
        }

        return closest;
    }

    window.addEventListener("keydown", handleKeyDown);

    return {
        update(newOptions: SpatialNavOptions) {
            isEnabled = newOptions.enabled ?? true;
        },
        destroy() {
            window.removeEventListener("keydown", handleKeyDown);
        }
    };
}