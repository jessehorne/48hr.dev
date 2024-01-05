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
            <div class="flex flex-row">
            <h5 class="basis-1/3 text-2xl font-bold tracking-tight text-gray-900">{{Title}}</h5>
            <p class="basis-2/3 text-right">created on {{CreatedAt}}</p>
            </div>
            <p class="mb-3 font-normal text-gray-700">{{Short}}</p>
            {{NeedsBackend}}
            {{NeedsFrontend}}
            {{NeedsInfra}}
            <p class="mb-3 font-normal"><a href="/user/{{UserID}}">by {{DisplayName}}</a></p>
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

function buildCard(p) {
    const options = { year: 'numeric', month: 'long', day: 'numeric' };
    const formattedDate = new Intl.DateTimeFormat('en-US', options).format(p.CreatedAt);

    return cardTmpl
        .replace("{{Title}}", p.Title ? p.Title : "")
        .replace("{{Short}}", p.Short ? p.Short : "")
        .replace("{{NeedsBackend}}", p.NeedsBackend ? buildNeeds(p.ProjectId, "backend") : "")
        .replace("{{NeedsFrontend}}", p.NeedsFrontend ? buildNeeds(p.ProjectId, "frontend") : "")
        .replace("{{NeedsInfra}}", p.NeedsInfra ? buildNeeds(p.ProjectId, "infra") : "")
        .replace("{{ProjectID}}", p.ProjectId)
        .replace("{{CreatedAt}}", formattedDate)
        .replace("{{UserID}}", p.UserID)
        .replace("{{DisplayName}}", p.DisplayName)
}

function addCard(p) {
    if (p.title != "Centrifuge") {
        $("#cards").append($(buildCard(p)))
    }
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
              addCard(data);
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
});
