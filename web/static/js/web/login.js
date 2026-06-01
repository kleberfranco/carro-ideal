$(function () {
    const form = $("#login-form");
    if (!form.length) {
        return;
    }

    const button = $("#login-submit");
    const alertBox = $("#login-error");

    function clearValidation() {
        alertBox
            .removeClass("alert alert-danger alert-success")
            .text("")
            .hide();

        form.find(".is-invalid").removeClass("is-invalid");
        form.find(".is-valid").removeClass("is-valid");
        form.find(".invalid-feedback").text("");
    }

    form.on("submit", function (e) {
        e.preventDefault();

        clearValidation();

        button.prop("disabled", true);
        button.addClass("button-loading");

        const payload = {
            email: form.find("[name=email]").val(),
            password: form.find("[name=password]").val()
        };

        $.ajax({
            url: "/api/login",
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify(payload)
        })
            .done(function (data, _textStatus, jqXHR) {
                if (jqXHR.status === 200) {
                    form.find("input").addClass("is-valid");

                    alertBox
                        .addClass("alert alert-success")
                        .text("Login realizado com sucesso!")
                        .show();

                    // Armazenar dados do usuário no localStorage
                    localStorage.setItem("user", JSON.stringify(data.user));

                    setTimeout(function () {
                        window.location.href = "/";
                    }, 1200);
                }
            })
            .fail(function (xhr) {
                let data;
                try {
                    data = JSON.parse(xhr.responseText);
                } catch (e) {
                }

                if (xhr.status === 422 && data && data.errors) {
                    Object.keys(data.errors).forEach(function (field) {
                        const input = form.find('[name="' + field + '"]');
                        const feedback = input.siblings(".invalid-feedback");

                        input.addClass("is-invalid");
                        if (feedback.length) {
                            feedback.text(data.errors[field]);
                        }
                    });

                    alertBox
                        .addClass("alert alert-danger")
                        .text("Por favor, corrija os campos destacados.")
                        .show();

                    return;
                }

                if (xhr.status === 401 && data && data.error) {
                    alertBox
                        .addClass("alert alert-danger")
                        .text(data.error)
                        .show();

                    return;
                }

                let msg = "Erro ao fazer login. Tente novamente em instantes.";
                if (data && data.error) {
                    msg = data.error;
                }

                alertBox
                    .addClass("alert alert-danger")
                    .text(msg)
                    .show();
            })
            .always(function () {
                button.prop("disabled", false);
                button.removeClass("button-loading");
            });
    });
});

