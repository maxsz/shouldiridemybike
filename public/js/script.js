var reason;

$(function() {
	$("#reason").hide()
	loadLocation();
});

function handleError() {
	var excuse = "Sorry, I don't know :(";
	$("#pos").html(excuse);
}

function loadLocation() {
	navigator.geolocation.getCurrentPosition(function(position){
		$.get('/shouldi', 
			{latitude: position.coords.latitude,
				longitude: position.coords.longitude}, 
		  	function(data, textStatus, xhr) {
		  		var response;
		  		
		  		if ("success" != textStatus) {
		  			handleError();
		  			return;
		  		}
		  		try {
		  			response = $.parseJSON(data);
		  		} catch(err) {
		  			handleError();
		  			return;
		  		}
		  		if (typeof response != 'object') {
		  			handleError();
		  			return;
		  		}

		  		reason = response["Reason"];
		  		if (response["Result"]) {
		  			$("#pos").html("Yes.");
		  		} else {
		  			$("#pos").html("No.");
		  		}
		  		$("#reason").show();
		});
		
	}, function(){
		$("#pos").html("Sorry, I don't know where you are.");
	});
}

function why() {
	$("#reason").show();
	$("#reason").html(reason);
}