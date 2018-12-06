// Copyright (c) 2013 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// chrome.runtime.onConnect.addListener(function(port) {
//   port.onMessage.addListener(function(msg) {
//     port.postMessage({counter: msg.counter+1});
//   });
// });
//
// chrome.extension.onRequest.addListener(
//   function(request, sender, sendResponse) {
//     sendResponse({counter: request.counter+1});
//   });

var glob

function addOverlay(el, id, autostart) {
  var bbox = el.getBBox()
  var div = document.createElement("div")
  
  var svg = document.querySelector(".punch-viewer-svgpage-svgcontainer > svg")
  
  div.style.position = "absolute"
  // div.style.border = "7px solid #fdf6e3"
  // div.style.boxSizing = "border-box"
  div.style.left = (bbox.x/svg.viewBox.baseVal.width*100)+"%"
  div.style.top = (bbox.y/svg.viewBox.baseVal.height*100)+"%"
  div.style.width = (bbox.width/svg.viewBox.baseVal.width*100)+"%"
  div.style.height = (bbox.height/svg.viewBox.baseVal.height*100)+"%"
  div.style.background = "rgba(253,246,227,.6)" // #fdf6e3"
  el._overlay_term = div
  div.id = "foobar"
	
	var added = false
  
	// div.addEventListener("click", function(e) {
	// 	if (added) {
	// 		console.log("focused", div.querySelector("iframe").contentWindow)
	// 		div.querySelector("iframe").contentWindow.focus()
	// 	}
	// 	e.stopPropagation()
	// 	e.preventDefault()
	// })
	
	function addIframe(e) {
		if (added) {
			return
		}
		top.navigator.keyboard.lock()
    div.innerHTML='<iframe style="border:0" src="http://127.0.0.1:8080/?demo='+id+'" width="100%" height="100%"></iframe>'
    if (e) {
			e.stopPropagation()
		}
		added = true
		
		var refresh = document.createElement("div")
	  refresh.style.position = "absolute"
	  refresh.style.left = "98%"
	  refresh.style.top = "100%"
	  refresh.style.width = "15px"
	  refresh.style.height = "15px"
	 	refresh.style.background = "rgba(0, 0, 0, .1)"
		div.appendChild(refresh)
		
	  refresh.addEventListener("mousedown", function() {
	  	added = false
			addIframe()
	  });		
	}
	
	if (autostart) {
		addIframe()
	} else {
	  div.addEventListener("mousedown", addIframe);		
	}
	
	glob = div
  

  document.querySelector(".punch-viewer-svgpage").append(div)
  // box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19)
}

function scanListener() {
  var terms = document.querySelectorAll(".punch-viewer-svgpage a")
  terms.forEach(function(term) {
    if (term._skip_term) {
      return
    }
    if (term._overlay_term) {
      if (document.body.contains(term._overlay_term)) {
        return
      }
    }
		var v = term.getAttribute("xlink:href").match(/.*^#term=(.+)$/)
    if (v) {
			var rest = v[1]
			rest = rest.replace(/,autostart$/g, "")
      addOverlay(term, rest, v[1].match(/,autostart$/) != null)
    } else {
      term._skiip_term = true
    }
  })
}

function setupListener() {
  setInterval(scanListener, 100)
  console.log("setupListener", document.URL, document.body)
}

setupListener()