$(function () {
    const form = $("#register-form");
    if (!form.length) {
        return;
    }

    const button = $("#register-submit");
    const alertBox = $("#register-error");

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
            name: form.find("[name=name]").val(),
            email: form.find("[name=email]").val(),
            password: form.find("[name=password]").val(),
            confirm_password: form.find("[name=confirm_password]").val()
        };

        $.ajax({
            url: "/api/register",
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify(payload)
        })
            .done(function (_, _textStatus, jqXHR) {
                if (jqXHR.status === 201) {
                    form.find("input").addClass("is-valid");

                    alertBox
                        .addClass("alert alert-success")
                        .text("Conta criada com sucesso! Você já pode entrar.")
                        .show();

                    setTimeout(function () {
                        window.location.href = "/login";
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

                if (xhr.status === 409 && data && data.error) {
                    const emailInput = form.find('[name="email"]');
                    emailInput.addClass("is-invalid");
                    emailInput.siblings(".invalid-feedback").text(data.error);

                    alertBox
                        .addClass("alert alert-danger")
                        .text(data.error)
                        .show();

                    return;
                }

                let msg = "Erro ao cadastrar. Tente novamente em instantes.";
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