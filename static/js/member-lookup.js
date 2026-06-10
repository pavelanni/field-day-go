(function () {
    const callsignInput = document.getElementById("callsign");
    if (!callsignInput) return;

    let debounceTimer;

    callsignInput.addEventListener("blur", function () {
        const callsign = callsignInput.value.trim();
        if (!callsign) return;

        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(function () {
            fetch("/member-lookup?callsign=" + encodeURIComponent(callsign))
                .then(function (resp) { return resp.json(); })
                .then(function (data) {
                    if (!data || !data.first_name) return;

                    document.getElementById("firstname").value = data.first_name;
                    document.getElementById("lastname").value = data.last_name || "";
                    document.getElementById("email").value = data.email || "";
                    document.getElementById("nfarl").checked = true;
                })
                .catch(function () {});
        }, 300);
    });
})();
