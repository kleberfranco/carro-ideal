function cookieValue(name) {
    const prefix = name + "=";
    const item = document.cookie.split(";").map(function (part) {
        return part.trim();
    }).find(function (part) {
        return part.startsWith(prefix);
    });
    return item ? decodeURIComponent(item.substring(prefix.length)) : "";
}

$.ajaxPrefilter(function (options, originalOptions, xhr) {
    const method = (options.type || "GET").toUpperCase();
    if (!["GET", "HEAD", "OPTIONS", "TRACE"].includes(method)) {
        xhr.setRequestHeader("X-CSRF-Token", cookieValue("carro_csrf"));
    }
});

$(function () {
    const root = document.documentElement;
    const theme = localStorage.getItem("carro-ideal-theme") || "light";
    root.dataset.theme = theme;
    $("#theme-toggle").attr("aria-pressed", theme === "dark").text(theme === "dark" ? "Tema claro" : "Tema escuro");

    $("#theme-toggle").on("click", function () {
        const next = root.dataset.theme === "dark" ? "light" : "dark";
        root.dataset.theme = next;
        localStorage.setItem("carro-ideal-theme", next);
        $(this).attr("aria-pressed", next === "dark").text(next === "dark" ? "Tema claro" : "Tema escuro");
    });

    if ($("#admin-nav-item").length) {
        $.getJSON("/api/auth/me").done(function (response) {
            if (response.data.user.role.toLowerCase() === "admin") {
                $("#admin-nav-item").removeClass("d-none");
            }
        });
    }

    $(document).ajaxError(function (_event, xhr, settings) {
        if (settings.global === false || xhr.status === 401 || xhr.status === 403) {
            return;
        }
        if (!$("#global-api-error").length) {
            $("body > .container").prepend(
                '<div id="global-api-error" class="alert alert-danger" role="alert">Não foi possível concluir uma solicitação. Tente novamente.</div>'
            );
        }
    });
});
