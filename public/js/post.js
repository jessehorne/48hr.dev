function initPost(firestore, userID) {
    $("#post-form").submit(function(e) {
        e.preventDefault();
    });

    $("#create-button").click(function(e) {
        var title = $("#title-input").val();
        var short = $("#short-input").val();
        var needsBackend = $("#backend-input").is(":checked");
        var needsFrontend = $("#frontend-input").is(":checked");
        var needsInfra = $("#infra-input").is(":checked");

        if (title.length < 3 || short.length < 3) {
            return;
        }

        data = {
            "title": title,
            "short": short,
            "needsBackend": needsBackend,
            "needsFrontend": needsFrontend,
            "needsInfra": needsInfra,
        }

        console.log(data)

        // get user token
        firebase.auth().currentUser
            .getIdToken()
            .then(function (token) {
                $.ajax({
                    url: "/api/projects",
                    type: "POST",
                    data: JSON.stringify(data),
                    contentType: "application/json",
                    dataType: "json",
                    headers: {
                        "Authorization": "Bearer " + token
                    },
                    success: function(data) {
                        window.location.href = "/";
                    }
                })
            });
    });
}
