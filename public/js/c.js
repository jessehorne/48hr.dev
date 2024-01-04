var applyAsBackend = `
    <a
        href=""
        id="{{ProjectID}}"
        data-which="{{Which}}"
        class="apply-button inline-flex items-center px-3 py-2 text-sm font-medium text-center text-white bg-bubble-gum rounded-lg hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300">
    Apply as Backend
    </a>
    `;

var applyAsFrontend = `
    <a
        href=""
        id="{{ProjectID}}"
        data-which="{{Which}}"
        class="apply-button inline-flex items-center px-3 py-2 text-sm font-medium text-center text-white bg-tahiti rounded-lg hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300">
    Apply as Frontend
    </a>
    `;

var applyAsInfra = `
    <a
        href=""
        id="{{ProjectID}}"
        data-which="{{Which}}"
        class="apply-button inline-flex items-center px-3 py-2 text-sm font-medium text-center text-white bg-purple rounded-lg hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300">
    Apply as Infrastructure
    </a>
    `;


var cardTmpl = `
   <div class="max-w-screen-xl rounded-lg border border-gray-200 overflow-hidden shadow-md bg-white p-4">
        <div class="p-4">
            <h5 class="text-2xl font-bold tracking-tight text-gray-900">{{Title}}</h5>
            <p class="mb-3 font-normal text-gray-700">{{Short}}</p>
            {{NeedsBackend}}
            {{NeedsFrontend}}
            {{NeedsInfra}}
        </div>
    </div>
   `

function buildNeeds(projectId, which) {
    var tmpl = "";

    if (which == "backend") {
        tmpl = applyAsBackend;
    } else if (which == "frontend") {
        tmpl = applyAsFrontend;
    } else if (which == "infra") {
        tmpl = applyAsInfra;
    }

    return tmpl.replace("{{ProjectID}}", projectId).replace("{{Which}}", which);
}

function buildCard(projectId, title, short, needsBackend, needsFrontend, needsInfra) {
    return cardTmpl
        .replace("{{Title}}", title ? title : "")
        .replace("{{Short}}", short ? short : "")
        .replace("{{NeedsBackend}}", needsBackend ? buildNeeds(projectId, "backend") : "")
        .replace("{{NeedsFrontend}}", needsFrontend ? buildNeeds(projectId, "frontend") : "")
        .replace("{{NeedsInfra}}", needsInfra ? buildNeeds(projectId, "infra") : "")
        .replace("{{ProjectID}}", projectId)
}

function addCard(projectId, title, short, needsBackend, needsFrontend, needsInfra) {
    $("#cards").append(
        $(
            buildCard(
                projectId,
                title,
                short,
                needsBackend, needsFrontend, needsInfra
            )
        )
    )
}

$(document).ready(function() {
    // Initialize Firebase

    firebase.initializeApp(firebaseConfig);
    firebase.analytics();

    // Initialize the FirebaseUI Widget using Firebase.
    var ui = new firebaseui.auth.AuthUI(firebase.auth());

    var firestore;
    if (firebase.firestore) {
        firestore = firebase.firestore();

        // get posts
        firestore.collection("posts").get().then((qs) => {
            if (qs.size > 0) {
                $("#cards").html("");
            }
          qs.forEach((p) => {
              const data = p.data();
              addCard(p.id, data.title, data.short, data.needBackend, data.needFrontend, data.needInfra);
          });
        })
    }

    firebase.auth().onAuthStateChanged(function(user) {
        if (user) {
            if (typeof initPost !== "undefined") {
                initPost(firestore, user.uid);
            }
            $("#login-link").addClass("hidden");
            $("#projects-link").removeClass("hidden");
        } else {
            $("#login-link").removeClass("hidden");
            $("#projects-link").addClass("hidden");
        }
    });

    var uiConfig = {
        callbacks: {
            signInSuccessWithAuthResult: function(authResult, redirectUrl) {
                // User successfully signed in.
                // Return type determines whether we continue the redirect automatically
                // or whether we leave that to developer to handle.
                return true;
            },
            uiShown: function() {
                // The widget is rendered.
                // Hide the loader.
                document.getElementById('loader').style.display = 'none';
            }
        },
        // Will use popup for IDP Providers sign-in flow instead of the default, redirect.
        signInFlow: 'popup',
        signInSuccessUrl: '/',
        signInOptions: [
            // Leave the lines as is for the providers you want to offer your users.
            firebase.auth.GoogleAuthProvider.PROVIDER_ID,
            firebase.auth.EmailAuthProvider.PROVIDER_ID,
        ],
        // Terms of service url.
        tosUrl: '<your-tos-url>',
        // Privacy policy url.
        privacyPolicyUrl: '<your-privacy-policy-url>'
    };


    if ($("#firebaseui").length) {
        ui.start('#firebaseui', uiConfig);
    }

    // addCard(555, "Pooper", "Log poops daily for science, built with Go.", true, true, true);
});
