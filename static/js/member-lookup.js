(function () {
    const callsignInput = document.getElementById("callsign");
    if (!callsignInput) return;

    callsignInput.addEventListener("blur", function () {
        const callsign = callsignInput.value.trim();
        if (!callsign) return;

        fetch("/member-lookup?callsign=" + encodeURIComponent(callsign))
            .then(function (resp) { return resp.json(); })
            .then(function (data) {
                if (!data || !data.first_name) return;

                var firstname = document.getElementById("firstname");
                var lastname = document.getElementById("lastname");
                var email = document.getElementById("email");

                if (firstname.value === "") firstname.value = data.first_name;
                if (lastname.value === "") lastname.value = data.last_name || "";
                if (email.value === "") email.value = data.email || "";
                document.getElementById("nfarl").checked = true;
            })
            .catch(function (err) { console.warn("member lookup failed:", err); });
    });
})();
