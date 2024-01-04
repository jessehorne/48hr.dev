function initPost(firestore, userID) {
    $("#post-form").submit(function(e) {
        e.preventDefault();
    });

    $("#create-button").click(function(e) {
        var title = $("#title-input").val();
        var short = $("#short-input").val();
        var needBackend = $("#backend-input").is(":checked");
        var needFrontend = $("#frontend-input").is(":checked");
        var needInfra = $("#infra-input").is(":checked");

        if (title.length < 3 || short.length < 3) {
            return;
        }

        firestore.collection("posts").add(
            {
                "user_id": userID,
                "title": title,
                "short": short,
                "needBackend": needBackend,
                "needFrontend": needFrontend,
                "needInfra": needInfra
            }
        ).then(function() {
            window.location.href = "/";
        });
    });
}
