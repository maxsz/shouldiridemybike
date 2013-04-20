var reason;
var geocoder;
var address;

$(function() {
  $("#reason").hide()
  geocoder = new google.maps.Geocoder();
  loadLocation();
});

function handleError(error) {
  var excuse;

  if (typeof error !== 'undefined' && error !== null) {
    reason = error.toString();
    $("#reason").show();
  }
  excuse = "Sorry, I don't know :(<br />Just look out of the window.";
  $("#pos").html(excuse);
}

function loadLocation() {
  navigator.geolocation.getCurrentPosition(function(position){
    var latlng = new google.maps.LatLng(
      position.coords.latitude, 
      position.coords.longitude
    );

    geocoder.geocode({'latLng': latlng}, function(results, status) {
      if (status == google.maps.GeocoderStatus.OK) {
        if (results[0]) {
          address = results[0]["formatted_address"];
        }
      }
    });

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
          if (response["Error"] !== null) {
            handleError(response["Error"]);
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
    
  }, function(error){
    switch(error.code) 
      { 
        case error.PERMISSION_DENIED:
          reason = "You need to allow shouldiridemybike to use your location.";
          break;
        case error.POSITION_UNAVAILABLE:
          reason = "Your location information is unavailable";
          break;
        case error.TIMEOUT:
          reason = "You need to use a modern browser and allow shouldiridemybike to use your location.";
          break;
        case error.UNKNOWN_ERROR:
          reason = "You need to use a modern browser and allow shouldiridemybike to use your location.";
          break;
      }
    $("#pos").html("Sorry, I don't know where you are.");
    $("#reason").show();
  });
}

function why() {
  var poweredby;
  var location;

  poweredby = '<br />Powered by <a href="http://forecast.io">Forecast.io</a>.<br />';
  location = '<span id="location">' + address + '</span>.';

  $("#reason").show();
  $("#reason").html(reason + poweredby + location);
}