$(function () {
    const state = {
        vehicles: [],
        categories: [],
        questions: [],
        vehiclePage: 1,
        vehicleLimit: 10,
        vehicleTotal: 0,
        vehicleSearch: "",
        users: [],
        userPage: 1,
        userLimit: 10,
        userTotal: 0,
        userSearch: ""
    };
    const vehicleModal = new bootstrap.Modal("#vehicle-form-modal");
    const categoryModal = new bootstrap.Modal("#category-form-modal");
    const questionModal = new bootstrap.Modal("#question-form-modal");
    const optionModal = new bootstrap.Modal("#option-form-modal");

    function escapeHtml(value) {
        return $("<div>").text(value == null ? "" : value).html();
    }

    function showAlert(message, type) {
        $("#admin-alert").removeClass("d-none alert-danger alert-success").addClass("alert-" + type).text(message);
    }

    function request(options) {
        return $.ajax(options).fail(function (xhr) {
            if (xhr.status === 401) {
                window.location.href = "/login";
                return;
            }
            const message = xhr.responseJSON && xhr.responseJSON.error
                ? xhr.responseJSON.error
                : "Não foi possível concluir a operação.";
            showAlert(message, "danger");
        });
    }

    function showView(view) {
        $("[id^=admin-][id$=-view]").addClass("d-none");
        $("#admin-" + view + "-view").removeClass("d-none");
        $("#admin-tabs .nav-link").removeClass("active");
        $('#admin-tabs [data-admin-view="' + view + '"]').addClass("active");
        if (view === "dashboard") loadDashboard();
        if (view === "users") loadUsers();
        if (view === "vehicles") loadVehicles();
        if (view === "categories") loadCategories();
        if (view === "questions") loadQuestions();
    }

    function roleBadge(role) {
        return role === "admin"
            ? '<span class="badge text-bg-primary">Admin</span>'
            : '<span class="badge text-bg-secondary">Usuário</span>';
    }

    function loadUsers() {
        const query = $.param({
            page: state.userPage,
            limit: state.userLimit,
            search: state.userSearch
        });
        request({url: "/api/admin/users?" + query}).done(function (response) {
            state.users = response.data.items || [];
            state.userTotal = response.data.total;
            $("#users-table").html(state.users.map(function (user) {
                const date = new Date(user.created_at).toLocaleDateString("pt-BR");
                return `<tr>
                    <td><strong>${escapeHtml(user.name)}</strong></td>
                    <td>${escapeHtml(user.email)}</td>
                    <td>${roleBadge(user.role)}</td>
                    <td>${statusBadge(user.active)}</td>
                    <td>${date}</td>
                </tr>`;
            }).join(""));
            const totalPages = Math.max(1, Math.ceil(state.userTotal / state.userLimit));
            $("#users-page").text("Página " + state.userPage + " de " + totalPages);
            $("#users-prev").prop("disabled", state.userPage <= 1);
            $("#users-next").prop("disabled", state.userPage >= totalPages);
        });
    }

    function loadDashboard() {
        request({url: "/api/admin/dashboard"}).done(function (response) {
            const stats = response.data;
            const cards = [
                ["Usuários", stats.users, "users"],
                ["Veículos", stats.vehicles, "vehicles"],
                ["Recomendações", stats.recommendations, null],
                ["Perguntas", stats.questions, "questions"],
                ["Usuários ativos (7 dias)", stats.active_users_week, "users"],
                ["Novos usuários (7 dias)", stats.new_users_week, "users"]
            ].map(function (item) {
                const nav = item[2] ? ` data-admin-nav="${item[2]}"` : "";
                const cls = item[2] ? " stat-card--link" : "";
                return `<div class="col-lg-4 col-6"><article class="stat-card${cls}"${nav}><strong>${item[1]}</strong><span>${item[0]}</span></article></div>`;
            }).join("");
            $("#admin-stats").html(cards);
        });
    }

    function loadCategories(callback) {
        request({url: "/api/admin/categories"}).done(function (response) {
            state.categories = response.data.items || [];
            $("#categories-table").html(state.categories.map(function (category) {
                return `<tr>
                    <td><strong>${escapeHtml(category.name)}</strong><small>${escapeHtml(category.description)}</small></td>
                    <td>${category.vehicle_count}</td>
                    <td>${statusBadge(category.active)}</td>
                    <td class="text-end">${actions("category", category.id, category.active)}</td>
                </tr>`;
            }).join(""));
            const options = state.categories.filter(function (item) { return item.active; }).map(function (item) {
                return `<option value="${item.id}">${escapeHtml(item.name)}</option>`;
            }).join("");
            $("#vehicle-form [name=category_id]").html(options);
            if (callback) callback();
        });
    }

    function loadVehicles() {
        const query = $.param({
            page: state.vehiclePage,
            limit: state.vehicleLimit,
            search: state.vehicleSearch
        });
        request({url: "/api/admin/vehicles?" + query}).done(function (response) {
            state.vehicles = response.data.items || [];
            state.vehicleTotal = response.data.total;
            $("#vehicles-table").html(state.vehicles.map(function (vehicle) {
                return `<tr>
                    <td><strong>${escapeHtml(vehicle.brand)} ${escapeHtml(vehicle.model)}</strong><small>${escapeHtml(vehicle.version)} · ${vehicle.year}</small></td>
                    <td>${escapeHtml(vehicle.category.name)}</td>
                    <td>R$ ${Number(vehicle.price_min).toLocaleString("pt-BR")} - R$ ${Number(vehicle.price_max).toLocaleString("pt-BR")}</td>
                    <td>${statusBadge(vehicle.active)}</td>
                    <td class="text-end">${actions("vehicle", vehicle.id, vehicle.active)}</td>
                </tr>`;
            }).join(""));
            const totalPages = Math.max(1, Math.ceil(state.vehicleTotal / state.vehicleLimit));
            $("#vehicles-page").text("Página " + state.vehiclePage + " de " + totalPages);
            $("#vehicles-prev").prop("disabled", state.vehiclePage <= 1);
            $("#vehicles-next").prop("disabled", state.vehiclePage >= totalPages);
        });
    }

    function loadQuestions() {
        request({url: "/api/admin/questions"}).done(function (response) {
            state.questions = response.data.items || [];
            renderQuestions();
        });
    }

    function renderQuestions() {
        const search = $("#question-search").val().toLowerCase();
        const questions = state.questions.filter(function (question) {
            return question.text.toLowerCase().includes(search);
        });
        $("#questions-admin-list").html(questions.map(function (question) {
                const options = question.answer_options.map(function (option) {
                    return `<li>
                        <span>${escapeHtml(option.text)} ${statusBadge(option.active)}</span>
                        <span>
                            <button class="btn btn-sm btn-outline-primary edit-option" data-question-id="${question.id}" data-id="${option.id}">Editar</button>
                            ${option.active ? `<button class="btn btn-sm btn-outline-danger delete-option" data-question-id="${question.id}" data-id="${option.id}">Desativar</button>` : ""}
                        </span>
                    </li>`;
                }).join("");
                return `<article class="admin-question-card">
                    <div class="admin-heading mb-3">
                        <div><span class="question-order">#${question.display_order}</span><strong>${escapeHtml(question.text)}</strong><small>Peso ${question.weight} · ${question.answer_options.length} opções</small></div>
                        <div>${statusBadge(question.active)} ${actions("question", question.id, question.active)}</div>
                    </div>
                    <ul class="admin-options">${options}</ul>
                    <button class="btn btn-sm btn-outline-primary w-auto new-option" data-question-id="${question.id}">Adicionar opção</button>
                </article>`;
        }).join(""));
    }

    function statusBadge(active) {
        return active
            ? '<span class="badge text-bg-success">Ativo</span>'
            : '<span class="badge text-bg-secondary">Inativo</span>';
    }

    function actions(type, id, active) {
        return `<button class="btn btn-sm btn-outline-primary edit-${type}" data-id="${id}">Editar</button>
            ${active ? `<button class="btn btn-sm btn-outline-danger delete-${type}" data-id="${id}">Desativar</button>` : ""}`;
    }

    function value(form, name) {
        return form.find("[name=" + name + "]").val();
    }

    $("#admin-tabs").on("click", "[data-admin-view]", function () {
        showView($(this).data("admin-view"));
    });

    $("#admin-stats").on("click", "[data-admin-nav]", function () {
        showView($(this).data("admin-nav"));
    });

    let vehicleSearchTimer;
    $("#vehicle-search").on("input", function () {
        clearTimeout(vehicleSearchTimer);
        state.vehicleSearch = $(this).val();
        state.vehiclePage = 1;
        vehicleSearchTimer = setTimeout(loadVehicles, 250);
    });

    $("#vehicles-prev").on("click", function () {
        if (state.vehiclePage > 1) {
            state.vehiclePage -= 1;
            loadVehicles();
        }
    });

    $("#vehicles-next").on("click", function () {
        if (state.vehiclePage * state.vehicleLimit < state.vehicleTotal) {
            state.vehiclePage += 1;
            loadVehicles();
        }
    });

    let userSearchTimer;
    $("#user-search").on("input", function () {
        clearTimeout(userSearchTimer);
        state.userSearch = $(this).val();
        state.userPage = 1;
        userSearchTimer = setTimeout(loadUsers, 250);
    });

    $("#users-prev").on("click", function () {
        if (state.userPage > 1) {
            state.userPage -= 1;
            loadUsers();
        }
    });

    $("#users-next").on("click", function () {
        if (state.userPage * state.userLimit < state.userTotal) {
            state.userPage += 1;
            loadUsers();
        }
    });

    $("#question-search").on("input", renderQuestions);

    $("#new-category").on("click", function () {
        $("#category-form")[0].reset();
        $("#category-form [name=id]").val("");
        categoryModal.show();
    });

    $(document).on("click", ".edit-category", function () {
        const category = state.categories.find(item => item.id === Number($(this).data("id")));
        $("#category-form [name=id]").val(category.id);
        $("#category-form [name=name]").val(category.name);
        $("#category-form [name=description]").val(category.description);
        categoryModal.show();
    });

    $("#category-form").on("submit", function (event) {
        event.preventDefault();
        const form = $(this);
        const id = value(form, "id");
        request({
            url: "/api/admin/categories" + (id ? "/" + id : ""),
            method: id ? "PUT" : "POST",
            contentType: "application/json",
            data: JSON.stringify({name: value(form, "name"), description: value(form, "description")})
        }).done(function () {
            categoryModal.hide(); loadCategories(); showAlert("Categoria salva.", "success");
        });
    });

    $("#new-vehicle").on("click", function () {
        loadCategories(function () {
            $("#vehicle-form")[0].reset();
            $("#vehicle-form [name=id]").val("");
            $("#vehicle-form [name=year]").val(new Date().getFullYear());
            $("#vehicle-form [name=match_profile]").val('{"urban": 1}');
            vehicleModal.show();
        });
    });

    $(document).on("click", ".edit-vehicle", function () {
        const vehicle = state.vehicles.find(item => item.id === Number($(this).data("id")));
        loadCategories(function () {
            Object.keys(vehicle).forEach(function (key) {
                if (key === "category" || typeof vehicle[key] === "object") return;
                $("#vehicle-form [name=" + key + "]").val(vehicle[key]);
            });
            $("#vehicle-form [name=match_profile]").val(JSON.stringify(vehicle.match_profile || {}, null, 2));
            vehicleModal.show();
        });
    });

    $("#vehicle-form").on("submit", function (event) {
        event.preventDefault();
        const form = $(this);
        let profile;
        try { profile = JSON.parse(value(form, "match_profile")); }
        catch (_) { showAlert("O perfil de compatibilidade deve ser um JSON válido.", "danger"); return; }
        const id = value(form, "id");
        const numeric = ["category_id","year","price_min","price_max","seats","trunk_capacity","consumption_city","consumption_highway"];
        const payload = {match_profile: profile};
        form.serializeArray().forEach(function (field) {
            if (field.name === "id" || field.name === "match_profile") return;
            payload[field.name] = numeric.includes(field.name) ? Number(field.value) : field.value;
        });
        request({
            url: "/api/admin/vehicles" + (id ? "/" + id : ""),
            method: id ? "PUT" : "POST",
            contentType: "application/json",
            data: JSON.stringify(payload)
        }).done(function () {
            vehicleModal.hide(); loadVehicles(); showAlert("Veículo salvo.", "success");
        });
    });

    $("#new-question").on("click", function () {
        $("#question-form")[0].reset();
        $("#question-form [name=id]").val("");
        $("#question-form [name=weight]").val(1);
        $("#question-form [name=display_order]").val(state.questions.length + 1);
        questionModal.show();
    });

    $(document).on("click", ".edit-question", function () {
        const question = state.questions.find(item => item.id === Number($(this).data("id")));
        $("#question-form [name=id]").val(question.id);
        $("#question-form [name=text]").val(question.text);
        $("#question-form [name=weight]").val(question.weight);
        $("#question-form [name=display_order]").val(question.display_order);
        questionModal.show();
    });

    $("#question-form").on("submit", function (event) {
        event.preventDefault();
        const form = $(this);
        const id = value(form, "id");
        request({
            url: "/api/admin/questions" + (id ? "/" + id : ""),
            method: id ? "PUT" : "POST",
            contentType: "application/json",
            data: JSON.stringify({
                text: value(form, "text"),
                type: "SINGLE_CHOICE",
                weight: Number(value(form, "weight")),
                display_order: Number(value(form, "display_order"))
            })
        }).done(function () {
            questionModal.hide(); loadQuestions(); showAlert("Pergunta salva.", "success");
        });
    });

    $(document).on("click", ".new-option", function () {
        $("#option-form")[0].reset();
        $("#option-form [name=question_id]").val($(this).data("question-id"));
        $("#option-form [name=id]").val("");
        $("#option-form [name=score_profile]").val('{"urban": 1}');
        optionModal.show();
    });

    $(document).on("click", ".edit-option", function () {
        const question = state.questions.find(item => item.id === Number($(this).data("question-id")));
        const option = question.answer_options.find(item => item.id === Number($(this).data("id")));
        $("#option-form [name=question_id]").val(question.id);
        $("#option-form [name=id]").val(option.id);
        $("#option-form [name=text]").val(option.text);
        $("#option-form [name=display_order]").val(option.display_order);
        $("#option-form [name=score_profile]").val(JSON.stringify(option.score_profile, null, 2));
        optionModal.show();
    });

    $("#option-form").on("submit", function (event) {
        event.preventDefault();
        const form = $(this);
        const questionID = value(form, "question_id");
        const id = value(form, "id");
        let profile;
        try { profile = JSON.parse(value(form, "score_profile")); }
        catch (_) { showAlert("O perfil de score deve ser um JSON válido.", "danger"); return; }
        request({
            url: "/api/admin/questions/" + questionID + "/options" + (id ? "/" + id : ""),
            method: id ? "PUT" : "POST",
            contentType: "application/json",
            data: JSON.stringify({
                text: value(form, "text"),
                display_order: Number(value(form, "display_order")),
                score_profile: profile
            })
        }).done(function () {
            optionModal.hide(); loadQuestions(); showAlert("Opção salva.", "success");
        });
    });

    $(document).on("click", ".delete-vehicle, .delete-category, .delete-question, .delete-option", function () {
        if (!window.confirm("Deseja desativar este registro?")) return;
        const button = $(this);
        let url;
        if (button.hasClass("delete-option")) {
            url = "/api/admin/questions/" + button.data("question-id") + "/options/" + button.data("id");
        } else if (button.hasClass("delete-vehicle")) {
            url = "/api/admin/vehicles/" + button.data("id");
        } else if (button.hasClass("delete-category")) {
            url = "/api/admin/categories/" + button.data("id");
        } else {
            url = "/api/admin/questions/" + button.data("id");
        }
        request({url: url, method: "DELETE"}).done(function () {
            if (url.includes("/vehicles/")) loadVehicles();
            else if (url.includes("/categories/")) loadCategories();
            else loadQuestions();
            showAlert("Registro desativado.", "success");
        });
    });

    loadDashboard();
});
