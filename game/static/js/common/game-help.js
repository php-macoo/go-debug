(function () {
    var overlay = document.getElementById("gameHelpOverlay");
    if (!overlay) return;

    var openBtns = document.querySelectorAll("[data-game-help-open]");
    var closeBtns = overlay.querySelectorAll("[data-game-help-close]");
    var panel = overlay.querySelector(".game-help-panel");

    function open() {
        overlay.hidden = false;
        overlay.setAttribute("aria-hidden", "false");
        document.body.style.overflow = "hidden";
        if (panel && typeof panel.focus === "function") {
            try { panel.focus(); } catch (e) {}
        }
    }

    function close() {
        overlay.hidden = true;
        overlay.setAttribute("aria-hidden", "true");
        document.body.style.overflow = "";
    }

    openBtns.forEach(function (b) {
        b.addEventListener("click", function (e) {
            e.preventDefault();
            open();
        });
    });
    closeBtns.forEach(function (b) {
        b.addEventListener("click", function (e) {
            e.preventDefault();
            close();
        });
    });
    overlay.addEventListener("click", function (e) {
        if (e.target === overlay) close();
    });
    document.addEventListener("keydown", function (e) {
        if (e.key === "Escape" && !overlay.hidden) {
            e.preventDefault();
            close();
        }
    });
})();
