var AddMemory = AddMemory || (function() {
    var map;
    var marker;

    var delay = function() {
        var timeout = 0;
        return function(callback, ms) {
            clearTimeout(timeout);
            timeout = setTimeout(callback, ms);
        };
    }();

    function initialize() {
        var mapOptions = {
            center: { lat: 38.6272222, lng: -90.1977778},
            zoom: 13
        };
        
        map = new google.maps.Map(document.getElementById('address_search_map'), mapOptions);
        
        $("#address_text").keyup(function() {
            delay(findAddress, 1500);
        });
      
        $("#address_text").keypress(function (e) {
            if ((e.which && e.which == 13) || (e.keyCode && e.keyCode == 13)) {
                findAddress();
            }
        });
    }

    function findAddress() {
        var addressText = $("#address_text").val();
        if (!addressText) {
            alert("Please enter an address to search for.");
            return;
        }

        if (marker) {
            marker.setMap(null);
        }
        
        var latitude = $("#latitude");
        var longitude = $("#longitude");
        var submitButton = $("#submit_button");
        var geocoder = new google.maps.Geocoder();
        
        $("#map_spinner").show();
        
        geocoder.geocode({address: addressText}, function(results, status) {
            $("#map_spinner").hide();
            
            if (status == google.maps.GeocoderStatus.OK) {
                map.setCenter(results[0].geometry.location);
                
                marker = new google.maps.Marker({
                    map: map,
                    position: results[0].geometry.location
                });
                
                latitude.val(results[0].geometry.location.lat());
                longitude.val(results[0].geometry.location.lng());
                
                submitButton.attr("disabled", false);
            } else {
                latitude.val("");
                longitude.val("");
                
                submitButton.attr("disabled", "disabled");
            }
        });
    }

    function verifyForm() {
        var latitude = $("#latitude").val();
        var longitude = $("#longitude").val();
        
        if (!latitude || !longitude) {
            alert("Please search and find an address.");
            $("#address_text").focus();
            return false;
        }
        
        var titleField = $("#title");
        if (!titleField.val()) {
            alert("Please enter a title.");
            titleField.focus();
            return false;
        }
        
        var detailsField = $("#details");
        if (!detailsField.val()) {
            alert("Please enter some details.");
            detailsField.focus();
            return false;
        }
        
        return true;
    }

    google.maps.event.addDomListener(window, 'load', initialize);
})();