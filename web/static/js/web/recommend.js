$(function () {
    const alertBox = $("#recommend-alert");
    const form = $("#questionnaire-form");
    const submitButton = $("#generate-submit");
    const vehicleModal = new bootstrap.Modal(document.getElementById("vehicle-modal"));
    const comparisonModal = new bootstrap.Modal(document.getElementById("comparison-modal"));
    let currentItems = [];

    function showAlert(message, type) {
        alertBox
            .removeClass("d-none alert-danger alert-success alert-warning")
            .addClass("alert-" + type)
            .text(message);
    }

    function clearAlert() {
        alertBox.addClass("d-none").text("");
    }

    function showView(view) {
        clearAlert();
        $("#questionnaire-view, #results-view, #history-view").addClass("d-none");
        $("#" + view + "-view").removeClass("d-none");
        $("#recommend-tabs .nav-link").removeClass("active");
        $('#recommend-tabs .nav-link[data-view="' + view + '"]').addClass("active");
        if (view === "history") {
            loadHistory();
        }
    }

    function escapeHtml(value) {
        return $("<div>").text(value || "").html();
    }

    function money(value) {
        return Number(value).toLocaleString("pt-BR", {
            style: "currency",
            currency: "BRL",
            maximumFractionDigits: 0
        });
    }

    function loadQuestions() {
        $.getJSON("/api/questions/")
            .done(function (response) {
                const questions = response.data.items || [];
                const html = questions.map(function (question, questionIndex) {
                    const options = question.answer_options.map(function (option) {
                        return `
                            <label class="answer-option">
                                <input class="form-check-input" type="radio"
                                    name="question_${question.id}"
                                    value="${option.id}"
                                    data-question-id="${question.id}">
                                <span>${escapeHtml(option.text)}</span>
                            </label>`;
                    }).join("");

                    return `
                        <fieldset class="question-card">
                            <legend>
                                <span class="question-number">${questionIndex + 1}</span>
                                ${escapeHtml(question.text)}
                            </legend>
                            <div class="answer-grid">${options}</div>
                        </fieldset>`;
                }).join("");

                $("#questions-list").html(html);
                $("#questionnaire-loading").addClass("d-none");
                form.removeClass("d-none");
            })
            .fail(function (xhr) {
                if (xhr.status === 401) {
                    window.location.href = "/login";
                    return;
                }
                showAlert("Não foi possível carregar o questionário.", "danger");
            });
    }

    function renderRecommendation(recommendation) {
        currentItems = recommendation.items || [];
        const cards = currentItems.map(function (item) {
            const vehicle = item.vehicle;
            return `
                <div class="col-12">
                    <article class="vehicle-card">
                        <div class="rank-badge">${item.rank}</div>
                        <div class="vehicle-main">
                            <div class="d-flex flex-wrap justify-content-between gap-2">
                                <div>
                                    <p class="vehicle-category">${escapeHtml(vehicle.category.name)}</p>
                                    <h3>${escapeHtml(vehicle.brand)} ${escapeHtml(vehicle.model)}</h3>
                                    <p class="text-secondary mb-2">${escapeHtml(vehicle.version)} · ${vehicle.year}</p>
                                </div>
                                <div class="score-ring">${Number(item.score).toFixed(0)}<small>%</small></div>
                            </div>
                            <p class="reason">${escapeHtml(item.reason)}</p>
                            <div class="vehicle-facts">
                                <span>${escapeHtml(vehicle.transmission)}</span>
                                <span>${escapeHtml(vehicle.fuel_type)}</span>
                                <span>${money(vehicle.price_min)} a ${money(vehicle.price_max)}</span>
                            </div>
                            <button class="btn btn-outline-primary w-auto mt-3 vehicle-details"
                                data-vehicle-id="${vehicle.id}">Ver detalhes</button>
                            <label class="compare-choice mt-3">
                                <input class="form-check-input compare-vehicle" type="checkbox"
                                    value="${vehicle.id}">
                                Comparar
                            </label>
                        </div>
                    </article>
                </div>`;
        }).join("");
        $("#recommendation-results").html(cards);
        updateComparison();
        showView("results");
    }

    function updateComparison() {
        let selected = $(".compare-vehicle:checked");
        if (selected.length > 3) {
            selected.last().prop("checked", false);
            selected = $(".compare-vehicle:checked");
            showAlert("Compare no máximo três veículos.", "warning");
        }
        $("#comparison-count").text(
            selected.length ? selected.length + " veículo(s) selecionado(s)." : "Selecione de 2 a 3 veículos para comparar."
        );
        $("#compare-selected").prop("disabled", selected.length < 2);
    }

    $(document).on("change", ".compare-vehicle", updateComparison);

    $("#compare-selected").on("click", function () {
        const ids = $(".compare-vehicle:checked").map(function () {
            return Number($(this).val());
        }).get();
        const items = currentItems.filter(function (item) {
            return ids.includes(item.vehicle.id);
        });
        const cell = function (value, best) {
            return `<td class="${best ? "comparison-best" : ""}">${value}</td>`;
        };
        const highestScore = Math.max.apply(null, items.map(item => Number(item.score)));
        const lowestPrice = Math.min.apply(null, items.map(item => Number(item.vehicle.price_min)));
        $("#comparison-body").html(`
            <table class="table comparison-table align-middle">
                <thead><tr><th>Critério</th>${items.map(item => `<th>${escapeHtml(item.vehicle.brand)} ${escapeHtml(item.vehicle.model)}</th>`).join("")}</tr></thead>
                <tbody>
                    <tr><th>Compatibilidade</th>${items.map(item => cell(Number(item.score).toFixed(0) + "%", Number(item.score) === highestScore)).join("")}</tr>
                    <tr><th>Preço inicial</th>${items.map(item => cell(money(item.vehicle.price_min), Number(item.vehicle.price_min) === lowestPrice)).join("")}</tr>
                    <tr><th>Câmbio</th>${items.map(item => cell(escapeHtml(item.vehicle.transmission), false)).join("")}</tr>
                    <tr><th>Combustível</th>${items.map(item => cell(escapeHtml(item.vehicle.fuel_type), false)).join("")}</tr>
                    <tr><th>Lugares</th>${items.map(item => cell(item.vehicle.seats, false)).join("")}</tr>
                    <tr><th>Porta-malas</th>${items.map(item => cell(item.vehicle.trunk_capacity + " L", false)).join("")}</tr>
                    <tr><th>Por que combina</th>${items.map(item => cell(escapeHtml(item.reason), false)).join("")}</tr>
                </tbody>
            </table>`);
        comparisonModal.show();
    });

    form.on("submit", function (event) {
        event.preventDefault();
        clearAlert();

        const answers = [];
        $("#questions-list fieldset").each(function () {
            const selected = $(this).find("input:checked");
            if (selected.length) {
                answers.push({
                    question_id: Number(selected.data("question-id")),
                    answer_option_id: Number(selected.val())
                });
            }
        });

        submitButton.prop("disabled", true).text("Calculando...");
        $.ajax({
            url: "/api/recommendations/generate",
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify({answers: answers})
        })
            .done(function (response) {
                renderRecommendation(response.data);
            })
            .fail(function (xhr) {
                if (xhr.status === 401) {
                    window.location.href = "/login";
                    return;
                }
                const message = xhr.responseJSON && xhr.responseJSON.error
                    ? xhr.responseJSON.error
                    : "Não foi possível gerar as recomendações.";
                showAlert(message, "danger");
            })
            .always(function () {
                submitButton.prop("disabled", false).text("Gerar minhas recomendações");
            });
    });

    function loadHistory() {
        $("#history-list").html('<p class="text-secondary">Carregando histórico...</p>');
        $.getJSON("/api/recommendations/")
            .done(function (response) {
                const items = response.data.items || [];
                if (!items.length) {
                    $("#history-list").html('<div class="empty-state">Você ainda não gerou recomendações.</div>');
                    return;
                }
                const html = items.map(function (item) {
                    const date = new Date(item.created_at).toLocaleString("pt-BR");
                    return `
                        <button class="history-item" data-recommendation-id="${item.id}">
                            <span>
                                <strong>Recomendação #${item.id}</strong>
                                <small>${date}</small>
                            </span>
                            <span>${item.item_count} veículos</span>
                        </button>`;
                }).join("");
                $("#history-list").html(html);
            })
            .fail(function () {
                showAlert("Não foi possível carregar o histórico.", "danger");
            });
    }

    $(document).on("click", "[data-view]", function () {
        showView($(this).data("view"));
    });

    $(document).on("click", ".history-item", function () {
        const id = $(this).data("recommendation-id");
        $.getJSON("/api/recommendations/" + id)
            .done(function (response) {
                renderRecommendation(response.data);
            })
            .fail(function () {
                showAlert("Não foi possível abrir essa recomendação.", "danger");
            });
    });

    $(document).on("click", ".vehicle-details", function () {
        const id = $(this).data("vehicle-id");
        $.getJSON("/api/vehicles/" + id)
            .done(function (response) {
                const vehicle = response.data;
                $("#vehicle-modal-title").text(vehicle.brand + " " + vehicle.model);
                $("#vehicle-modal-body").html(`
                    <p class="lead">${escapeHtml(vehicle.description)}</p>
                    <div class="detail-grid">
                        <div><small>Versão</small><strong>${escapeHtml(vehicle.version)}</strong></div>
                        <div><small>Ano</small><strong>${vehicle.year}</strong></div>
                        <div><small>Câmbio</small><strong>${escapeHtml(vehicle.transmission)}</strong></div>
                        <div><small>Combustível</small><strong>${escapeHtml(vehicle.fuel_type)}</strong></div>
                        <div><small>Porta-malas</small><strong>${vehicle.trunk_capacity} L</strong></div>
                        <div><small>Consumo urbano</small><strong>${vehicle.consumption_city} km/l</strong></div>
                    </div>
                    <h3 class="fs-5 mt-4">Pontos fortes</h3>
                    <p>${escapeHtml(vehicle.strengths)}</p>
                    <h3 class="fs-5">Pontos de atenção</h3>
                    <p>${escapeHtml(vehicle.weaknesses)}</p>
                `);
                vehicleModal.show();
            })
            .fail(function () {
                showAlert("Não foi possível carregar os detalhes do veículo.", "danger");
            });
    });

    loadQuestions();
});
