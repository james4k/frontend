(function (global) {
	"use strict";

	// Basic parallax via background-position-y
	global.initParallax = function() {
		var elems = Sizzle('.parallax');
		var update = function() {
			var yoffset = window.pageYOffset;
			var height = document.documentElement.clientHeight;
			elems.forEach(function(e) {
				var offs = yoffset - Math.max(0, e.offsetTop - height);
				var speed = Math.min(100.0 / height, 1/10.0);
				e.style.backgroundPositionY = -offs * speed + 'px';
			});
		};
		// Disable on mobile. Typically, scroll events are not fired until the
		// scroll physics animation stops. So, experience is very poor and
		// probably confusing. There are hacks out there for this, but meh.
		if (navigator.userAgent.indexOf("Mobile") < 0) {
			window.onscroll = update;
			window.onresize = update;
		}
	};
})(window);

