var sb = document.getElementById("searchbox");

sb.addEventListener("keypress", function (e) {
	// If pressed enter
	if (e.keyCode == 13) {
		window.location = "https://google.com/#q=" + encodeURIComponent(this.value);

		return false
	}
});
